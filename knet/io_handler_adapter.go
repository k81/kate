package knet

type IoHandlerAdapter struct {
}

func (h *IoHandlerAdapter) OnConnected(*IoSession) error {
	return nil
}

func (h *IoHandlerAdapter) OnDisconnected(*IoSession) {
}

func (h *IoHandlerAdapter) OnIdle(*IoSession) error {
	return nil
}

func (h *IoHandlerAdapter) OnError(*IoSession, error) {
}

func (h *IoHandlerAdapter) OnMessage(*IoSession, Message) error {
	return nil
}
