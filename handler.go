package zlog

import (
	"context"
	"io"
	"log/slog"
	"slices"
)

type handler struct {
	out   slog.Handler
	up    *handler
	attrs []slog.Attr
	group string
}

var _ slog.Handler = &handler{}

// NewDefaultHandler returns new json handler with presumtuous options
func NewDefaultHandler(w io.Writer) slog.Handler {
	return &handler{out: slog.NewJSONHandler(
		w, &slog.HandlerOptions{
			AddSource: true,
			Level:     Level(),
		},
	)}
}

// Handler returns an instance of our handler
func Handler(h slog.Handler) slog.Handler {
	if h, ok := h.(*handler); ok {
		return &handler{
			out: h.out,
			up:  h,
		}
	}

	return &handler{out: h}
}

// WithAttrs implements [slog.Handler]
func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &handler{
		out:   h.out.WithAttrs(attrs),
		up:    h,
		attrs: slices.Clone(attrs),
	}
}

// WithGroup implements [slog.Handler]
func (h *handler) WithGroup(name string) slog.Handler {
	return &handler{
		out:   h.out.WithGroup(name),
		up:    h,
		group: name,
	}
}

// Enabled implements [slog.Handler]
func (h *handler) Enabled(ctx context.Context, level slog.Level) bool {
	if hh, ok := ctx.Value(handlerValue).(*handler); ok {
		h = hh
	}
	return h.out.Enabled(ctx, level)
}

// Handle implements [slog.Handler]
func (h *handler) Handle(ctx context.Context, record slog.Record) error {
	if hh, ok := ctx.Value(handlerValue).(*handler); ok {
		h = hh
	}
	return h.out.Handle(ctx, record)
}
