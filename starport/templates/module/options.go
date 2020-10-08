package module

// Options ...
type Options struct {
	AppName string
	Feature string
}

// Validate that options are usable
func (opts *Options) Validate() error {
	return nil
}
