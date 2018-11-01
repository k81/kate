package kate

import (
	"context"
	"net/http"

	"github.com/k81/kate/log"
	"github.com/k81/kate/utils"
)

func Recovery(h ContextHandler) ContextHandler {
	f := func(ctx context.Context, w ResponseWriter, r *Request) {
		defer func() {
			if r := recover(); r != nil {
				switch v := r.(type) {
				case ErrorInfo:
					Error(ctx, w, v)
				default:
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
					log.Error(ctx, "got panic", "error", v, "stack", utils.GetPanicStack())
				}
			}
		}()

		h.ServeHTTP(ctx, w, r)
	}
	return ContextHandlerFunc(f)
}
