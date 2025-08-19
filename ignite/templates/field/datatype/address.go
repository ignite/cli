package datatype

import (
	"fmt"

	"github.com/emicklei/proto"

	"github.com/ignite/cli/v29/ignite/pkg/multiformatname"
	"github.com/ignite/cli/v29/ignite/pkg/protoanalysis/protoutil"
)

// DataAddress address (string) data type definition.
var DataAddress = DataType{
	Name:                    Address,
	DataType:                func(string) string { return "string" },
	CollectionsKeyValueName: func(string) string { return "collections.StringKey" },
	DefaultTestValue:        "cosmos1abcdefghijklmnopqrstuvwxyz0123456",
	ValueLoop:               "fmt.Sprintf(`cosmos1abcdef%d`, i)",
	ValueIndex:              "`cosmos1abcdefghijklmnopqrstuvwxyz0123456`",
	ValueInvalidIndex:       "`cosmos1invalid`",
	ProtoType: func(_, name string, index int) string {
		return fmt.Sprintf("string %s = %d", name, index)
	},
	GenesisArgs: func(name multiformatname.Name, value int) string {
		return fmt.Sprintf("%s: \"%d\",\n", name.UpperCamel, value)
	},
	CLIArgs: func(name multiformatname.Name, _, prefix string, argIndex int) string {
		return fmt.Sprintf("%s%s := args[%d]", prefix, name.UpperCamel, argIndex)
	},
	ToBytes: func(name string) string {
		return fmt.Sprintf("%[1]vBytes := []byte(%[1]v)", name)
	},
	ToString: func(name string) string {
		return name
	},
	ToProtoField: func(_, name string, index int) *proto.NormalField {
		field := protoutil.NewField(name, "string", index)
		option := protoutil.NewOption("cosmos_proto.scalar", "cosmos.AddressString", protoutil.Custom())
		field.Options = append(field.Options, option)
		return field
	},
	ProtoImports: []string{"cosmos_proto/cosmos.proto"},
}
