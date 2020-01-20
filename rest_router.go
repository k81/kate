package kate

import (
	"context"

	"github.com/julienschmidt/httprouter"
	"github.com/k81/kate/log/ctxzap"
	"go.uber.org/zap"
)

// RESTRouter define the REST router
type RESTRouter struct {
	*httprouter.Router
	maxBodyBytes int64
	ctx          context.Context
}

// NewRESTRouter create a REST router
func NewRESTRouter(ctx context.Context, logger *zap.Logger) *RESTRouter {
	r := &RESTRouter{
		Router: httprouter.New(),
		ctx:    ctxzap.ToContext(ctx, logger),
	}
	r.Router.RedirectTrailingSlash = false
	r.Router.RedirectFixedPath = false
	return r
}

// SetMaxBodyBytes set the body size limit
func (r *RESTRouter) SetMaxBodyBytes(n int64) {
	r.maxBodyBytes = n
}

// Handle register a http handler for the specified method and path
func (r *RESTRouter) Handle(method, pattern string, h ContextHandler) {
	r.Router.Handle(method, pattern, Handle(r.ctx, h, r.maxBodyBytes))
}

// HandleFunc register a http handler for the specified method and path
func (r *RESTRouter) HandleFunc(method, pattern string, h func(context.Context, ResponseWriter, *Request)) {
	r.Handle(method, pattern, ContextHandlerFunc(h))
}

// HEAD register a handler for HEAD request
func (r *RESTRouter) HEAD(pattern string, h ContextHandler) {
	r.Handle("HEAD", pattern, h)
}

// OPTIONS register a handler for OPTIONS request
func (r *RESTRouter) OPTIONS(pattern string, h ContextHandler) {
	r.Handle("OPTIONS", pattern, h)
}

// GET register a handler for GET request
func (r *RESTRouter) GET(pattern string, h ContextHandler) {
	r.Handle("GET", pattern, h)
}

// POST register a handler for POST request
func (r *RESTRouter) POST(pattern string, h ContextHandler) {
	r.Handle("POST", pattern, h)
}

// PUT register a handler for PUT request
func (r *RESTRouter) PUT(pattern string, h ContextHandler) {
	r.Handle("PUT", pattern, h)
}

// DELETE register a handler for DELETE request
func (r *RESTRouter) DELETE(pattern string, h ContextHandler) {
	r.Handle("DELETE", pattern, h)
}

// PATCH register a handler for PATCH request
func (r *RESTRouter) PATCH(pattern string, h ContextHandler) {
	r.Handle("PATCH", pattern, h)
}
