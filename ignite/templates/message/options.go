package message

import (
	"github.com/ignite-hq/cli/ignite/pkg/multiformatname"
	"github.com/ignite-hq/cli/ignite/templates/field"
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
