package nrgrpc

import (
	"strings"

	"google.golang.org/grpc/codes"
)

// Options contains some configurations for an interceptor
type Options struct {
	IgnoredServices map[string]struct{}
	IgnoredMethods  map[string]struct{}
	IgnoredCodes    map[codes.Code]struct{}
}

func composeOptions(funcs []Option) Options {
	o := Options{
		IgnoredServices: map[string]struct{}{},
		IgnoredMethods:  map[string]struct{}{},
		IgnoredCodes:    map[codes.Code]struct{}{},
	}
	if len(funcs) > 0 && funcs[0] != nil {
		for _, f := range funcs {
			o = f(o)
		}
	}
	return o
}

// IsIgnored returned true if the given method is ignored
func (o *Options) IsIgnored(fullMethod string) bool {
	if _, ok := o.IgnoredMethods[fullMethod]; ok {
		return true
	}
	if len(o.IgnoredServices) == 0 {
		return false
	}
	strs := strings.Split(fullMethod, "/")
	_, ok := o.IgnoredServices[strs[0]]
	if ok {
		o.IgnoredMethods[fullMethod] = struct{}{}
	}
	return ok
}

// IsCodeIgnored returned true if the given code is ignored
func (o *Options) IsCodeIgnored(c codes.Code) bool {
	_, ok := o.IgnoredCodes[c]
	return ok
}

// Option is a function for building configurations object for an interceptor
type Option func(Options) Options

// WithIgnoredServices receives service names to get newrelic to ignore them
func WithIgnoredServices(services ...string) Option {
	return func(o Options) Options {
		for _, s := range services {
			o.IgnoredServices[s] = struct{}{}
		}
		return o
	}
}

// WithIgnoredMethods receives full method names to get newrelic to ignore them
func WithIgnoredMethods(fullMethods ...string) Option {
	return func(o Options) Options {
		for _, m := range fullMethods {
			o.IgnoredMethods[m] = struct{}{}
		}
		return o
	}
}

// WithIgnoredCodes receives codes to get newrelic to ignore them
func WithIgnoredCodes(cs ...codes.Code) Option {
	return func(o Options) Options {
		for _, c := range cs {
			o.IgnoredCodes[c] = struct{}{}
		}
		return o
	}
}
