package typed

import (
	"fmt"
	"path/filepath"

	"github.com/emicklei/proto"

	"github.com/ignite/cli/ignite/pkg/multiformatname"
	"github.com/ignite/cli/ignite/pkg/protoanalysis/protoutil"
	"github.com/ignite/cli/ignite/templates/field"
)

// Options ...
type Options struct {
	AppName      string
	AppPath      string
	ModuleName   string
	ModulePath   string
	TypeName     multiformatname.Name
	MsgSigner    multiformatname.Name
	Fields       field.Fields
	Indexes      field.Fields
	NoMessage    bool
	NoSimulation bool
	IsIBC        bool
}

// Validate that options are usable.
func (opts *Options) Validate() error {
	return nil
}

// ProtoPath returns the path to the proto folder within the generated app.
func (opts *Options) ProtoPath(fname string) string {
	return filepath.Join(opts.AppPath, "proto", opts.AppName, opts.ModuleName, fname)
}

// ProtoTypeImport Return the protobuf import statement for this type.
func (opts *Options) ProtoTypeImport() *proto.Import {
	return protoutil.NewImport(fmt.Sprintf("%s/%s/%s.proto", opts.AppName, opts.ModuleName, opts.TypeName.Snake))
}
