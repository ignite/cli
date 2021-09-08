package protoanalysis

import (
	"fmt"

	"github.com/tendermint/starport/starport/pkg/placeholder"
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
func EnsureProtoImported(content, protoImport, protoPath, importPlaceholder string, replacer placeholder.Replacer) string {
	isImported, err := IsImported(protoPath, protoImport)
	if err != nil {
		replacer.AppendMiscError(fmt.Sprintf("failed to check %s dependency %s", protoImport, err.Error()))
		return content
	}
	if !isImported {
		templateGogoProtoImport := `%[1]v
import "%[2]v.proto";`
		replacementGogoProtoImport := fmt.Sprintf(
			templateGogoProtoImport,
			importPlaceholder,
			protoImport,
		)
		content = replacer.Replace(content, importPlaceholder, replacementGogoProtoImport)
	}
	return content
}
