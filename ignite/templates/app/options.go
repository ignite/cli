package app

// Options ...
type Options struct {
	AppName          string
	AppPath          string
	ProtoPath        string
	GitHubPath       string
	BinaryNamePrefix string
	ModulePath       string
	AddressPrefix    string
	// IncludePrefixes is used to filter the files to include from the generator
	IncludePrefixes []string
	IsChainMinimal  bool
	IsConsumerChain bool
}

// Validate that options are usable.
func (opts *Options) Validate() error {
	return nil
}
