package knet

import (
	"context"
	"time"
)

type Client interface {
	Dial(addr string) error
	Close()
	Disconnect()
	SetProtocol(Protocol)
	SetIoHandler(h IoHandler)
	Call(ctx context.Context, req Request) (Response, error)
	CallWithTimeout(ctx context.Context, req Request, timeout time.Duration) (Response, error)
	Send(ctx context.Context, msg Message) error
	SendWithTimeout(ctx context.Context, msg Message, timeout time.Duration) error
	GetSession() *IoSession
	IsClosed() bool
	IsConnected() bool
}
