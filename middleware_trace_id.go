package kate

import (
	"context"

	"github.com/k81/kate/log"
	"github.com/k81/kate/utils"
)

func TraceId(h ContextHandler) ContextHandler {
	f := func(ctx context.Context, w ResponseWriter, r *Request) {
		traceId := r.Header.Get("X-Trace-Id")
		if traceId == "" {
			traceId = utils.FastUUIDStr()
		}
		ctx = context.WithValue(ctx, "X-Trace-Id", traceId)
		ctx = log.SetContext(ctx, "trace_id", traceId)
		h.ServeHTTP(ctx, w, r)
	}
	return ContextHandlerFunc(f)
}
