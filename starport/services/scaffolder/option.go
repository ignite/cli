package scaffolder

import (
	"strings"

	"github.com/tendermint/starport/starport/pkg/cosmosver"
)

// Option configures scaffolding.
type Option func(*scaffoldingOptions)

// scaffoldingOptions keeps set of options to apply scaffolding.
type scaffoldingOptions struct {
	addressPrefix string
	sdkVersion    cosmosver.MajorVersion
}

func newOptions(options ...Option) *scaffoldingOptions {
	opts := &scaffoldingOptions{
		sdkVersion: cosmosver.Launchpad,
	}
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

// SdkVersion specifies Cosmos-SDK version.
func SdkVersion(v cosmosver.MajorVersion) Option {
	return func(o *scaffoldingOptions) {
		o.sdkVersion = v
	}
}
