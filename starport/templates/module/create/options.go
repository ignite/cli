package modulecreate

// CreateOptions ...
type CreateOptions struct {
	ModuleName string
	ModulePath string
	AppName    string
	OwnerName  string
}

// Validate that options are usable
func (opts *CreateOptions) Validate() error {
	return nil
}
