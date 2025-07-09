package datatype

import (
	"fmt"

	"github.com/emicklei/proto"

	"github.com/ignite/cli/v29/ignite/pkg/multiformatname"
	"github.com/ignite/cli/v29/ignite/pkg/protoanalysis/protoutil"
)

// DataBytes is a string data type definition.
var DataBytes = DataType{
	Name:                    Bytes,
	DataType:                func(string) string { return "[]byte" },
	CollectionsKeyValueName: func(string) string { return "collections.BytesKey" },
	DefaultTestValue:        "3,2,3,5",
	ValueLoop:               "[]byte{1+i%1, 2+i%2, 3+i%3}",
	ProtoType: func(_, name string, index int) string {
		return fmt.Sprintf("bytes %s = %d", name, index)
	},
	GenesisArgs: func(name multiformatname.Name, value int) string {
		return fmt.Sprintf("%s: []byte(\"%d\"),\n", name.UpperCamel, value)
	},
	CLIArgs: func(name multiformatname.Name, _, prefix string, argIndex int) string {
		return fmt.Sprintf("%s%s := []byte(args[%d])", prefix, name.UpperCamel, argIndex)
	},
	ToBytes: func(name string) string {
		return name
	},
	ToString: func(name string) string {
		return fmt.Sprintf("string(%s)", name)
	},
	ToProtoField: func(_, name string, index int) *proto.NormalField {
		return protoutil.NewField(name, "bytes", index)
	},
	NonIndex: true,
}
