package add

// Options ...
type Options struct {
	AppName string
	Feature string
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	return nil
}
