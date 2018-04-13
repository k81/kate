package kate

import (
	"context"
	"net/http"
	"strings"
)

func MethodOnly(method string, h ContextHandler) ContextHandler {
	f := func(ctx context.Context, w ResponseWriter, r *Request) {
		if strings.ToUpper(r.Method) != method {
			//log.Debug(ctx, "method not allowed", "method", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
			return
		}

		h.ServeHTTP(ctx, w, r)
	}
	return ContextHandlerFunc(f)
}

func HEAD(h ContextHandler) ContextHandler {
	return MethodOnly("HEAD", h)
}

func OPTIONS(h ContextHandler) ContextHandler {
	return MethodOnly("OPTIONS", h)
}

func GET(h ContextHandler) ContextHandler {
	return MethodOnly("GET", h)
}

func POST(h ContextHandler) ContextHandler {
	return MethodOnly("POST", h)
}

func PUT(h ContextHandler) ContextHandler {
	return MethodOnly("PUT", h)
}

func DELETE(h ContextHandler) ContextHandler {
	return MethodOnly("DELETE", h)
}

func PATCH(h ContextHandler) ContextHandler {
	return MethodOnly("PATCH", h)
}

// Deprecated
// for backward compability only
func PostOnly(h ContextHandler) ContextHandler {
	return POST(h)
}
