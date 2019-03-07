package kate

import (
	"context"
	"net/http"

	"github.com/k81/kate/utils"
	"github.com/k81/log"
)

func Recovery(h ContextHandler) ContextHandler {
	f := func(ctx context.Context, w ResponseWriter, r *Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
				log.Error(ctx, "got panic", "error", err, "stack", utils.GetPanicStack())
			}
		}()

		h.ServeHTTP(ctx, w, r)
	}
	return ContextHandlerFunc(f)
}
