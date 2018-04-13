package knet

import "sync/atomic"

type IoServiceBase struct {
	nextSessionId uint64
	conf          *IoConfig
	protocol      Protocol
	handler       IoHandler
}

func NewIoServiceBase(conf *IoConfig) *IoServiceBase {
	return &IoServiceBase{
		conf: conf,
	}
}

func (srv *IoServiceBase) SetIoHandler(h IoHandler) {
	srv.handler = h
}

func (srv *IoServiceBase) Protocol() Protocol {
	return srv.protocol
}

func (srv *IoServiceBase) SetProtocol(p Protocol) {
	srv.protocol = p
}

func (srv *IoServiceBase) IoHandler() IoHandler {
	return srv.handler
}

func (srv *IoServiceBase) IoConfig() *IoConfig {
	return srv.conf
}

func (srv *IoServiceBase) AddRef() {
}

func (srv *IoServiceBase) DecRef() {
}

func (srv *IoServiceBase) NextSessionId() uint64 {
	return atomic.AddUint64(&srv.nextSessionId, 1)
}
