package kate

import (
	"context"

	"github.com/julienschmidt/httprouter"
)

type RESTRouter struct {
	*httprouter.Router
	maxBodyBytes int64
	ctx          context.Context
}

func NewRESTRouter(ctx context.Context) *RESTRouter {
	r := &RESTRouter{
		Router: httprouter.New(),
		ctx:    ctx,
	}
	r.Router.RedirectTrailingSlash = false
	r.Router.RedirectFixedPath = false
	return r
}

func (r *RESTRouter) SetMaxBodyBytes(n int64) {
	r.maxBodyBytes = n
}

func (r *RESTRouter) Handle(method, pattern string, h ContextHandler) {
	r.Router.Handle(method, pattern, Handle(r.ctx, h, r.maxBodyBytes))
}

func (r *RESTRouter) HandleFunc(method, pattern string, h func(context.Context, ResponseWriter, *Request)) {
	r.Handle(method, pattern, ContextHandlerFunc(h))
}

func (r *RESTRouter) HEAD(pattern string, h ContextHandler) {
	r.Handle("HEAD", pattern, h)
}

func (r *RESTRouter) OPTIONS(pattern string, h ContextHandler) {
	r.Handle("OPTIONS", pattern, h)
}

func (r *RESTRouter) GET(pattern string, h ContextHandler) {
	r.Handle("GET", pattern, h)
}

func (r *RESTRouter) POST(pattern string, h ContextHandler) {
	r.Handle("POST", pattern, h)
}

func (r *RESTRouter) PUT(pattern string, h ContextHandler) {
	r.Handle("PUT", pattern, h)
}

func (r *RESTRouter) DELETE(pattern string, h ContextHandler) {
	r.Handle("DELETE", pattern, h)
}

func (r *RESTRouter) PATCH(pattern string, h ContextHandler) {
	r.Handle("PATCH", pattern, h)
}
