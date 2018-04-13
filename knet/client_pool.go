package knet

import (
	"context"
	"errors"
	"sync"
	"time"
)

var ErrClientPoolClosed = errors.New("client manager closed")

type ClientPoolConfig struct {
	IdleMin int `json:"idle_min"`
	IdleMax int `json:"idle_max"`
	Max     int `json:"max"`
}

type ClientFactory interface {
	NewClient() (Client, error)
}

type ClientPool struct {
	sync.Mutex
	conf     ClientPoolConfig
	factory  ClientFactory
	num      int
	freeList chan Client
	closed   bool

	ctx       context.Context
	cancel    context.CancelFunc
	closeOnce sync.Once
}

func NewClientPool(ctx context.Context, factory ClientFactory, conf ClientPoolConfig) *ClientPool {
	newctx, cancel := context.WithCancel(ctx)

	if conf.IdleMax < conf.Max {
		conf.IdleMax = conf.Max
	}

	if conf.IdleMin > conf.IdleMax {
		conf.IdleMin = conf.IdleMax
	}

	p := &ClientPool{
		conf:     conf,
		factory:  factory,
		freeList: make(chan Client, conf.IdleMax),
		ctx:      newctx,
		cancel:   cancel,
	}
	return p
}

func (p *ClientPool) Open() error {
	if p.factory == nil {
		panic("client factory not defined")
	}

	for i := 0; i < p.conf.IdleMin; i++ {
		client, err := p.factory.NewClient()
		if err != nil {
			return err
		}

		p.Lock()
		p.num++
		p.Unlock()
		p.freeList <- client
	}

	return nil
}

func (p *ClientPool) Close() {
	p.closeOnce.Do(func() {
		p.Lock()
		p.closed = true
		p.Unlock()

		p.cancel()

		close(p.freeList)

		for c := range p.freeList {
			c.Close()
		}

	})
}

func (p *ClientPool) Get() Client {
	c, err := p.get()
	if err != nil {
		return &errClient{err}
	}

	return &pooledClient{
		p:      p,
		Client: c,
	}
}

func (p *ClientPool) get() (c Client, err error) {
	if p.isClosed() {
		return nil, ErrClientPoolClosed
	}

	select {
	case c = <-p.freeList:
		return
	default:
	}

	if p.isMaxReached() {
		select {
		case <-p.ctx.Done():
			return nil, ErrClientPoolClosed
		case c = <-p.freeList:
		}
		return
	}

	if c, err = p.factory.NewClient(); err != nil {
		return
	}

	p.Lock()
	p.num++
	p.Unlock()
	return
}

func (p *ClientPool) put(c Client) {
	if p.isClosed() {
		c.Close()
		return
	}

	select {
	case p.freeList <- c:
		return
	default:
	}

	p.Lock()
	p.num--
	p.Unlock()

	c.Close()
}

func (p *ClientPool) isClosed() (closed bool) {
	p.Lock()
	closed = p.closed
	p.Unlock()
	return
}

func (p *ClientPool) isMaxReached() (reached bool) {
	p.Lock()
	if p.num >= p.conf.Max {
		reached = true
	}
	p.Unlock()
	return
}

type pooledClient struct {
	p *ClientPool
	Client
}

func (pc *pooledClient) Close() {
	c := pc.Client
	if _, ok := c.(*errClient); ok {
		return
	}

	pc.Client = &errClient{ErrClientClosed}
	pc.p.put(c)
}

type errClient struct {
	err error
}

func (c *errClient) Dial(addr string) error                          { return c.err }
func (c *errClient) Close()                                          {}
func (c *errClient) Disconnect()                                     {}
func (c *errClient) SetProtocol(Protocol)                            {}
func (c *errClient) SetIoHandler(IoHandler)                          {}
func (c *errClient) Call(context.Context, Request) (Response, error) { return nil, c.err }
func (c *errClient) CallWithTimeout(context.Context, Request, time.Duration) (Response, error) {
	return nil, c.err
}
func (c *errClient) Send(context.Context, Message) error { return c.err }
func (c *errClient) SendWithTimeout(context.Context, Message, time.Duration) error {
	return c.err
}
func (c *errClient) GetSession() *IoSession { return nil }
func (c *errClient) IsClosed() bool         { return true }
func (c *errClient) IsConnected() bool      { return false }
