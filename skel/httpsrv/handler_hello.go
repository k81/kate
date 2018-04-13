package httpsrv

import (
	"context"

	"github.com/k81/kate"
)

type Hello struct{}

func (h *Hello) ServeHTTP(ctx context.Context, w kate.ResponseWriter, r *kate.Request) {
	kate.OkData(ctx, w, "hello world")
}
