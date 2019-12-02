package kate

import (
	"context"
	"net/http"

	"go.uber.org/zap"
)

// Router defines the standard http outer
type Router struct {
	*http.ServeMux
	maxBodyBytes int64
	ctx          context.Context
	logger       *zap.Logger
}

// NewRouter create a http router
func NewRouter(ctx context.Context, logger *zap.Logger) *Router {
	return &Router{
		ServeMux: http.NewServeMux(),
		ctx:      ctx,
		logger:   logger,
	}
}

// SetMaxBodyBytes set the body size limit
func (r *Router) SetMaxBodyBytes(n int64) {
	r.maxBodyBytes = n
}

// Handle register a http handler for the specified path
func (r *Router) Handle(pattern string, h ContextHandler) {
	r.ServeMux.Handle(pattern, StdHandler(r.ctx, h, r.maxBodyBytes, r.logger))
}

// HandleFunc register a http handler for the specified path
func (r *Router) HandleFunc(pattern string, h func(context.Context, ResponseWriter, *Request)) {
	r.Handle(pattern, ContextHandlerFunc(h))
}

// HEAD register a handler for HEAD request
func (r *Router) HEAD(pattern string, h ContextHandler) {
	r.ServeMux.Handle(pattern, StdHandler(r.ctx, HEAD(h), r.maxBodyBytes, r.logger))
}

// OPTIONS register a handler for OPTIONS request
func (r *Router) OPTIONS(pattern string, h ContextHandler) {
	r.ServeMux.Handle(pattern, StdHandler(r.ctx, OPTIONS(h), r.maxBodyBytes, r.logger))
}

// GET register a handler for GET request
func (r *Router) GET(pattern string, h ContextHandler) {
	r.ServeMux.Handle(pattern, StdHandler(r.ctx, GET(h), r.maxBodyBytes, r.logger))
}

// POST register a handler for POST request
func (r *Router) POST(pattern string, h ContextHandler) {
	r.ServeMux.Handle(pattern, StdHandler(r.ctx, POST(h), r.maxBodyBytes, r.logger))
}

// PUT register a handler for PUT request
func (r *Router) PUT(pattern string, h ContextHandler) {
	r.ServeMux.Handle(pattern, StdHandler(r.ctx, PUT(h), r.maxBodyBytes, r.logger))
}

// DELETE register a handler for DELETE request
func (r *Router) DELETE(pattern string, h ContextHandler) {
	r.ServeMux.Handle(pattern, StdHandler(r.ctx, DELETE(h), r.maxBodyBytes, r.logger))
}

// PATCH register a handler for PATCH request
func (r *Router) PATCH(pattern string, h ContextHandler) {
	r.ServeMux.Handle(pattern, StdHandler(r.ctx, PATCH(h), r.maxBodyBytes, r.logger))
}
