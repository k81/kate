package log

import (
	"net/url"

	"go.uber.org/zap"
)

func NewSink(u *url.URL) (zap.Sink, error) {
	return NewWriter(u.Path)
}

func init() {
	zap.RegisterSink("rfile", NewSink)
}
