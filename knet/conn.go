package knet

import (
	"net"
	"sync/atomic"
	"time"
)

type Conn struct {
	net.Conn
	readTimeout  time.Duration
	writeTimeout time.Duration
	bytesIn      uint32
	bytesOut     uint32
}

func (c *Conn) SetTimeout(d time.Duration) {
	c.readTimeout = d
	c.writeTimeout = d
}

func (c *Conn) SetReadTimeout(d time.Duration) (err error) {
	c.readTimeout = d

	if d == 0 {
		if err = c.Conn.SetReadDeadline(time.Time{}); err != nil {
			return
		}
	}
	return
}

func (c *Conn) SetWriteTimeout(d time.Duration) (err error) {
	c.writeTimeout = d

	if d == 0 {
		if err = c.Conn.SetWriteDeadline(time.Time{}); err != nil {
			return
		}
	}
	return
}

func (c *Conn) Read(b []byte) (n int, err error) {
	if c.readTimeout > 0 {
		if err = c.Conn.SetReadDeadline(time.Now().Add(c.readTimeout)); err != nil {
			return
		}
	}
	n, err = c.Conn.Read(b)
	atomic.AddUint32(&c.bytesIn, uint32(n))
	return
}

func (c *Conn) Write(b []byte) (n int, err error) {
	if c.writeTimeout > 0 {
		if err = c.Conn.SetWriteDeadline(time.Now().Add(c.writeTimeout)); err != nil {
			return
		}
	}
	n, err = c.Conn.Write(b)
	atomic.AddUint32(&c.bytesOut, uint32(n))
	return
}

func (c *Conn) GetReadBytes() uint32 {
	return atomic.LoadUint32(&c.bytesIn)
}

func (c *Conn) GetWriteBytes() uint32 {
	return atomic.LoadUint32(&c.bytesOut)
}
