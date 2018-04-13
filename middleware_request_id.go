package kate

import (
	"context"
	"sync/atomic"

	"github.com/k81/kate/log"
)

var (
	gReqId = uint64(0)
)

func nextReqId() uint64 {
	return atomic.AddUint64(&gReqId, uint64(1))
}

func RequestId(h ContextHandler) ContextHandler {
	f := func(ctx context.Context, w ResponseWriter, r *Request) {
		ctx = log.SetContext(ctx, "session", nextReqId())
		h.ServeHTTP(ctx, w, r)
	}
	return ContextHandlerFunc(f)
}
