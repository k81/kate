package httpsrv

import (
	"context"

	"github.com/k81/kate"
)

type HelloHandler struct {
	BaseHandler
}

func (h *HelloHandler) ServeHTTP(ctx context.Context, w kate.ResponseWriter, r *kate.Request) {
	h.OKData(ctx, w, "hello world")
}
