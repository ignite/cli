package plushhelpers

import (
	"fmt"
)

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

// GenerateUniqueArg returns the line of code for the iterated value i for the type datatypeName
// The value is unique depending on i, except for bool which always returns true
// This method must be placed in the template inside a loop with an iterator i
func GenerateUniqueArg(datatypeName string) string {
	switch datatypeName {
	case "string":
		return "strconv.Itoa(i)"
	case "uint":
		return "uint64(i)"
	case "int":
		return "int32(i)"
	case "bool":
		return "true"
	default:
		panic(fmt.Sprintf("unknown type %s", datatypeName))
	}
}

// GenerateNotFoundIndex returns the line of code for an index that doesn't exist for a map
// This is used for map tests generation, for test cases where the type is not found for the specified index
// NOTE: This method is not reliable for tests with a map with only booleans as indexes
func GenerateNotFoundIndex(datatypeName string) string {
	switch datatypeName {
	case "string":
		return "not_found"
	case "uint":
		return "^uint64(0)" // max uint64
	case "int":
		return "-1"
	case "bool":
		return "true"
	default:
		panic(fmt.Sprintf("unknown type %s", datatypeName))
	}
}