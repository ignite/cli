package message

import (
	"github.com/tendermint/starport/starport/pkg/multiformatname"
	"github.com/tendermint/starport/starport/templates/field"
)

// Options ...
type Options struct {
	AppName    string
	AppPath    string
	ModuleName string
	ModulePath string
	OwnerName  string
	MsgName    multiformatname.Name
	MsgSigner  multiformatname.Name
	MsgDesc    string
	Fields     field.Fields
	ResFields  field.Fields
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	return nil
}
