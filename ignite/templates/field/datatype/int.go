package datatype

import (
	"fmt"

	"github.com/emicklei/proto"

	"github.com/ignite/cli/v29/ignite/pkg/multiformatname"
	"github.com/ignite/cli/v29/ignite/pkg/protoanalysis/protoutil"
)

var (
	// DataInt is an int data type definition.
	DataInt = DataType{
		Name:                    Int,
		DataType:                func(string) string { return "int64" },
		CollectionsKeyValueName: func(string) string { return "collections.Int64Key" },
		DefaultTestValue:        "111",
		ValueLoop:               "int64(i)",
		ValueIndex:              "0",
		ValueInvalidIndex:       "100000",
		ProtoType: func(_, name string, index int) string {
			return fmt.Sprintf("int64 %s = %d", name, index)
		},
		GenesisArgs: func(name multiformatname.Name, value int) string {
			return fmt.Sprintf("%s: %d,\n", name.UpperCamel, value)
		},
		CLIArgs: func(name multiformatname.Name, _, prefix string, argIndex int) string {
			return fmt.Sprintf(`%s%s, err := cast.ToInt64E(args[%d])
            		if err != nil {
                		return err
            		}`,
				prefix, name.UpperCamel, argIndex)
		},
		ToBytes: func(name string) string {
			return fmt.Sprintf(`%[1]vBytes := make([]byte, 4)
  					binary.BigEndian.PutUint64(%[1]vBytes, uint64(%[1]v))`, name)
		},
		ToString: func(name string) string {
			return fmt.Sprintf("strconv.FormatInt(%s, 10)", name)
		},
		ToProtoField: func(_, name string, index int) *proto.NormalField {
			return protoutil.NewField(name, "int64", index)
		},
		GoCLIImports: []GoImport{{Name: "github.com/spf13/cast"}},
	}

	// DataIntSlice is an int array data type definition.
	DataIntSlice = DataType{
		Name:                    IntSlice,
		DataType:                func(string) string { return "[]int64" },
		CollectionsKeyValueName: func(string) string { return collectionValueComment },
		DefaultTestValue:        "5,4,3,2,1",
		ValueLoop:               "[]int64{int64(i+i%1), int64(i+i%2), int64(i+i%3)}",
		ProtoType: func(_, name string, index int) string {
			return fmt.Sprintf("repeated int64 %s = %d", name, index)
		},
		GenesisArgs: func(name multiformatname.Name, value int) string {
			return fmt.Sprintf("%s: []int64{%d},\n", name.UpperCamel, value)
		},
		CLIArgs: func(name multiformatname.Name, _, prefix string, argIndex int) string {
			return fmt.Sprintf(`%[1]vCast%[2]v := strings.Split(args[%[3]v], listSeparator)
					%[1]v%[2]v := make([]int64, len(%[1]vCast%[2]v))
					for i, arg := range %[1]vCast%[2]v {
						value, err := cast.ToInt64E(arg)
						if err != nil {
							return err
						}
						%[1]v%[2]v[i] = value
					}`, prefix, name.UpperCamel, argIndex)
		},
		ToProtoField: func(_, name string, index int) *proto.NormalField {
			return protoutil.NewField(name, "int64", index, protoutil.Repeated())
		},
		GoCLIImports: []GoImport{{Name: "github.com/spf13/cast"}, {Name: "strings"}},
		NonIndex:     true,
	}
)
