package plushhelpers

import (
	"fmt"
	"strings"
)

// CastArgs returns the line of code to cast a value received from CLI of type string into its datatype
// Don't forget to import github.com/spf13/cast in templates
func CastArgs(datatype string, argIndex int) string {
	return fmt.Sprintf("cast.To%sE(args[%d])", strings.Title(datatype), argIndex)
}

// GenerateValidArg will produce a valid value for the specified type
// This function doesn't guarantee to produce unique values
// Note that return value needs to be wrapped into a string
func GenerateValidArg(datatypeName string) string {
	switch datatypeName {
	case "string":
		return "xyz"
	case "uint":
		return "111"
	case "int":
		return "111"
	case "bool":
		return "true"
	default:
		panic(fmt.Sprintf("unknown type %s", datatypeName))
	}
}
