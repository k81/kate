package knet

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/k81/kate/log"
	"github.com/k81/kate/utils"
)

var (
	ErrSessionClosed = errors.New("session closed")
	ErrPeerDead      = errors.New("peer dead")
	ErrTimeout       = errors.New("timeout")
)

type IoSession struct {
	id        uint64
	srv       IoService
	conf      *IoConfig
	handler   IoHandler
	protocol  Protocol
	conn      *Conn
	attrs     map[string]interface{}
	attrsLock sync.RWMutex
	sendQ     chan Message
	recvQ     chan Message

	idleCount     uint32
	readMsgCount  uint32
	writeMsgCount uint32

	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	connected uint32
	closed    uint32
}

func NewIoSession(ctx context.Context, srv IoService, conn net.Conn) *IoSession {
	var (
		id             = srv.NextSessionId()
		newctx, cancel = context.WithCancel(ctx)
	)

	newctx = log.SetContext(newctx, "session", id)

	s := &IoSession{
		id:       id,
		srv:      srv,
		handler:  srv.IoHandler(),
		conf:     srv.IoConfig(),
		protocol: srv.Protocol(),
		conn:     &Conn{Conn: conn},
		attrs:    make(map[string]interface{}),
		ctx:      newctx,
		cancel:   cancel,
		sendQ:    make(chan Message, srv.IoConfig().SendQueueSize),
		recvQ:    make(chan Message, srv.IoConfig().RecvQueueSize),
	}

	_ = s.conn.SetReadTimeout(s.conf.ReadTimeout)
	_ = s.conn.SetWriteTimeout(s.conf.WriteTimeout)

	return s
}

func (s *IoSession) Id() uint64 {
	return s.id
}

func (s *IoSession) IoService() IoService {
	return s.srv
}

func (s *IoSession) Context() context.Context {
	return s.ctx
}

func (s *IoSession) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *IoSession) LocalAddr() net.Addr {
	return s.conn.LocalAddr()
}

func (s *IoSession) GetAttr(key string) (v interface{}) {
	s.attrsLock.RLock()
	v = s.attrs[key]
	s.attrsLock.RUnlock()
	return
}

func (s *IoSession) SetAttr(key string, value interface{}) {
	s.attrsLock.Lock()
	s.attrs[key] = value
	s.attrsLock.Unlock()
}

func (s *IoSession) RemoveAttr(key string) {
	s.attrsLock.Lock()
	delete(s.attrs, key)
	s.attrsLock.Unlock()
}

func (s *IoSession) Open() {
	s.srv.AddRef()
	s.wg.Add(3)
	go s.handleLoop()
	go s.readLoop()
	go s.writeLoop()
	atomic.StoreUint32(&s.connected, 1)

	if err := s.handler.OnConnected(s); err != nil {
		s.Close()
	}
}

func (s *IoSession) Close() {
	if !atomic.CompareAndSwapUint32(&s.closed, 0, 1) {
		return
	}

	atomic.StoreUint32(&s.connected, 0)
	s.cancel()
	s.conn.Close()

	go func() {
		s.wg.Wait()

		close(s.sendQ)
		close(s.recvQ)

		s.handler.OnDisconnected(s)
		s.srv.DecRef()
	}()
}

func (s *IoSession) IsConnected() bool {
	return atomic.LoadUint32(&s.connected) == 1
}

func (s *IoSession) IsClosed() bool {
	return atomic.LoadUint32(&s.closed) == 1
}

func (s *IoSession) Send(ctx context.Context, m Message) error {
	return s.SendWithTimeout(ctx, m, 0)
}

func (s *IoSession) SendWithTimeout(ctx context.Context, m Message, timeout time.Duration) error {
	if timeout == 0 {
		select {
		case <-s.ctx.Done():
			return ErrSessionClosed
		case <-ctx.Done():
			return ctx.Err()
		case s.sendQ <- m:
		}
	} else {
		select {
		case <-s.ctx.Done():
			return ErrSessionClosed
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(timeout):
			return ErrTimeout
		case s.sendQ <- m:
		}
	}

	return nil
}

func (s *IoSession) GetIdleCount() uint32 {
	return atomic.LoadUint32(&s.idleCount)
}

func (s *IoSession) String() string {
	return fmt.Sprintf("session %d, Read Byte Count: %d, Write Byte Count: %d, Read Msg Count: %d, Write Msg Count: %d",
		s.id,
		s.conn.GetReadBytes(),
		s.conn.GetWriteBytes(),
		atomic.LoadUint32(&s.readMsgCount),
		atomic.LoadUint32(&s.writeMsgCount),
	)
}

func (s *IoSession) handleLoop() {
	var (
		m   Message
		err error
	)

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("got panic in handle loop: error=%v, stack=%v", r, utils.GetPanicStack())
		}

		if err != nil {
			s.handler.OnError(s, err)
		}

		s.wg.Done()
		s.Close()
	}()

	for {
		select {
		case <-s.ctx.Done():
			return

		case m = <-s.recvQ:
			if err = s.handler.OnMessage(s, m); err != nil {
				return
			}
		}
	}
}

func (s *IoSession) readLoop() {
	var (
		m   Message
		err error
	)

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("got panic in read loop: error=%v, stack=%v", r, utils.GetPanicStack())
		}

		if !s.IsClosed() && err != nil && err != io.EOF {
			s.handler.OnError(s, err)
		}

		s.wg.Done()
		s.Close()
	}()

	for {
		select {
		case <-s.ctx.Done():
			return

		default:
		}

		if m, err = s.protocol.Decode(s, s.conn); err != nil {
			if e, ok := err.(net.Error); ok && e.Timeout() {
				atomic.AddUint32(&s.idleCount, 1)

				if err = s.handler.OnIdle(s); err != nil {
					return
				}

				continue
			}
			return
		}

		atomic.StoreUint32(&s.idleCount, 0)
		atomic.AddUint32(&s.readMsgCount, 1)

		select {
		case <-s.ctx.Done():
		case s.recvQ <- m:
		}
	}
}

func (s *IoSession) writeLoop() {
	var (
		m    Message
		data []byte
		err  error
	)

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("got panic in write loop: error=%v, stack=%v", r, utils.GetPanicStack())
		}

		if !s.IsClosed() && err != nil {
			s.handler.OnError(s, err)
		}

		s.wg.Done()
		s.Close()
	}()

	for {
		select {
		case <-s.ctx.Done():
			return
		case m = <-s.sendQ:
			if data, err = s.protocol.Encode(s, m); err != nil {
				return
			}

			if _, err = s.conn.Write(data); err != nil {
				return
			}
			atomic.AddUint32(&s.writeMsgCount, 1)
		}
	}
}
