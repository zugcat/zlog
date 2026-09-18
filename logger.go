package zlog

import (
	"io"
	"log/slog"
)

// New returns a new logger for the given handler.
// A zlog style handler will be wrapped around handler.
func New(handler slog.Handler) *slog.Logger {
	return slog.New(Handler(handler))
}

// NewDefault returns a new logger with a default handler.
func NewDefault(w io.Writer) *slog.Logger {
	return slog.New(NewDefaultHandler(w))
}
