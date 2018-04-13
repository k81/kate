package knet

import "io"

type Message interface{}

type Request interface {
	Message
	Id() uint64
}

type Response interface {
	Message
	Id() uint64
}

type ProtocolDecoder interface {
	// if not full message is read, should return (nil, nil)
	Decode(*IoSession, io.Reader) (Message, error)
}

type ProtocolEncoder interface {
	Encode(*IoSession, Message) ([]byte, error)
}

type Protocol interface {
	ProtocolEncoder
	ProtocolDecoder
}
