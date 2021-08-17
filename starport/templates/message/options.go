package message

import (
	"github.com/tendermint/starport/starport/pkg/field"
	"github.com/tendermint/starport/starport/pkg/multiformatname"
)

// Options ...
type Options struct {
	AppName    string
	ModuleName string
	ModulePath string
	OwnerName  string
	MsgName    multiformatname.Name
	MsgDesc    string
	Fields     field.Fields
	ResFields  field.Fields
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	return nil
}
