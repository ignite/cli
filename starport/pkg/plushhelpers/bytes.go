package plushhelpers

import (
	"fmt"
)

// CastToBytes returns the lines of code to cast a value into bytes
// the converted variable name is the second returned value
func CastToBytes(varName string, datatypeName string) (string, string) {
	bytesVar := varName + "Bytes"

	switch datatypeName {
	case "string":
		return fmt.Sprintf("%[1]vBytes := []byte(%[1]v)", varName), bytesVar
	case "uint":
		return fmt.Sprintf(`%[1]vBytes := make([]byte, 8)
  		binary.BigEndian.PutUint64(%[1]vBytes, %[1]v)`, varName), bytesVar
	case "int":
		return fmt.Sprintf(`%[1]vBytes := make([]byte, 4)
  		binary.BigEndian.PutUint32(%[1]vBytes, uint32(%[1]v))`, varName), bytesVar
	case "bool":
		return fmt.Sprintf(`%[1]vBytes := []byte("0")
		if %[1]v {
			%[1]vBytes = []byte("1")
		}`, varName), bytesVar
	default:
		panic(fmt.Sprintf("unknown type %s", datatypeName))
	}
}
