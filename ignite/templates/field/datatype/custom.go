package datatype

import (
	"fmt"

	"github.com/emicklei/proto"

	"github.com/ignite/cli/ignite/pkg/multiformatname"
	"github.com/ignite/cli/ignite/pkg/protoanalysis/protoutil"
)

// DataCustom is a custom data type definition.
var DataCustom = DataType{
	DataType:         func(datatype string) string { return fmt.Sprintf("*%s", datatype) },
	DefaultTestValue: "null",
	ProtoType: func(datatype, name string, index int) string {
		return fmt.Sprintf("%s %s = %d", datatype, name, index)
	},
	GenesisArgs: func(name multiformatname.Name, value int) string {
		return fmt.Sprintf("%s: new(types.%s),\n", name.UpperCamel, name.UpperCamel)
	},
	CLIArgs: func(name multiformatname.Name, datatype, prefix string, argIndex int) string {
		return fmt.Sprintf(`%[1]v%[2]v := new(types.%[3]v)
					err = json.Unmarshal([]byte(args[%[4]v]), %[1]v%[2]v)
    				if err != nil {
                		return err
            		}`, prefix, name.UpperCamel, datatype, argIndex)
	},
	ToProtoField: func(datatype, name string, index int) *proto.NormalField {
		return protoutil.NewField(name, datatype, index)
	},
	GoCLIImports: []GoImport{{Name: "encoding/json"}},
	NonIndex:     true,
}
