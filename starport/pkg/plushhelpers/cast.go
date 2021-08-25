package plushhelpers

import (
	"fmt"
	"strings"

	"github.com/tendermint/starport/starport/pkg/field"
)

// castArg returns the line of code to cast a value received from CLI of type string into its datatype
// Don't forget to import github.com/spf13/cast in templates
func castArg(prefix string, field field.Field, argIndex int) string {
	switch field.DatatypeName {
	case datatypeString:
		return fmt.Sprintf("%s%s := args[%d]", prefix, field.Name.UpperCamel, argIndex)
	case datatypeUint, datatypeInt, datatypeBool:
		return fmt.Sprintf(`%s%s, err := cast.To%sE(args[%d])
            if err != nil {
                return err
            }`,
			prefix, field.Name.UpperCamel, strings.Title(field.Datatype), argIndex)
	case datatypeCustom:
		return fmt.Sprintf(`%[1]v%[2]v := new(types.%[3]v)
			err = json.Unmarshal([]byte(args[%[4]v]), %[1]v%[2]v)
    		if err != nil {
                return err
            }`, prefix, field.Name.UpperCamel, field.Datatype, argIndex)
	default:
		panic(fmt.Sprintf("unknown type %s", field.DatatypeName))
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
	case datatypeCustom:
		return fmt.Sprintf(`%[1]vBufferBytes := new(bytes.Buffer)
		json.NewEncoder(%[1]vBytes).Encode(%[1]v)
		%[1]vBytes := reqBodyBytes.Bytes()`, varName)
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
	case datatypeCustom:
		return fmt.Sprintf("fmt.Sprintf(\"%s\", %s)", "%+v", varName)
	default:
		panic(fmt.Sprintf("unknown type %s", datatypeName))
	}
}
