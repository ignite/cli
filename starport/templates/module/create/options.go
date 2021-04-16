package modulecreate

// CreateOptions represents the options to scaffold a Cosmos SDK module
type CreateOptions struct {
	ModuleName  string
	ModulePath  string
	AppName     string
	OwnerName   string
	IsIBC       bool
	IBCOrdering string
}

// Validate that options are usable
func (opts *CreateOptions) Validate() error {
	return nil
}
