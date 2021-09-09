package protoanalysis

import (
	"fmt"
)

// IsImported returns true if the proto file is imported in the provided proto file
func IsImported(path, protoImport string) (bool, error) {
	f, err := ParseFile(path)
	if err != nil {
		return false, err
	}

	protoImport = fmt.Sprintf("%s.proto", protoImport)
	for _, dep := range f.Dependencies {
		if dep == protoImport {
			return true, nil
		}
	}

	return false, nil
}

// EnsureProtoImported add the proto file import in the proto file content in case it's not defined
func EnsureProtoImported(protoImport, protoPath, importPlaceholder string) string {
	isImported, err := IsImported(protoPath, protoImport)
	if err != nil {
		return importPlaceholder
	}
	if !isImported {
		templateGogoProtoImport := `%[1]v
import "%[2]v.proto";`
		replacementGogoProtoImport := fmt.Sprintf(
			templateGogoProtoImport,
			importPlaceholder,
			protoImport,
		)
		return replacementGogoProtoImport
	}
	return importPlaceholder
}
