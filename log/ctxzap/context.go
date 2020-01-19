package ctxzap

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ctxMarker struct{}

var (
	ctxMarkerKey = &ctxMarker{}
	nullLogger   = zap.NewNop()
)

// Extract takes the call-scoped Logger from grpc_zap middleware.
//
// It always returns a Logger that has all the grpc_ctxtags updated.
func Extract(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(ctxMarkerKey).(*zap.Logger); ok {
		return logger
	}

	return nullLogger
}

// ToContext adds the zap.Logger to the context for extraction later.
// Returning the new context that has been created.
func ToContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxMarkerKey, logger)
}

// With append logging fields to context
func With(ctx context.Context, fields ...zapcore.Field) context.Context {
	return ToContext(ctx, Extract(ctx).With(fields...))
}
