package message

import "github.com/tendermint/starport/starport/templates/typed"

// Options ...
type Options struct {
	AppName    string
	ModuleName string
	ModulePath string
	OwnerName  string
	MsgName    string
	MsgDesc    string
	Fields     []typed.Field
	ResFields  []typed.Field
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	return nil
}
