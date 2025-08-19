package query

import (
	"path/filepath"

	"github.com/ignite/cli/v29/ignite/pkg/multiformatname"
	"github.com/ignite/cli/v29/ignite/templates/field"
)

// Options ...
type Options struct {
	AppName     string
	ProtoDir    string
	ProtoVer    string
	ModuleName  string
	ModulePath  string
	QueryName   multiformatname.Name
	Description string
	ResFields   field.Fields
	ReqFields   field.Fields
	Paginated   bool
}

// ProtoFile returns the path to the proto folder.
func (opts *Options) ProtoFile(fname string) string {
	return filepath.Join(opts.ProtoDir, opts.AppName, opts.ModuleName, opts.ProtoVer, fname)
}
