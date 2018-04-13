package knet

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"
)

var (
	ErrClientClosed       = errors.New("client closed")
	ErrClientDisconnected = errors.New("client disconnected")
)

type DialFunc func(addr string) (net.Conn, error)

type ClientConfig struct {
	Io struct {
		SendQueueSize int
		RecvQueueSize int
		ReadTimeout   time.Duration
		WriteTimeout  time.Duration
	}
	AutoReconnect bool
}

func NewClientConfig() *ClientConfig {
	conf := &ClientConfig{}
	conf.Io.SendQueueSize = 16
	conf.Io.RecvQueueSize = 16
	return conf
}

type pendingRequest struct {
	Request
	errorCh chan error
	replyCh chan Response
}

type ClientBase struct {
	*IoServiceBase
	sync.Mutex
	remoteAddr   string
	conf         *ClientConfig
	handler      IoHandler
	dial         DialFunc
	session      *IoSession
	sessionError error
	closed       bool
	pendingLock  sync.Mutex
	pendingMap   map[uint64]*pendingRequest

	ctx       context.Context
	cancel    context.CancelFunc
	closeOnce sync.Once
}

func NewClientBase(ctx context.Context, dial DialFunc, conf *ClientConfig) *ClientBase {
	var (
		newctx, cancel = context.WithCancel(ctx)
		ioConf         = &IoConfig{}
	)

	*ioConf = conf.Io

	c := &ClientBase{
		IoServiceBase: NewIoServiceBase(ioConf),
		conf:          conf,
		dial:          dial,
		pendingMap:    make(map[uint64]*pendingRequest),
		ctx:           newctx,
		cancel:        cancel,
	}
	c.IoServiceBase.SetIoHandler(c)

	return c
}

func (c *ClientBase) SetIoHandler(h IoHandler) {
	c.handler = h
}

func (c *ClientBase) Dial(addr string) (err error) {
	if c.dial == nil {
		panic("not dail func defined")
	}

	c.remoteAddr = addr
	return c.ensureConnected(true)
}

func (c *ClientBase) Close() {
	c.closeOnce.Do(func() {
		c.Lock()
		c.closed = true
		c.Unlock()

		session := c.GetSession()
		if session != nil {
			session.Close()
		}

		c.cancel()
	})
}

func (c *ClientBase) Send(ctx context.Context, msg Message) error {
	return c.SendWithTimeout(ctx, msg, 0)
}

func (c *ClientBase) SendWithTimeout(ctx context.Context, msg Message, timeout time.Duration) (err error) {
	if c.IsClosed() {
		return ErrClientClosed
	}

	if err = c.ensureConnected(false); err != nil {
		return
	}

	return c.GetSession().SendWithTimeout(ctx, msg, timeout)
}

func (c *ClientBase) Call(ctx context.Context, req Request) (Response, error) {
	return c.CallWithTimeout(ctx, req, 0)
}

func (c *ClientBase) CallWithTimeout(ctx context.Context, req Request, timeout time.Duration) (resp Response, err error) {
	if c.IsClosed() {
		return nil, ErrClientClosed
	}

	if err = c.ensureConnected(false); err != nil {
		return
	}

	pendingReq := &pendingRequest{
		Request: req,
		replyCh: make(chan Response, 1),
		errorCh: make(chan error, 1),
	}

	c.pendingLock.Lock()
	c.pendingMap[req.Id()] = pendingReq
	c.pendingLock.Unlock()

	tBegin := time.Now()

	if err = c.GetSession().SendWithTimeout(ctx, req, timeout); err != nil {
		return
	}

	if timeout == 0 {
		select {
		case <-c.ctx.Done():
			return nil, ErrClientClosed
		case <-ctx.Done():
			return nil, ctx.Err()
		case resp = <-pendingReq.replyCh:
		case err = <-pendingReq.errorCh:
		}
	} else {
		wait := time.Since(tBegin)

		select {
		case <-c.ctx.Done():
			return nil, ErrClientClosed
		case <-ctx.Done():
			return nil, ctx.Err()
		case resp = <-pendingReq.replyCh:
		case err = <-pendingReq.errorCh:
		case <-time.After(wait):
			{
				c.pendingLock.Lock()
				delete(c.pendingMap, req.Id())
				c.pendingLock.Unlock()
				err = ErrTimeout
			}
		}
	}
	return
}

func (c *ClientBase) GetSession() (session *IoSession) {
	c.Lock()
	session = c.session
	c.Unlock()
	return
}

func (c *ClientBase) Disconnect() {
	session := c.GetSession()
	c.OnError(session, ErrClientDisconnected)
	session.Close()
}

func (c *ClientBase) IsConnected() bool {
	session := c.GetSession()
	return session != nil && session.IsConnected()
}

func (c *ClientBase) IsClosed() (closed bool) {
	c.Lock()
	closed = c.closed
	c.Unlock()
	return
}

func (c *ClientBase) OnConnected(session *IoSession) error {
	c.Lock()
	c.session = session
	c.sessionError = nil
	c.Unlock()

	if h := c.handler; h != nil {
		return h.OnConnected(session)
	}
	return nil
}

func (c *ClientBase) OnDisconnected(session *IoSession) {
	if c.IsClosed() {
		return
	}

	c.settlePendingRequests()

	if h := c.handler; h != nil {
		h.OnDisconnected(session)
	}
}

func (c *ClientBase) OnIdle(session *IoSession) error {
	if h := c.handler; h != nil {
		return h.OnIdle(session)
	}
	return nil
}

func (c *ClientBase) OnError(session *IoSession, err error) {

	c.Lock()
	c.sessionError = err
	c.Unlock()

	if h := c.handler; h != nil {
		h.OnError(session, err)
	}
}

func (c *ClientBase) OnMessage(session *IoSession, msg Message) error {

	if resp, ok := msg.(Response); ok {
		c.pendingLock.Lock()
		req, exists := c.pendingMap[resp.Id()]
		delete(c.pendingMap, resp.Id())
		c.pendingLock.Unlock()

		if exists {
			req.replyCh <- resp
			return nil
		}
	}

	if h := c.handler; h != nil {
		return h.OnMessage(session, msg)
	}
	return nil
}

func (c *ClientBase) ensureConnected(force bool) (err error) {
	var (
		conn    net.Conn
		session *IoSession
	)

	if c.IsConnected() {
		return
	}

	if !force {
		if !c.conf.AutoReconnect {
			return ErrClientDisconnected
		}
	}

	if conn, err = c.dial(c.remoteAddr); err != nil {
		return
	}

	session = NewIoSession(c.ctx, c, conn)
	session.Open()
	return
}

func (c *ClientBase) settlePendingRequests() {
	var (
		pendingRequests map[uint64]*pendingRequest
		newPendingMap   = make(map[uint64]*pendingRequest)
		err             error
	)

	c.pendingLock.Lock()
	pendingRequests = c.pendingMap
	c.pendingMap = newPendingMap
	err = c.sessionError
	c.pendingLock.Unlock()

	if err == nil {
		err = ErrClientClosed
	}

	for _, req := range pendingRequests {
		req.errorCh <- err
	}

	return
}
