package typed

import (
	"fmt"
	"path/filepath"

	"github.com/emicklei/proto"

	"github.com/ignite/cli/v29/ignite/pkg/multiformatname"
	"github.com/ignite/cli/v29/ignite/pkg/protoanalysis/protoutil"
	"github.com/ignite/cli/v29/ignite/templates/field"
)

// Options ...
type Options struct {
	AppName      string
	ProtoDir     string
	ProtoVer     string
	ModuleName   string
	ModulePath   string
	TypeName     multiformatname.Name
	MsgSigner    multiformatname.Name
	Fields       field.Fields
	Index        field.Field
	NoMessage    bool
	NoSimulation bool
	IsIBC        bool
}

// ProtoFile returns the path to the proto folder within the generated app.
func (opts *Options) ProtoFile(fname string) string {
	return filepath.Join(opts.ProtoDir, opts.AppName, opts.ModuleName, opts.ProtoVer, fname)
}

// ProtoTypeImport Return the protobuf import statement for this type.
func (opts *Options) ProtoTypeImport() *proto.Import {
	return protoutil.NewImport(fmt.Sprintf("%s/%s/%s/%s.proto", opts.AppName, opts.ModuleName, opts.ProtoVer, opts.TypeName.Snake))
}
