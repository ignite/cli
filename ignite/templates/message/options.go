package message

import (
	"github.com/ignite/cli/v29/ignite/pkg/multiformatname"
	"github.com/ignite/cli/v29/ignite/templates/field"
)

// Options ...
type Options struct {
	AppName      string
	AppPath      string
	ProtoPath    string
	ModuleName   string
	ModulePath   string
	MsgName      multiformatname.Name
	MsgSigner    multiformatname.Name
	MsgDesc      string
	Fields       field.Fields
	ResFields    field.Fields
	NoSimulation bool
}

// Validate that options are usable.
func (opts *Options) Validate() error {
	return nil
}
