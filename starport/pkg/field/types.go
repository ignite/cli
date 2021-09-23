package field

import (
	"fmt"

	"github.com/tendermint/starport/starport/pkg/multiformatname"
)

const (
	DataTypeString      DataTypeName = "string"
	DataTypeStringSlice DataTypeName = "array.string"
	DataTypeBool        DataTypeName = "bool"
	DataTypeInt         DataTypeName = "int"
	DataTypeIntSlice    DataTypeName = "array.int"
	DataTypeUint        DataTypeName = "uint"
	DataTypeUintSlice   DataTypeName = "array.uint"
	DataTypeCoin        DataTypeName = "coin"
	DataTypeCoinSlice   DataTypeName = "array.coin"
	DataTypeCustom      DataTypeName = DataTypeName(TypeCustom)

	DataTypeStringSliceAlias DataTypeName = "strings"
	DataTypeIntSliceAlias    DataTypeName = "ints"
	DataTypeUintSliceAlias   DataTypeName = "uints"
	DataTypeCoinSliceAlias   DataTypeName = "coins"

	TypeCustom    = "customstarporttype"
	TypeSeparator = ":"
)

type (
	// DataTypeName represents the Alias Name for the data type
	DataTypeName string
	// GoImport represents the go import repo name with the alias
	GoImport struct {
		Name  string
		Alias string
	}
	dataType struct {
		DataDeclaration   func(datatype string) string
		ProtoDeclaration  func(datatype, name string, index int) string
		GenesisArgs       func(name multiformatname.Name, value int) string
		ProtoImports      []string
		GoCLIImports      []GoImport
		ValueDefault      string
		ValueLoop         string
		ValueIndex        string
		ValueInvalidIndex string
		ToBytes           func(name string) string
		ToString          func(name string) string
		CLIArgs           func(name multiformatname.Name, datatype, prefix string, argIndex int) string
		NonIndex          bool
	}
)

var (

	// SupportedTypes all support data types and definitions
	SupportedTypes = map[DataTypeName]dataType{
		DataTypeString:           typeString,
		DataTypeStringSlice:      typeStringSlice,
		DataTypeStringSliceAlias: typeStringSlice,
		DataTypeBool:             typeBool,
		DataTypeInt:              typeInt,
		DataTypeIntSlice:         typeIntSlice,
		DataTypeIntSliceAlias:    typeIntSlice,
		DataTypeUint:             typeUint,
		DataTypeUintSlice:        typeUintSlice,
		DataTypeUintSliceAlias:   typeUintSlice,
		DataTypeCoin:             typeCoin,
		DataTypeCoinSlice:        typeCoinSlice,
		DataTypeCoinSliceAlias:   typeCoinSlice,
		DataTypeCustom:           typeCustom,
	}

	typeString = dataType{
		DataDeclaration:   func(string) string { return "string" },
		ValueDefault:      "xyz",
		ValueLoop:         "strconv.Itoa(i)",
		ValueIndex:        "strconv.Itoa(0)",
		ValueInvalidIndex: "strconv.Itoa(100000)",
		ProtoDeclaration: func(_, name string, index int) string {
			return fmt.Sprintf("string %s = %d;", name, index)
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

	typeBool = dataType{
		DataDeclaration:   func(string) string { return "bool" },
		ValueDefault:      "false",
		ValueLoop:         "false",
		ValueIndex:        "false",
		ValueInvalidIndex: "false",
		ProtoDeclaration: func(_, name string, index int) string {
			return fmt.Sprintf("bool %s = %d;", name, index)
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
		GoCLIImports: []GoImport{{Name: "github.com/spf13/cast"}},
	}

	typeInt = dataType{
		DataDeclaration:   func(string) string { return "int32" },
		ValueDefault:      "111",
		ValueLoop:         "int32(i)",
		ValueIndex:        "0",
		ValueInvalidIndex: "100000",
		ProtoDeclaration: func(_, name string, index int) string {
			return fmt.Sprintf("int32 %s = %d;", name, index)
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
		GoCLIImports: []GoImport{{Name: "github.com/spf13/cast"}},
	}

	typeUint = dataType{
		DataDeclaration:   func(string) string { return "uint64" },
		ValueDefault:      "111",
		ValueLoop:         "uint64(i)",
		ValueIndex:        "0",
		ValueInvalidIndex: "100000",
		ProtoDeclaration: func(_, name string, index int) string {
			return fmt.Sprintf("uint64 %s = %d;", name, index)
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
		GoCLIImports: []GoImport{{Name: "github.com/spf13/cast"}},
	}

	typeStringSlice = dataType{
		DataDeclaration: func(string) string { return "[]string" },
		ValueDefault:    "abc,xyz",
		ProtoDeclaration: func(_, name string, index int) string {
			return fmt.Sprintf("repeated string %s = %d;", name, index)
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

	typeIntSlice = dataType{
		DataDeclaration: func(string) string { return "[]int32" },
		ValueDefault:    "1,2,3,4,5",
		ProtoDeclaration: func(_, name string, index int) string {
			return fmt.Sprintf("repeated int32 %s = %d;", name, index)
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
		GoCLIImports: []GoImport{{Name: "github.com/spf13/cast"}, {Name: "strings"}},
		NonIndex:     true,
	}

	typeUintSlice = dataType{
		DataDeclaration: func(string) string { return "[]uint64" },
		ValueDefault:    "1,2,3,4,5",
		ProtoDeclaration: func(_, name string, index int) string {
			return fmt.Sprintf("repeated uint64 %s = %d;", name, index)
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
		GoCLIImports: []GoImport{{Name: "github.com/spf13/cast"}, {Name: "strings"}},
		NonIndex:     true,
	}

	typeCustom = dataType{
		DataDeclaration: func(datatype string) string { return fmt.Sprintf("*%s", datatype) },
		ValueDefault:    "null",
		ProtoDeclaration: func(datatype, name string, index int) string {
			return fmt.Sprintf("%s %s = %d;", datatype, name, index)
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
		GoCLIImports: []GoImport{{Name: "encoding/json"}},
		NonIndex:     true,
	}

	typeCoin = dataType{
		DataDeclaration: func(string) string { return "sdk.Coin" },
		ValueDefault:    "10token",
		ProtoDeclaration: func(_, name string, index int) string {
			return fmt.Sprintf("cosmos.base.v1beta1.Coin %s = %d [(gogoproto.nullable) = false];",
				name, index)
		},
		GenesisArgs: func(multiformatname.Name, int) string { return "" },
		CLIArgs: func(name multiformatname.Name, _, prefix string, argIndex int) string {
			return fmt.Sprintf(`%s%s, err := sdk.ParseCoinNormalized(args[%d])
					if err != nil {
						return err
					}`, prefix, name.UpperCamel, argIndex)
		},
		GoCLIImports: []GoImport{{Name: "github.com/cosmos/cosmos-sdk/types", Alias: "sdk"}},
		ProtoImports: []string{"gogoproto/gogo.proto", "cosmos/base/v1beta1/coin.proto"},
		NonIndex:     true,
	}

	typeCoinSlice = dataType{
		DataDeclaration: func(string) string { return "sdk.Coins" },
		ValueDefault:    "10token,20stake",
		ProtoDeclaration: func(_, name string, index int) string {
			return fmt.Sprintf("repeated cosmos.base.v1beta1.Coin %s = %d [(gogoproto.nullable) = false];",
				name, index)
		},
		GenesisArgs: func(multiformatname.Name, int) string { return "" },
		CLIArgs: func(name multiformatname.Name, _, prefix string, argIndex int) string {
			return fmt.Sprintf(`%s%s, err := sdk.ParseCoinsNormalized(args[%d])
					if err != nil {
						return err
					}`, prefix, name.UpperCamel, argIndex)
		},
		GoCLIImports: []GoImport{{Name: "github.com/cosmos/cosmos-sdk/types", Alias: "sdk"}},
		ProtoImports: []string{"gogoproto/gogo.proto", "cosmos/base/v1beta1/coin.proto"},
		NonIndex:     true,
	}
)
