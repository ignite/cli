package datatype

import (
	"fmt"

	"github.com/emicklei/proto"

	"github.com/ignite/cli/ignite/pkg/multiformatname"
	"github.com/ignite/cli/ignite/pkg/protoanalysis/protoutil"
)

var (
	// DataInt is an int data type definition.
	DataInt = DataType{
		DataType:          func(string) string { return "int32" },
		DefaultTestValue:  "111",
		ValueLoop:         "int32(i)",
		ValueIndex:        "0",
		ValueInvalidIndex: "100000",
		ProtoType: func(_, name string, index int) string {
			return fmt.Sprintf("int32 %s = %d", name, index)
		},
		GenesisArgs: func(name multiformatname.Name, value int) string {
			return fmt.Sprintf("%s: %d,\n", name.UpperCamel, value)
		},
		CLIArgs: func(name multiformatname.Name, _, prefix string, argIndex int) string {
			return fmt.Sprintf(`%s%s, err := cast.ToInt32E(args[%d])
            		if err != nil {
                		return err
            		}`,
				prefix, name.UpperCamel, argIndex)
		},
		ToBytes: func(name string) string {
			return fmt.Sprintf(`%[1]vBytes := make([]byte, 4)
  					binary.BigEndian.PutUint32(%[1]vBytes, uint32(%[1]v))`, name)
		},
		ToString: func(name string) string {
			return fmt.Sprintf("strconv.Itoa(int(%s))", name)
		},
		ToProtoField: func(_, name string, index int) *proto.NormalField {
			return protoutil.NewField(name, "int32", index)
		},
		GoCLIImports: []GoImport{{Name: "github.com/spf13/cast"}},
	}

	// DataIntSlice is an int array data type definition.
	DataIntSlice = DataType{
		DataType:         func(string) string { return "[]int32" },
		DefaultTestValue: "1,2,3,4,5",
		ProtoType: func(_, name string, index int) string {
			return fmt.Sprintf("repeated int32 %s = %d", name, index)
		},
		GenesisArgs: func(name multiformatname.Name, value int) string {
			return fmt.Sprintf("%s: []int32{%d},\n", name.UpperCamel, value)
		},
		CLIArgs: func(name multiformatname.Name, _, prefix string, argIndex int) string {
			return fmt.Sprintf(`%[1]vCast%[2]v := strings.Split(args[%[3]v], listSeparator)
					%[1]v%[2]v := make([]int32, len(%[1]vCast%[2]v))
					for i, arg := range %[1]vCast%[2]v {
						value, err := cast.ToInt32E(arg)
						if err != nil {
							return err
						}
						%[1]v%[2]v[i] = value
					}`, prefix, name.UpperCamel, argIndex)
		},
		ToProtoField: func(_, name string, index int) *proto.NormalField {
			return protoutil.NewField(name, "int32", index, protoutil.Repeated())
		},
		GoCLIImports: []GoImport{{Name: "github.com/spf13/cast"}, {Name: "strings"}},
		NonIndex:     true,
	}
)
