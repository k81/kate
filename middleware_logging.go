package kate

import (
	"context"
	"time"

	"github.com/k81/log"
)

// Logging implements the request in/out logging middleware
func Logging(h ContextHandler) ContextHandler {
	f := func(ctx context.Context, w ResponseWriter, r *Request) {
		var start = time.Now()

		log.Tag("_com_request_in").Info(ctx, "request in",
			"remote", r.RemoteAddr,
			"method", r.Method,
			"url", r.RequestURI,
			"body", string(r.RawBody))

		h.ServeHTTP(ctx, w, r)

		log.Tag("_com_request_out").Info(ctx, "request finished",
			"status_code", w.StatusCode(),
			"body", string(w.RawBody()),
			"duration_ms", int64(time.Since(start)/time.Millisecond))
	}
	return ContextHandlerFunc(f)
}
