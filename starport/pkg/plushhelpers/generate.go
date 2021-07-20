package plushhelpers

import (
	"fmt"
)

// GenerateValidArg will produce a valid value for the specified type
// This function doesn't guarantee to produce unique values
// Note that return value needs to be wrapped into a string
func GenerateValidArg(datatypeName string) string {
	switch datatypeName {
	case datatypeString:
		return "xyz"
	case datatypeUint, datatypeInt:
		return "111"
	case datatypeBool:
		return valueFalse
	default:
		panic(fmt.Sprintf("unknown type %s", datatypeName))
	}
}

// GenerateUniqueArg returns the line of code for the iterated value i for the type datatypeName
// The value is unique depending on i, except for bool which always returns true
// This method must be placed in the template inside a loop with an iterator i
func GenerateUniqueArg(datatypeName string) string {
	switch datatypeName {
	case datatypeString:
		return "strconv.Itoa(i)"
	case datatypeUint:
		return "uint64(i)"
	case datatypeInt:
		return "int32(i)"
	case datatypeBool:
		return valueFalse
	default:
		panic(fmt.Sprintf("unknown type %s", datatypeName))
	}
}

// GenerateValidIndex returns the line of code for a valid index for a map depending on the type
func GenerateValidIndex(datatypeName string) string {
	switch datatypeName {
	case datatypeString:
		return "strconv.Itoa(0)"
	case datatypeUint, datatypeInt:
		return "0"
	case datatypeBool:
		return valueFalse
	default:
		panic(fmt.Sprintf("unknown type %s", datatypeName))
	}
}

// GenerateNotFoundIndex returns the line of code for an index that doesn't exist for a map
// This is used for map tests generation, for test cases where the type is not found for the specified index
// NOTE: This method is not reliable for tests with a map with only booleans as indexes
func GenerateNotFoundIndex(datatypeName string) string {
	switch datatypeName {
	case datatypeString:
		return "strconv.Itoa(100000)"
	case datatypeUint, datatypeInt:
		return "100000"
	case datatypeBool:
		return valueFalse
	default:
		panic(fmt.Sprintf("unknown type %s", datatypeName))
	}
}
