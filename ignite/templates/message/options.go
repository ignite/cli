package message

import (
	"path/filepath"

	"github.com/ignite/cli/v29/ignite/pkg/multiformatname"
	"github.com/ignite/cli/v29/ignite/templates/field"
)

// Options ...
type Options struct {
	AppName      string
	ProtoDir     string
	ProtoVer     string
	ModuleName   string
	ModulePath   string
	MsgName      multiformatname.Name
	MsgSigner    multiformatname.Name
	MsgDesc      string
	Fields       field.Fields
	ResFields    field.Fields
	NoSimulation bool
}

// ProtoFile returns the path to the proto folder.
func (opts *Options) ProtoFile(fname string) string {
	return filepath.Join(opts.ProtoDir, opts.AppName, opts.ModuleName, opts.ProtoVer, fname)
}
