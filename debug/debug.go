package debug

import "context"

type ctxMarker struct{}

var ctxMarkerKey = &ctxMarker{}

func Get(ctx context.Context) bool {
	if debug, ok := ctx.Value(ctxMarkerKey).(bool); ok {
		return debug
	}

	return false
}

func Wrap(ctx context.Context, enabled bool) context.Context {
	return context.WithValue(ctx, ctxMarkerKey, enabled)
}
