package datatype

import (
	"fmt"

	"github.com/emicklei/proto"

	"github.com/ignite/cli/v29/ignite/pkg/multiformatname"
	"github.com/ignite/cli/v29/ignite/pkg/protoanalysis/protoutil"
)

var (
	// DataUint uint data type definition.
	DataUint = DataType{
		Name:                    Uint,
		DataType:                func(string) string { return "uint64" },
		CollectionsKeyValueName: func(string) string { return "collections.Uint64Key" },
		DefaultTestValue:        "111",
		ValueLoop:               "uint64(i)",
		ValueIndex:              "0",
		ValueInvalidIndex:       "100000",
		ProtoType: func(_, name string, index int) string {
			return fmt.Sprintf("uint64 %s = %d", name, index)
		},
		GenesisArgs: func(name multiformatname.Name, value int) string {
			return fmt.Sprintf("%s: %d,\n", name.UpperCamel, value)
		},
		CLIArgs: func(name multiformatname.Name, _, prefix string, argIndex int) string {
			return fmt.Sprintf(`%s%s, err := cast.ToUint64E(args[%d])
            		if err != nil {
                		return err
            		}`,
				prefix, name.UpperCamel, argIndex)
		},
		ToBytes: func(name string) string {
			return fmt.Sprintf(`%[1]vBytes := make([]byte, 8)
  					binary.BigEndian.PutUint64(%[1]vBytes, %[1]v)`, name)
		},
		ToString: func(name string) string {
			return fmt.Sprintf("strconv.Itoa(int(%s))", name)
		},
		ToProtoField: func(_, name string, index int) *proto.NormalField {
			return protoutil.NewField(name, "uint64", index)
		},
		GoCLIImports: []GoImport{{Name: "github.com/spf13/cast"}},
	}

	// DataUintSlice uint array data type definition.
	DataUintSlice = DataType{
		Name:                    UintSlice,
		DataType:                func(string) string { return "[]uint64" },
		CollectionsKeyValueName: func(string) string { return collectionValueComment },
		DefaultTestValue:        "13,26,31,40",
		ValueLoop:               "[]uint64{uint64(i+i%1), uint64(i+i%2), uint64(i+i%3)}",
		ProtoType: func(_, name string, index int) string {
			return fmt.Sprintf("repeated uint64 %s = %d", name, index)
		},
		GenesisArgs: func(name multiformatname.Name, value int) string {
			return fmt.Sprintf("%s: []uint64{%d},\n", name.UpperCamel, value)
		},
		CLIArgs: func(name multiformatname.Name, _, prefix string, argIndex int) string {
			return fmt.Sprintf(`%[1]vCast%[2]v := strings.Split(args[%[3]v], listSeparator)
					%[1]v%[2]v := make([]uint64, len(%[1]vCast%[2]v))
					for i, arg := range %[1]vCast%[2]v {
						value, err := cast.ToUint64E(arg)
						if err != nil {
							return err
						}
						%[1]v%[2]v[i] = value
					}`,
				prefix, name.UpperCamel, argIndex)
		},
		ToProtoField: func(_, name string, index int) *proto.NormalField {
			return protoutil.NewField(name, "uint64", index, protoutil.Repeated())
		},
		GoCLIImports: []GoImport{{Name: "github.com/spf13/cast"}, {Name: "strings"}},
		NonIndex:     true,
	}
)
