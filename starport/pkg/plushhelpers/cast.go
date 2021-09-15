package plushhelpers

import (
	"fmt"
	"strings"

	"github.com/tendermint/starport/starport/pkg/field"
)

// castArg returns the line of code to cast a value received from CLI of type string into its datatype
// Don't forget to import github.com/spf13/cast in templates
func castArg(prefix string, f field.Field, argIndex int) string {
	switch f.DatatypeName {
	case field.TypeString:
		return fmt.Sprintf("%s%s := args[%d]", prefix, f.Name.UpperCamel, argIndex)
	case field.TypeUint, field.TypeInt, field.TypeBool:
		return fmt.Sprintf(`%s%s, err := cast.To%sE(args[%d])
            if err != nil {
                return err
            }`,
			prefix, f.Name.UpperCamel, strings.Title(f.Datatype), argIndex)
	case field.TypeUintSlice, field.TypeIntSlice, field.TypeStringSlice:
		return fmt.Sprintf(`%s%s, err := cast.To%sSliceE(args[%d])
            if err != nil {
                return err
            }`,
			prefix, f.Name.UpperCamel, strings.Title(f.Datatype), argIndex)
	case field.TypeCustom:
		return fmt.Sprintf(`%[1]v%[2]v := new(types.%[3]v)
			err = json.Unmarshal([]byte(args[%[4]v]), %[1]v%[2]v)
    		if err != nil {
                return err
            }`, prefix, f.Name.UpperCamel, f.Datatype, argIndex)
	default:
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
}

// CastToBytes returns the lines of code to cast a value into bytes
// the name of the cast type variable is [name]Bytes
func CastToBytes(varName string, datatypeName string) string {
	switch datatypeName {
	case field.TypeString:
		return fmt.Sprintf("%[1]vBytes := []byte(%[1]v)", varName)
	case field.TypeStringSlice:
		return fmt.Sprintf(`%[1]vBytes := make([]byte, 0)
		for _, v := range %[1]v {
  			%[1]vBytes = append(%[1]vBytes, []byte(v)...)
		}`, varName)
	case field.TypeUint:
		return fmt.Sprintf(`%[1]vBytes := make([]byte, 8)
  		binary.BigEndian.PutUint64(%[1]vBytes, %[1]v)`, varName)
	case field.TypeUintSlice:
		return fmt.Sprintf(`%[1]vBytes := make([]byte, 0)
		for _, v := range %[1]v {
  			binary.BigEndian.PutUint64(%[1]vBytes, v)
		}`, varName)
	case field.TypeInt:
		return fmt.Sprintf(`%[1]vBytes := make([]byte, 4)
  		binary.BigEndian.PutUint32(%[1]vBytes, uint32(%[1]v))`, varName)
	case field.TypeIntSlice:
		return fmt.Sprintf(`%[1]vBytes := make([]byte, 0)
		for _, v := range %[1]v {
  			binary.BigEndian.PutUint64(%[1]vBytes, uint32(v))
		}`, varName)
	case field.TypeBool:
		return fmt.Sprintf(`%[1]vBytes := []byte{0}
		if %[1]v {
			%[1]vBytes = []byte{1}
		}`, varName)
	case field.TypeCustom:
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
	case field.TypeString:
		return varName
	case field.TypeStringSlice:
		return fmt.Sprintf("strings.Join(%s, \",\")", varName)
	case field.TypeUint, field.TypeInt:
		return fmt.Sprintf("strconv.Itoa(int(%s))", varName)
	case field.TypeUintSlice, field.TypeIntSlice:
		return fmt.Sprintf("strings.Trim(strings.Replace(fmt.Sprint(%s), \" \", \",\", -1), \"[]\")", varName)
	case field.TypeBool:
		return fmt.Sprintf("strconv.FormatBool(%s)", varName)
	case field.TypeCustom:
		return fmt.Sprintf("fmt.Sprintf(\"%s\", %s)", "%+v", varName)
	default:
		panic(fmt.Sprintf("unknown type %s", datatypeName))
	}
}
