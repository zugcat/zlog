package zlog

import (
	"context"
	"log/slog"
)

type handlerType int

var handlerValue = handlerType(0)

// withHandler sets a handler to carry in a context
func withHandler(parent context.Context, handler slog.Handler) context.Context {
	return context.WithValue(parent, handlerValue, handler)
}

// getHandler gets the current handler from a context, if one if present,
// otherwise a new or existing handler from [slog.Default]
func getHandler(ctx context.Context) *handler {
	if handler, ok := ctx.Value(handlerValue).(*handler); ok {
		return handler
	}

	h := slog.Default().Handler()
	if h, ok := h.(*handler); ok {
		return h
	}
	return Handler(h).(*handler)
}

// WithAttrs acts per [With] except taking a slice of [slog.Attr]
func WithAttrs(ctx context.Context, attrs []slog.Attr) context.Context {
	return withHandler(ctx, getHandler(ctx).WithAttrs(attrs))
}

// WithGroup acts as per [slog.Logger.WithGroup] in the returned context
func WithGroup(ctx context.Context, group string) context.Context {
	return withHandler(ctx, getHandler(ctx).WithGroup(group))
}

// With acts a per [slog.Logger.With] in the returned context
func With(ctx context.Context, args ...any) context.Context {
	if len(args) == 0 {
		return ctx
	}

	attrs := make([]slog.Attr, 0, len(args)>>1)
	for len(args) > 0 {
		switch x := args[0].(type) {
		case string:
			if len(args) == 1 {
				attrs = append(attrs, slog.String("!badvalue", x))
				args = nil
			} else {
				attrs = append(attrs, slog.Any(x, args[1]))
				args = args[2:]
			}

		case slog.Attr:
			attrs = append(attrs, x)
			args = args[1:]

		default:
			attrs = append(attrs, slog.Any("!badkey", x))
			args = args[1:]
		}
	}

	return WithAttrs(ctx, attrs)
}
