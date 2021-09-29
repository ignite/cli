package datatype

import (
	"github.com/tendermint/starport/starport/pkg/multiformatname"
)

const (
	String      Name = "string"
	StringSlice Name = "array.string"
	Bool        Name = "bool"
	Int         Name = "int"
	IntSlice    Name = "array.int"
	Uint        Name = "uint"
	UintSlice   Name = "array.uint"
	Coin        Name = "coin"
	Coins       Name = "array.coin"
	Custom      Name = Name(TypeCustom)

	StringSliceAlias Name = "strings"
	IntSliceAlias    Name = "ints"
	UintSliceAlias   Name = "uints"
	CoinSliceAlias   Name = "coins"

	TypeCustom    = "customstarporttype"
	TypeSeparator = ":"
)

var (
	// SupportedTypes all support data types and definitions
	SupportedTypes = map[Name]dataType{
		String:           typeString,
		StringSlice:      typeStringSlice,
		StringSliceAlias: typeStringSlice,
		Bool:             typeBool,
		Int:              typeInt,
		IntSlice:         typeIntSlice,
		IntSliceAlias:    typeIntSlice,
		Uint:             typeUint,
		UintSlice:        typeUintSlice,
		UintSliceAlias:   typeUintSlice,
		Coin:             typeCoin,
		Coins:            typeCoinSlice,
		CoinSliceAlias:   typeCoinSlice,
		Custom:           typeCustom,
	}
)

type (
	// Name represents the Alias Name for the data type
	Name string
	// GoImport represents the go import repo name with the alias
	GoImport struct {
		Name  string
		Alias string
	}
	dataType struct {
		DataType          func(datatype string) string
		ProtoType         func(datatype, name string, index int) string
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
