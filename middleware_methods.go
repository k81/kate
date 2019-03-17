package kate

import (
	"context"
	"net/http"
	"strings"
)

// MethodOnly is a middleware to restrict http method for standard http router
func MethodOnly(method string, h ContextHandler) ContextHandler {
	f := func(ctx context.Context, w ResponseWriter, r *Request) {
		if strings.ToUpper(r.Method) != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			// nolint:errcheck
			w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
			return
		}

		h.ServeHTTP(ctx, w, r)
	}
	return ContextHandlerFunc(f)
}

// HEAD only allow HEAD method
func HEAD(h ContextHandler) ContextHandler {
	return MethodOnly("HEAD", h)
}

// OPTIONS only allow OPTIONS method
func OPTIONS(h ContextHandler) ContextHandler {
	return MethodOnly("OPTIONS", h)
}

// GET only allow GET method
func GET(h ContextHandler) ContextHandler {
	return MethodOnly("GET", h)
}

// POST only allow POST method
func POST(h ContextHandler) ContextHandler {
	return MethodOnly("POST", h)
}

// PUT only allow PUT method
func PUT(h ContextHandler) ContextHandler {
	return MethodOnly("PUT", h)
}

// DELETE only allow DELETE method
func DELETE(h ContextHandler) ContextHandler {
	return MethodOnly("DELETE", h)
}

// PATCH only allow PATCH method
func PATCH(h ContextHandler) ContextHandler {
	return MethodOnly("PATCH", h)
}
