package main

import (
	"context"

	"github.com/k81/kate/log"
)

func foo(ctx context.Context) {
	ctx = log.SetContext(ctx, "context", "foo")

	log.Info(ctx, "foo called")
}

func bar(ctx context.Context) {
	ctx = log.SetContext(ctx, "context", "bar")

	log.Info(ctx, "bar called")
}

func main() {
	ctx := log.SetContext(context.Background(), "module", "example")
	log.Info(ctx, "program started")
	foo(ctx)
	bar(ctx)
	log.Info(ctx, "program exited")
}
