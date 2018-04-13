package knet

import "net"

type Server interface {
	ListenAndServe(addr string) error
	Serve(ln net.Listener) error
	Close()
}
