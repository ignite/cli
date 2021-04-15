package modulecreate

// CreateOptions defines options to create a module
type CreateOptions struct {
	ModuleName  string
	ModulePath  string
	AppName     string
	OwnerName   string
	IsIBC       bool
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
