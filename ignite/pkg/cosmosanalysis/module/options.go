package module

// DiscoverOption configures calls to Discovery function.
type DiscoverOption func(*discoverOptions)

type discoverOptions struct {
	protoDir, sdkDir string
}

// WithProtoDir sets the relative proto directory path.
func WithProtoDir(path string) DiscoverOption {
	return func(o *discoverOptions) {
		o.protoDir = path
	}
}

// WithSDKDir sets the absolute directory path to the Cosmos SDK Go package.
func WithSDKDir(path string) DiscoverOption {
	return func(o *discoverOptions) {
		o.sdkDir = path
	}
}
