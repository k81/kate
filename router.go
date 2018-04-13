package kate

import (
	"context"
	"net/http"
)

type Router struct {
	*http.ServeMux
	maxBodyBytes int64
	ctx          context.Context
}

func NewRouter(ctx context.Context) *Router {
	return &Router{
		ServeMux: http.NewServeMux(),
		ctx:      ctx,
	}
}

func (r *Router) SetMaxBodyBytes(n int64) {
	r.maxBodyBytes = n
}

func (r *Router) Handle(pattern string, h ContextHandler) {
	r.ServeMux.Handle(pattern, StdHandler(r.ctx, h, r.maxBodyBytes))
}

func (r *Router) HandleFunc(pattern string, h func(context.Context, ResponseWriter, *Request)) {
	r.Handle(pattern, ContextHandlerFunc(h))
}

func (r *Router) HEAD(pattern string, h ContextHandler) {
	r.ServeMux.Handle(pattern, StdHandler(r.ctx, HEAD(h), r.maxBodyBytes))
}

func (r *Router) OPTIONS(pattern string, h ContextHandler) {
	r.ServeMux.Handle(pattern, StdHandler(r.ctx, OPTIONS(h), r.maxBodyBytes))
}

func (r *Router) GET(pattern string, h ContextHandler) {
	r.ServeMux.Handle(pattern, StdHandler(r.ctx, GET(h), r.maxBodyBytes))
}

func (r *Router) POST(pattern string, h ContextHandler) {
	r.ServeMux.Handle(pattern, StdHandler(r.ctx, POST(h), r.maxBodyBytes))
}

func (r *Router) PUT(pattern string, h ContextHandler) {
	r.ServeMux.Handle(pattern, StdHandler(r.ctx, PUT(h), r.maxBodyBytes))
}

func (r *Router) DELETE(pattern string, h ContextHandler) {
	r.ServeMux.Handle(pattern, StdHandler(r.ctx, DELETE(h), r.maxBodyBytes))
}

func (r *Router) PATCH(pattern string, h ContextHandler) {
	r.ServeMux.Handle(pattern, StdHandler(r.ctx, PATCH(h), r.maxBodyBytes))
}
