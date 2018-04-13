package knet

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"

	"github.com/k81/kate/log"
	"github.com/k81/kate/utils"
)

var ErrServerClosed = errors.New("server closed")

type ListenFunc func(addr string) (net.Listener, error)

type ServerConfig struct {
	Io struct {
		SendQueueSize int
		RecvQueueSize int
		ReadTimeout   time.Duration
		WriteTimeout  time.Duration
	}
	MaxConnection int
}

func NewServerConfig() *ServerConfig {
	conf := &ServerConfig{}
	conf.Io.SendQueueSize = 16
	conf.Io.RecvQueueSize = 16
	return conf
}

type ServerBase struct {
	*IoServiceBase
	listen ListenFunc
	conf   *ServerConfig

	ctx       context.Context
	cancel    context.CancelFunc
	closeOnce sync.Once
	wg        sync.WaitGroup
}

func NewServerBase(ctx context.Context, listen ListenFunc, conf *ServerConfig) *ServerBase {
	var (
		newctx, cancel = context.WithCancel(ctx)
		ioConf         = &IoConfig{}
	)

	*ioConf = conf.Io

	srv := &ServerBase{
		IoServiceBase: NewIoServiceBase(ioConf),
		listen:        listen,
		conf:          conf,
		ctx:           newctx,
		cancel:        cancel,
	}
	return srv
}

func (srv *ServerBase) ListenAndServe(addr string) error {
	if srv.listen == nil {
		panic("no listen func defined")
	}

	ln, err := srv.listen(addr)
	if err != nil {
		return err
	}

	if srv.conf.MaxConnection > 0 {
		ln = LimitListener(ln, srv.conf.MaxConnection)
	}
	return srv.Serve(ln)
}

func (srv *ServerBase) Serve(l net.Listener) error {
	defer l.Close()

	var (
		tempDelay time.Duration
		conn      net.Conn
		err       error
	)

	for {
		if conn, err = l.Accept(); err != nil {

			select {
			case <-srv.ctx.Done():
				return ErrServerClosed
			default:
			}

			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.Debug(srv.ctx, "Accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}

			return err
		}

		tempDelay = 0
		session := srv.newSession(conn)
		srv.wg.Add(1)
		go srv.serve(session)
	}
}

func (srv *ServerBase) Close() {
	srv.closeOnce.Do(func() {
		srv.cancel()
		srv.wg.Wait()
	})
}

func (srv *ServerBase) AddRef() {
	srv.wg.Add(1)
}

func (srv *ServerBase) DecRef() {
	srv.wg.Done()
}

func (srv *ServerBase) newSession(conn net.Conn) *IoSession {
	session := NewIoSession(srv.ctx, srv, conn)
	return session
}

func (srv *ServerBase) serve(session *IoSession) {
	defer func() {
		if r := recover(); r != nil {
			log.Error(srv.ctx, "got panic in serve session", "error", r, "stack", utils.GetPanicStack())
		}
		srv.wg.Done()
	}()

	session.Open()
}
