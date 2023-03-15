package recovery

import (
	"context"
)

var (
	defaultOptions = &options{
		recoveryHandlerFunc: nil,
	}
)

// HandlerFunc is a function that recovers from the panic `p` by returning an `error`.
type HandlerFunc func(p interface{}) (err error)

// HandlerFuncContext is a function that recovers from the panic `p` by returning an `error`.
// The context can be used to extract request scoped metadata and context values.
type HandlerFuncContext func(ctx context.Context, p interface{}) (err error)

type options struct {
	recoveryHandlerFunc HandlerFuncContext
}

func evaluateOptions(opts []Option) *options {
	optCopy := &options{}
	*optCopy = *defaultOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

// Option return UnaryServerInterceptor and StreamServerInterceptor with Option.
type Option func(*options)

// WithRecoveryHandler customizes the function for recovering from a panic.
func WithRecoveryHandler(f HandlerFunc) Option {
	return func(o *options) {
		o.recoveryHandlerFunc = func(ctx context.Context, p interface{}) error {
			return f(p)
		}
	}
}

// WithRecoveryHandlerContext customizes the function for recovering from a panic.
func WithRecoveryHandlerContext(f HandlerFuncContext) Option {
	return func(o *options) {
		o.recoveryHandlerFunc = f
	}
}
