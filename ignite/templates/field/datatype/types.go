package datatype

import (
	"fmt"
	"io"
	"sort"

	"github.com/emicklei/proto"

	"github.com/ignite/cli/v29/ignite/pkg/cliui/entrywriter"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/multiformatname"
)

const (
	// Separator represents the type separator.
	Separator = ":"

	// String represents the string type name.
	String Name = "string"
	// StringSlice represents the string array type name.
	StringSlice Name = "array.string"
	// Bool represents the bool type name.
	Bool Name = "bool"
	// Int represents the int type name.
	Int Name = "int"
	// Int64 represents the int64 type name.
	Int64 Name = "int64"
	// IntSlice represents the int array type name.
	IntSlice Name = "array.int"
	// Uint represents the uint type name.
	Uint Name = "uint"
	// Uint64 represents the uint64 type name.
	Uint64 Name = "uint64"
	// UintSlice represents the uint array type name.
	UintSlice Name = "array.uint"
	// Coin represents the coin type name.
	Coin Name = "coin"
	// Coins represents the coin array type name.
	Coins Name = "array.coin"
	// DecCoin represents the coin type name.
	DecCoin Name = "dec.coin"
	// DecCoins represents the decimal coin array type name.
	DecCoins Name = "array.dec.coin"
	// Bytes represents the bytes type name.
	Bytes Name = "bytes"
	// Address represents the address type name.
	Address Name = "address"
	// Custom represents the custom type name.
	Custom Name = Name(TypeCustom)

	// StringSliceAlias represents the string array type name alias.
	StringSliceAlias Name = "strings"
	// IntSliceAlias represents the int array type name alias.
	IntSliceAlias Name = "ints"
	// UintSliceAlias represents the uint array type name alias.
	UintSliceAlias Name = "uints"
	// CoinSliceAlias represents the coin array type name alias.
	CoinSliceAlias Name = "coins"
	// DecCoinSliceAlias represents the coin array type name alias.
	DecCoinSliceAlias Name = "dec.coins"

	// TypeCustom represents the string type name id.
	TypeCustom = "customignitetype"

	collectionValueComment = "/* Add collection key value */"
)

// supportedTypes all support data types and definitions.
var supportedTypes = map[Name]DataType{
	Bytes:             DataBytes,
	String:            DataString,
	StringSlice:       DataStringSlice,
	StringSliceAlias:  DataStringSlice,
	Bool:              DataBool,
	Int:               DataInt,
	Int64:             DataInt,
	IntSlice:          DataIntSlice,
	IntSliceAlias:     DataIntSlice,
	Uint:              DataUint,
	Uint64:            DataUint,
	UintSlice:         DataUintSlice,
	UintSliceAlias:    DataUintSlice,
	Coin:              DataCoin,
	Coins:             DataCoinSlice,
	CoinSliceAlias:    DataCoinSlice,
	DecCoin:           DataDecCoin,
	DecCoins:          DataDecCoinSlice,
	DecCoinSliceAlias: DataDecCoinSlice,
	Address:           DataAddress,
	Custom:            DataCustom,
}

// Name represents the Alias Name for the data type.
type Name string

// DataType represents the data types for code replacement.
type DataType struct {
	Name                    Name
	DataType                func(datatype string) string
	ProtoType               func(datatype, name string, index int) string
	CollectionsKeyValueName func(datatype string) string
	GenesisArgs             func(name multiformatname.Name, value int) string
	ProtoImports            []string
	GoCLIImports            GoImports
	DefaultTestValue        string
	ValueLoop               string
	ValueIndex              string
	ValueInvalidIndex       string
	ToBytes                 func(name string) string
	ToString                func(name string) string
	ToProtoField            func(datatype, name string, index int) *proto.NormalField
	CLIArgs                 func(name multiformatname.Name, datatype, prefix string, argIndex int) string
	NonIndex                bool
}

// Usage returns the usage of the data type.
// It provides a description of how to use the data type in scaffolding.
func (t DataType) Usage() string {
	if t.Name == Custom {
		return "use the custom type to scaffold already created chain types."
	}
	usage := fmt.Sprintf("use '<FIELD_NAME>:%s' to scaffold %s types (eg: %s).", t.Name, t.DataType(""), t.DefaultTestValue)
	if t.Name == Coins || t.Name == DecCoins ||
		t.Name == CoinSliceAlias || t.Name == DecCoinSliceAlias {
		return usage + " Disclaimer: Only one `coins` or `dec.coins` field can accept multiple CLI values per command due to AutoCLI limitations."
	}
	return usage
}

// GoImports represents a list of go import.
type GoImports []GoImport

// GoImport represents the go import repo name with the alias.
type GoImport struct {
	Name  string
	Alias string
}

// IsSupportedType type checks if the given typename is supported by ignite scaffolding.
// Returns corresponding Datatype if supported.
func IsSupportedType(typename Name) (dt DataType, ok bool) {
	dt, ok = supportedTypes[typename]
	return
}

// SupportedTypes return a list of supported types.
func SupportedTypes() map[string]string {
	supported := make(map[string]string)
	for name, dataType := range supportedTypes {
		if dataType.Name == Custom {
			name = "custom"
		}
		supported[string(name)] = dataType.Usage()
	}
	return supported
}

// PrintScaffoldTypeList prints the list of supported scaffold types to the given writer.
func PrintScaffoldTypeList(writer io.Writer) error {
	supported := SupportedTypes()
	entries := make([][]string, 0, len(supported))
	for name, usage := range supported {
		entries = append(entries, []string{name, usage})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i][0] < entries[j][0]
	})

	if err := entrywriter.MustWrite(writer, []string{"types", "usage"}, entries...); err != nil {
		return errors.Errorf("failed to write scaffold types: %w", err)
	}

	const footer = `Field Usage:
    - fieldName
    - fieldName:fieldType

If no :fieldType, default (string) is used
`

	_, err := fmt.Fprint(writer, footer)
	return err
}
