package modulecreate

// CreateOptions represents the options to scaffold a Cosmos SDK module
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

// CreateOptions defines options to add MsgServer
type MsgServerOptions struct {
	ModuleName string
	ModulePath string
	AppName    string
	OwnerName  string
}

// Validate that options are usable
func (opts *CreateOptions) Validate() error {
	return nil
}
