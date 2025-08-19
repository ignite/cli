package datatype

import (
	"fmt"

	"github.com/emicklei/proto"

	"github.com/ignite/cli/v29/ignite/pkg/multiformatname"
	"github.com/ignite/cli/v29/ignite/pkg/protoanalysis/protoutil"
)

// DataBool bool data type definition.
var DataBool = DataType{
	Name:                    Bool,
	DataType:                func(string) string { return "bool" },
	CollectionsKeyValueName: func(string) string { return "collections.BoolKey" },
	DefaultTestValue:        "true",
	ValueLoop:               "true",
	ValueIndex:              "true",
	ValueInvalidIndex:       "true",
	ProtoType: func(_, name string, index int) string {
		return fmt.Sprintf("bool %s = %d", name, index)
	},
	GenesisArgs: func(name multiformatname.Name, value int) string {
		return fmt.Sprintf("%s: %t,\n", name.UpperCamel, value%2 == 0)
	},
	CLIArgs: func(name multiformatname.Name, _, prefix string, argIndex int) string {
		return fmt.Sprintf(`%s%s, err := cast.ToBoolE(args[%d])
            		if err != nil {
                		return err
            		}`,
			prefix, name.UpperCamel, argIndex)
	},
	ToBytes: func(name string) string {
		return fmt.Sprintf(`%[1]vBytes := []byte{0}
					if %[1]v {
						%[1]vBytes = []byte{1}
					}`, name)
	},
	ToString: func(name string) string {
		return fmt.Sprintf("strconv.FormatBool(%s)", name)
	},
	ToProtoField: func(_, name string, index int) *proto.NormalField {
		return protoutil.NewField(name, "bool", index)
	},
	GoCLIImports: []GoImport{{Name: "github.com/spf13/cast"}},
}
