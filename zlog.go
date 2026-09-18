// Package zlog provides context based logging that works with any [slog.Handler].
// It enables two main things:
//   - context based logging
//   - context propagation support for temporal
//
// Instead of passing a logger instance around which was created via With, WithAttrs, or WithGroup,
// one may simply:
//
//	ctx = zlog.With(ctx, "operation", operation)
//	slog.InfoContext(ctx, "log message")
//
// This allows any called function anywhere in the call stack to inherit the logging attributes.
//
// If supplied as a ContextPropagators option to temporal [client.Options], it allows this same
// mechanism to cross temporal workflow boundaries into the activities:
//
//	cli, err := client.DialContext(
//		ctx,
//		client.Options{
//			...
//			ContextPropagators: []workflow.ContextPropagator{zlog.NewPropagator(), ...},
//		}
//	}
package zlog

//go:generate go run cloudeng.io/go/cmd/gomarkdown --overwrite .
