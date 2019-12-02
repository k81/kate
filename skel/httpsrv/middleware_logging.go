package httpsrv

import (
	"context"
	"time"

	"github.com/k81/kate"
	"go.uber.org/zap"
)

// Logging implements the request in/out logging middleware
func Logging(h kate.ContextHandler) kate.ContextHandler {
	f := func(ctx context.Context, w kate.ResponseWriter, r *kate.Request) {
		var start = time.Now()

		logger.Info("request in",
			zap.String("remote", r.RemoteAddr),
			zap.String("method", r.Method),
			zap.String("url", r.RequestURI),
			zap.String("body", string(r.RawBody)))

		h.ServeHTTP(ctx, w, r)

		logger.Info("request finished",
			zap.Int("status_code", w.StatusCode()),
			zap.String("body", string(w.RawBody())),
			zap.Int64("duration_ms", int64(time.Since(start)/time.Millisecond)))
	}
	return kate.ContextHandlerFunc(f)
}
