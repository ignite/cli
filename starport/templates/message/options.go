package message

import (
	"github.com/tendermint/starport/starport/pkg/field"
)

// Options ...
type Options struct {
	AppName    string
	ModuleName string
	ModulePath string
	OwnerName  string
	MsgName    string
	MsgDesc    string
	Fields     []field.Field
	ResFields  []field.Field
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	return nil
}
