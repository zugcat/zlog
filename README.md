# Package [github.com/zugcat/zlog](https://pkg.go.dev/github.com/zugcat/zlog?tab=doc)

```go
import github.com/zugcat/zlog
```

Package zlog provides context based logging that works with any
[slog.Handler]. It enables two main things:
  - context based logging
  - context propagation support for temporal

Instead of passing a logger instance around which was created via With,
WithAttrs, or WithGroup, one may simply:

    ctx = zlog.With(ctx, "operation", operation)
    slog.InfoContext(ctx, "log message")

This allows any called function anywhere in the call stack to inherit the
logging attributes.

If supplied as a ContextPropagators option to temporal [client.Options],
it allows this same mechanism to cross temporal workflow boundaries into the
activities:

    cli, err := client.DialContext(
    	ctx,
    	client.Options{
    		...
    		ContextPropagators: []workflow.ContextPropagator{zlog.NewPropagator(), ...},
    	}
    }

## Functions
### Func Handler
```go
func Handler(h slog.Handler) slog.Handler
```
Handler returns an instance of our handler

### Func New
```go
func New(handler slog.Handler) *slog.Logger
```
New returns a new logger for the given handler. A zlog style handler will be
wrapped around handler.

### Func NewDefault
```go
func NewDefault(w io.Writer) *slog.Logger
```
NewDefault returns a new logger with a default handler.

### Func NewDefaultHandler
```go
func NewDefaultHandler(w io.Writer) slog.Handler
```
NewDefaultHandler returns new json handler with presumtuous options

### Func NewPropagator
```go
func NewPropagator() workflow.ContextPropagator
```
NewPropagator creates a new [workflow.ContextPropagator] for zlog

### Func With
```go
func With(ctx context.Context, args ...any) context.Context
```
With acts a per [slog.Logger.With] in the returned context

### Func WithAttrs
```go
func WithAttrs(ctx context.Context, attrs []slog.Attr) context.Context
```
WithAttrs acts per [With] except taking a slice of [slog.Attr]

### Func WithGroup
```go
func WithGroup(ctx context.Context, group string) context.Context
```
WithGroup acts as per [slog.Logger.WithGroup] in the returned context



## Types
### Type Leveler
```go
type Leveler struct {
	slog.LevelVar
}
```
Leveler allows the internal singleton instance to export as an expvar
published as "logLevel". It is a flag.Getter to be usable as a cli option
variable, as well as a [slog.Leveler].

    flag.Var(zlog.Level(), "loglevel", "sets the log level")
    flag.Parse()

    slog.SetDefault(
    	slog.NewTextHandler(
    		os.Stderr,
    		&slog.HandlerOptions{
    			Level: zlog.Level(),
    		},
    	),
    )

### Functions

```go
func Level() *Leveler
```
Level retuns the singleton [Leveler] instance.



### Methods

```go
func (l *Leveler) Get() any
```
Get returns the [slog.Level], implementing flag.Getter


```go
func (l *Leveler) Set(value string) error
```
Set parses and sets the [slog.Level], implementing flag.Value


```go
func (l *Leveler) String() string
```
String returns a JSON safe string of the level, parseable by [Leveler.Set]







