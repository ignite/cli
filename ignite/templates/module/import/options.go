package moduleimport

// ImportOptions ...
type ImportOptions struct {
	AppName          string
	AppPath          string
	Feature          string
	BinaryNamePrefix string
}

// Validate that options are usable.
func (opts *ImportOptions) Validate() error {
	return nil
}
