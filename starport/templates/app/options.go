package app

// Options ...
type Options struct {
	AppName          string
	OwnerName        string
	OwnerAndRepoName string
	BinaryNamePrefix string
	ModulePath       string
	AddressPrefix    string
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	return nil
}
