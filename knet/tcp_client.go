package knet

import (
	"context"
	"net"
	"time"
)

type TCPClientConfig struct {
	Io struct {
		SendQueueSize int
		RecvQueueSize int
		ReadTimeout   time.Duration
		WriteTimeout  time.Duration
	}
	DialTimeout   time.Duration
	AutoReconnect bool
}

func NewTCPClientConfig() *TCPClientConfig {
	conf := &TCPClientConfig{}
	conf.Io.SendQueueSize = 16
	conf.Io.RecvQueueSize = 16
	conf.DialTimeout = 30 * time.Second
	return conf
}

type TCPClient struct {
	*ClientBase
}

func TCPDialFunc(timeout time.Duration) DialFunc {
	return func(addr string) (conn net.Conn, err error) {
		if timeout > 0 {
			return net.DialTimeout("tcp", addr, timeout)
		}
		return net.Dial("tcp", addr)
	}
}

func NewTCPClient(ctx context.Context, conf *TCPClientConfig) *TCPClient {
	clientConf := &ClientConfig{}
	clientConf.Io = conf.Io
	clientConf.AutoReconnect = conf.AutoReconnect

	c := &TCPClient{
		ClientBase: NewClientBase(ctx, TCPDialFunc(conf.DialTimeout), clientConf),
	}
	return c
}
