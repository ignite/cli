// Package field provides methods to parse a field provided in a command with the format name:type
package field

import (
	"fmt"
	"math/rand"

	"github.com/tendermint/starport/starport/pkg/multiformatname"
)

const (
	DataTypeCustom       DataTypeName = "custom"
	DataTypeString       DataTypeName = "string"
	DataTypeStringSlice  DataTypeName = "[]string"
	DataTypeBool         DataTypeName = "bool"
	DataTypeInt          DataTypeName = "int"
	DataTypeIntSlice     DataTypeName = "[]int"
	DataTypeUint         DataTypeName = "uint"
	DataTypeUintSlice    DataTypeName = "[]uint"
	DataTypeSDKCoin      DataTypeName = "sdk.Coin"
	DataTypeSDKCoinSlice DataTypeName = "[]sdk.Coin"

	TypeCustom    = "custom"
	TypeSeparator = ":"
)

type (
	// DataTypeName represents the alias name for the data type
	DataTypeName string
	dataType     struct {
		DataDeclaration   func(datatype string) string
		ProtoDeclaration  func(datatype string) string
		GenesisArgs       func(name multiformatname.Name, value int) string
		ProtoImports      []string
		GoCLIImports      []string
		ValueDefault      string
		ValueLoop         string
		ValueIndex        string
		ValueInvalidIndex string
		ToBytes           func(name string) string
		ToString          func(name string) string
		CLIArgs           func(name multiformatname.Name, datatype, prefix string, argIndex int) string
	}
)

var (
	// SupportedTypes all support data types and definitions
	SupportedTypes = map[DataTypeName]dataType{

		DataTypeString: {
			DataDeclaration:   func(string) string { return "string" },
			ProtoDeclaration:  func(string) string { return "string" },
			ValueDefault:      "xyz",
			ValueLoop:         "strconv.Itoa(i)",
			ValueIndex:        "strconv.Itoa(0)",
			ValueInvalidIndex: "strconv.Itoa(100000)",
			GenesisArgs: func(name multiformatname.Name, value int) string {
				return fmt.Sprintf("%s: \"%s\",\n", name.UpperCamel, name.LowerCamel)
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
		},

		DataTypeStringSlice: {
			DataDeclaration:   func(string) string { return "[]string" },
			ProtoDeclaration:  func(string) string { return "repeated string" },
			ValueDefault:      "abc,xyz",
			ValueLoop:         "[]string{strconv.Itoa(i), strconv.Itoa(i)}",
			ValueIndex:        "[]string{\"0\", \"1\"}",
			ValueInvalidIndex: "[]string{strconv.Itoa(100000), strconv.Itoa(100001)}",
			GenesisArgs: func(name multiformatname.Name, value int) string {
				return fmt.Sprintf("%s: []string{\"%s\"},\n", name.UpperCamel, name.LowerCamel)
			},
			CLIArgs: func(name multiformatname.Name, _, prefix string, argIndex int) string {
				return fmt.Sprintf(`%s%s, err := cast.ToStringSliceE(args[%d])
            		if err != nil {
                		return err
            		}`,
					prefix, name.UpperCamel, argIndex)
			},
			ToBytes: func(name string) string {
				return fmt.Sprintf(`%[1]vBytes := make([]byte, 0)
					for _, v := range %[1]v {
  						%[1]vBytes = append(%[1]vBytes, []byte(v)...)
					}`, name)
			},
			ToString: func(name string) string {
				return fmt.Sprintf("strings.Join(%s, \",\")", name)
			},
			GoCLIImports: []string{"github.com/spf13/cast"},
		},

		DataTypeBool: {
			DataDeclaration:   func(string) string { return "bool" },
			ProtoDeclaration:  func(string) string { return "bool" },
			ValueDefault:      "false",
			ValueLoop:         "false",
			ValueIndex:        "false",
			ValueInvalidIndex: "false",
			GenesisArgs: func(name multiformatname.Name, value int) string {
				return fmt.Sprintf("%s: %t,\n", name.UpperCamel, rand.Intn(value)%2 == 0)
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
			GoCLIImports: []string{"github.com/spf13/cast"},
		},

		DataTypeInt: {
			DataDeclaration:   func(string) string { return "int32" },
			ProtoDeclaration:  func(string) string { return "int32" },
			ValueDefault:      "111",
			ValueLoop:         "int32(i)",
			ValueIndex:        "0",
			ValueInvalidIndex: "100000",
			GenesisArgs: func(name multiformatname.Name, value int) string {
				return fmt.Sprintf("%s: %d,\n", name.UpperCamel, rand.Intn(value))
			},
			CLIArgs: func(name multiformatname.Name, _, prefix string, argIndex int) string {
				return fmt.Sprintf(`%s%s, err := cast.ToIntE(args[%d])
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
			GoCLIImports: []string{"github.com/spf13/cast"},
		},

		DataTypeIntSlice: {
			DataDeclaration:   func(string) string { return "[]int32" },
			ProtoDeclaration:  func(string) string { return "repeated int32" },
			ValueDefault:      "1,2,3,4,5",
			ValueLoop:         "[]int32{int32(i), int32(i)}",
			ValueIndex:        "0,1,2,3",
			ValueInvalidIndex: "100000,100001",
			GenesisArgs: func(name multiformatname.Name, value int) string {
				return fmt.Sprintf("%s: []int32{%d},\n", name.UpperCamel, rand.Intn(value))
			},
			CLIArgs: func(name multiformatname.Name, _, prefix string, argIndex int) string {
				return fmt.Sprintf(`%s%s, err := cast.ToIntSliceE(args[%d])
            		if err != nil {
                		return err
            		}`,
					prefix, name.UpperCamel, argIndex)
			},
			ToBytes: func(name string) string {
				return fmt.Sprintf(`%[1]vBytes := make([]byte, 0)
					for _, v := range %[1]v {
  						binary.BigEndian.PutUint64(%[1]vBytes, uint32(v))
					}`, name)
			},
			ToString: func(name string) string {
				return fmt.Sprintf("strings.Trim(strings.Replace(fmt.Sprint(%s), \" \", \",\", -1), \"[]\")", name)
			},
			GoCLIImports: []string{"github.com/spf13/cast"},
		},

		DataTypeUint: {
			DataDeclaration:   func(string) string { return "uint64" },
			ProtoDeclaration:  func(string) string { return "uint64" },
			ValueDefault:      "111",
			ValueLoop:         "uint64(i)",
			ValueIndex:        "0",
			ValueInvalidIndex: "100000",
			GenesisArgs: func(name multiformatname.Name, value int) string {
				return fmt.Sprintf("%s: %d,\n", name.UpperCamel, rand.Intn(value))
			},
			CLIArgs: func(name multiformatname.Name, _, prefix string, argIndex int) string {
				return fmt.Sprintf(`%s%s, err := cast.ToUintE(args[%d])
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
			GoCLIImports: []string{"github.com/spf13/cast"},
		},

		DataTypeUintSlice: {
			DataDeclaration:   func(string) string { return "[]uint64" },
			ProtoDeclaration:  func(string) string { return "repeated uint64" },
			ValueDefault:      "1,2,3,4,5",
			ValueLoop:         "[]uint64{uint64(i), uint64(i)}",
			ValueIndex:        "0,1,2,3",
			ValueInvalidIndex: "100000,100001",
			GenesisArgs: func(name multiformatname.Name, value int) string {
				return fmt.Sprintf("%s: []uint64{%d},\n", name.UpperCamel, rand.Intn(value))
			},
			CLIArgs: func(name multiformatname.Name, _, prefix string, argIndex int) string {
				return fmt.Sprintf(`%s%s, err := cast.ToUintSliceE(args[%d])
            		if err != nil {
                		return err
            		}`,
					prefix, name.UpperCamel, argIndex)
			},
			ToBytes: func(name string) string {
				return fmt.Sprintf(`%[1]vBytes := make([]byte, 0)
					for _, v := range %[1]v {
  						binary.BigEndian.PutUint64(%[1]vBytes, v)
					}`, name)
			},
			ToString: func(name string) string {
				return fmt.Sprintf("strings.Trim(strings.Replace(fmt.Sprint(%s), \" \", \",\", -1), \"[]\")", name)
			},
			GoCLIImports: []string{"github.com/spf13/cast"},
		},

		DataTypeCustom: {
			DataDeclaration:   func(datatype string) string { return datatype },
			ProtoDeclaration:  func(datatype string) string { return datatype },
			ValueDefault:      "null",
			ValueLoop:         "null",
			ValueIndex:        "null",
			ValueInvalidIndex: "null",
			GenesisArgs: func(name multiformatname.Name, value int) string {
				return fmt.Sprintf("%s: \"%s\",\n", name.UpperCamel, name.LowerCamel)
			},
			CLIArgs: func(name multiformatname.Name, datatype, prefix string, argIndex int) string {
				return fmt.Sprintf(`%[1]v%[2]v := new(types.%[3]v)
					err = json.Unmarshal([]byte(args[%[4]v]), %[1]v%[2]v)
    				if err != nil {
                		return err
            		}`, prefix, name.UpperCamel, datatype, argIndex)
			},
			ToBytes: func(name string) string {
				return fmt.Sprintf(`%[1]vBufferBytes := new(bytes.Buffer)
					json.NewEncoder(%[1]vBytes).Encode(%[1]v)
					%[1]vBytes := reqBodyBytes.Bytes()`, name)
			},
			ToString: func(name string) string {
				return fmt.Sprintf("fmt.Sprintf(\"%s\", %s)", "%+v", name)
			},
			GoCLIImports: []string{"encoding/json"},
		},
		//DataTypeSDKCoin: {
		//	DataDeclaration:  DefaultDataHandler,
		//	ProtoDeclaration: DefaultDataHandler,
		//	GenesisArgs: func(name multiformatname.Name, value int) string {
		//		return fmt.Sprintf("%s: \"%s\",\n", name.UpperCamel, name.LowerCamel)
		//	},
		//	ValueDefault:         "100token",
		//	ValueLoop:   "100token",
		//	ValueIndex:         "100token",
		//	ValueInvalidIndex: "100token",
		//	CLIArgs:       nil,
		//	ToBytes:       nil,
		//	ToString:      nil,
		//	ProtoImports:  nil,
		//	GoCLIImports:     nil,
		//},

		//DataTypeSDKCoinSlice: {
		//	DataDeclaration:  DefaultDataHandler,
		//	ProtoDeclaration: DefaultDataHandler,
		//	GenesisArgs: func(name multiformatname.Name, value int) string {
		//		return fmt.Sprintf("%s: \"%s\",\n", name.UpperCamel, name.LowerCamel)
		//	},
		//	ValueDefault:         "100token,50stake",
		//	ValueLoop:   "100token,50stake",
		//	ValueIndex:         "100token,50stake",
		//	ValueInvalidIndex: "100token,50stake",
		//	CLIArgs:       nil,
		//	ToBytes:       nil,
		//	ToString:      nil,
		//	ProtoImports:  nil,
		//	GoCLIImports:     nil,
		//},
	}
)
