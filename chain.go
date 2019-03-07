package kate

import (
	"context"
)

type Chain struct {
	middlewares []Middleware
}

func NewChain(middlewares ...Middleware) Chain {
	c := Chain{}
	c.middlewares = append(c.middlewares, middlewares...)

	return c
}

func (c Chain) Then(h ContextHandler) ContextHandler {
	if h == nil {
		panic("handler == nil")
	}

	final := h

	for i := len(c.middlewares) - 1; i >= 0; i-- {
		final = c.middlewares[i](final)
	}

	return final
}

func (c Chain) ThenFunc(h func(context.Context, ResponseWriter, *Request)) ContextHandler {
	return c.Then(ContextHandlerFunc(h))
}

func (c Chain) Append(middlewares ...Middleware) Chain {
	newMws := make([]Middleware, len(c.middlewares)+len(middlewares))
	copy(newMws, c.middlewares)
	copy(newMws[len(c.middlewares):], middlewares)

	newChain := NewChain(newMws...)
	return newChain
}
