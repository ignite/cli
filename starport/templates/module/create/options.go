package modulecreate

// CreateOptions ...
type CreateOptions struct {
	ModuleName string
	ModulePath string
	AppName    string
	OwnerName  string

	// True if the module should implement the IBC module interface
	IsIBC bool

	// Channel ordering of the IBC module: ordered, unordered or none
	IBCOrdering string
}

// Validate that options are usable
func (opts *CreateOptions) Validate() error {
	return nil
}
