package module

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

// ImportOptions ...
type ImportOptions struct {
	AppName string
	Feature string
	BinaryNamePrefix string
}

// Validate that options are usable
func (opts *ImportOptions) Validate() error {
	return nil
}
