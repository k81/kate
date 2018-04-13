package knet

type IoHandler interface {
	OnConnected(*IoSession) error
	OnDisconnected(*IoSession)
	OnIdle(*IoSession) error
	OnError(*IoSession, error)
	OnMessage(*IoSession, Message) error
}
