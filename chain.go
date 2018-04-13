package kate

import (
	"context"
	"fmt"
	"path"
	"reflect"
	"runtime"
	"strings"
)

type Chain struct {
	middlewares []Middleware
}

func NewChain(middlewares ...Middleware) Chain {
	recoveryFound := false

	for _, m := range middlewares {
		var (
			fn   = GetFunctionName(m)
			name = path.Base(fn)
		)

		if !recoveryFound && name == "kate.Recovery" {
			recoveryFound = true
		}

		if recoveryFound && strings.HasPrefix(name, "kate.Timeout") {
			warning := "middleware \"kate.Recovery\" must be registered after \"kate.Timeout\""
			fmt.Println(warning)
			panic(warning)
		}
	}

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

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
