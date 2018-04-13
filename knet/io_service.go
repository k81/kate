package knet

type IoService interface {
	Protocol() Protocol
	IoHandler() IoHandler
	IoConfig() *IoConfig
	AddRef()
	DecRef()
	NextSessionId() uint64
}
