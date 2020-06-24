package app

// Options ...
type Options struct {
	AppName    string
	ModulePath string
	Denom      string
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	return nil
}
