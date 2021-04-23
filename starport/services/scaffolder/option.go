package scaffolder

import (
	"strings"
)

// Option configures scaffolding.
type Option func(*scaffoldingOptions)

// scaffoldingOptions keeps set of options to apply scaffolding.
type scaffoldingOptions struct {
	addressPrefix string
}

func newOptions(options ...Option) *scaffoldingOptions {
	opts := &scaffoldingOptions{}
	opts.apply(options...)
	return opts
}

func (s *scaffoldingOptions) apply(options ...Option) {
	for _, o := range options {
		o(s)
	}
}

// AddressPrefix configures address prefix for the app.
func AddressPrefix(prefix string) Option {
	return func(o *scaffoldingOptions) {
		o.addressPrefix = strings.ToLower(prefix)
	}
}
