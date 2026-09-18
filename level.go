package zlog

import (
	"expvar"
	"flag"
	"log/slog"
	"os"
)

// Leveler allows the internal singleton instance to export
// as an [expvar] published as "logLevel".
// It is a [flag.Getter] to be usable as a cli option variable,
// as well as a [slog.Leveler].
//
//	flag.Var(zlog.Level(), "loglevel", "sets the log level")
//	flag.Parse()
//
//	slog.SetDefault(
//		slog.NewTextHandler(
//			os.Stderr,
//			&slog.HandlerOptions{
//				Level: zlog.Level(),
//			},
//		),
//	)
type Leveler struct {
	slog.LevelVar
}

var _ flag.Getter = &Leveler{}
var _ slog.Leveler = &Leveler{}
var _ expvar.Var = &Leveler{}

var logLevel Leveler

func init() {
	expvar.Publish("logLevel", &logLevel)
	logLevel.Set(os.Getenv("LOG_LEVEL"))
}

// String returns a JSON safe string of the level, parseable by [Leveler.Set]
func (l *Leveler) String() string {
	out, _ := l.MarshalText()
	return string(out)
}

// Get returns the [slog.Level], implementing [flag.Getter]
func (l *Leveler) Get() any {
	return l.Level()
}

// Set parses and sets the [slog.Level], implementing [flag.Value]
func (l *Leveler) Set(value string) error {
	return l.UnmarshalText([]byte(value))
}

// Level retuns the singleton [Leveler] instance.
func Level() *Leveler {
	return &logLevel
}
