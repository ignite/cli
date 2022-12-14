package app

// Options ...
type Options struct {
	AppName          string
	AppPath          string
	GitHubPath       string
	BinaryNamePrefix string
	ModulePath       string
	AddressPrefix    string
}

// Validate that options are usable.
func (opts *Options) Validate() error {
	return nil
}
