package plushhelpers

import (
	"fmt"
	"strings"
)

// CastArgs returns the line of code to cast a value received from CLI of type string into its datatype
// Don't forget to import github.com/spf13/cast in templates
func CastArgs(name, datatypeName, datatype string, argIndex int) string {
	switch datatypeName {
	case datatypeString:
		return fmt.Sprintf("%s := args[%d]", name, argIndex)
	case datatypeUint, datatypeInt, datatypeBool:
		return fmt.Sprintf(`%s, err := cast.To%sE(args[%d])
            if err != nil {
                return err
            }`, name, strings.Title(datatype), argIndex)
	default:
		panic(fmt.Sprintf("unknown type %s", datatype))
	}
}

// CastToBytes returns the lines of code to cast a value into bytes
// the name of the cast type variable is [name]Bytes
func CastToBytes(varName string, datatypeName string) string {
	switch datatypeName {
	case datatypeString:
		return fmt.Sprintf("%[1]vBytes := []byte(%[1]v)", varName)
	case datatypeUint:
		return fmt.Sprintf(`%[1]vBytes := make([]byte, 8)
  		binary.BigEndian.PutUint64(%[1]vBytes, %[1]v)`, varName)
	case datatypeInt:
		return fmt.Sprintf(`%[1]vBytes := make([]byte, 4)
  		binary.BigEndian.PutUint32(%[1]vBytes, uint32(%[1]v))`, varName)
	case datatypeBool:
		return fmt.Sprintf(`%[1]vBytes := []byte{0}
		if %[1]v {
			%[1]vBytes = []byte{1}
		}`, varName)
	default:
		panic(fmt.Sprintf("unknown type %s", datatypeName))
	}
}

// CastToString returns the lines of code to cast a value into bytes
func CastToString(varName string, datatypeName string) string {
	switch datatypeName {
	case datatypeString:
		return varName
	case datatypeUint, datatypeInt:
		return fmt.Sprintf("strconv.Itoa(int(%s))", varName)
	case datatypeBool:
		return fmt.Sprintf("strconv.FormatBool(%s)", varName)
	default:
		panic(fmt.Sprintf("unknown type %s", datatypeName))
	}
}
