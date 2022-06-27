package datatype

import (
	"fmt"

	"github.com/ignite/cli/ignite/pkg/multiformatname"
)

var (
	// DataString string data type definition
	DataString = DataType{
		DataType:          func(string) string { return "string" },
		DefaultTestValue:  "xyz",
		ValueLoop:         "strconv.Itoa(i)",
		ValueIndex:        "strconv.Itoa(0)",
		ValueInvalidIndex: "strconv.Itoa(100000)",
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
	}

	// DataStringSlice string array data type definition
	DataStringSlice = DataType{
		DataType:         func(string) string { return "[]string" },
		DefaultTestValue: "abc,xyz",
		ProtoType: func(_, name string, index int) string {
			return fmt.Sprintf("repeated string %s = %d", name, index)
		},
		GenesisArgs: func(name multiformatname.Name, value int) string {
			return fmt.Sprintf("%s: []string{\"%d\"},\n", name.UpperCamel, value)
		},
		CLIArgs: func(name multiformatname.Name, _, prefix string, argIndex int) string {
			return fmt.Sprintf(`%[1]v%[2]v := strings.Split(args[%[3]v], listSeparator)`,
				prefix, name.UpperCamel, argIndex)
		},
		GoCLIImports: []GoImport{{Name: "strings"}},
		NonIndex:     true,
	}
)
