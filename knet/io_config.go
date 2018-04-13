package knet

import "time"

type IoConfig struct {
	SendQueueSize int
	RecvQueueSize int
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
}
