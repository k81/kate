package knet

import (
	"context"
	"net"
)

type TCPServerConfig ServerConfig

type TCPServer struct {
	*ServerBase
}

func TCPListen(addr string) (ln net.Listener, err error) {
	if ln, err = net.Listen("tcp", addr); err != nil {
		return
	}
	ln = &TCPListener{ln.(*net.TCPListener)}
	return
}

func NewTCPServerConfig() *TCPServerConfig {
	conf := &TCPServerConfig{}
	conf.Io.SendQueueSize = 16
	conf.Io.RecvQueueSize = 16
	return conf
}

func NewTCPServer(ctx context.Context, conf *TCPServerConfig) *TCPServer {
	srv := &TCPServer{
		ServerBase: NewServerBase(ctx, TCPListen, (*ServerConfig)(conf)),
	}
	return srv
}
