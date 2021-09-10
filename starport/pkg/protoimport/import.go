package protoimport

import (
	"fmt"
	"strings"

	"github.com/tendermint/starport/starport/pkg/protoanalysis"
)

const protoFileSuffix = ".proto"

// IsImported returns true if the proto file is imported in the provided proto file
func IsImported(protoImport, protoPath string) (bool, error) {
	f, err := protoanalysis.ParseFile(protoPath)
	if err != nil {
		return false, err
	}

	for _, dep := range f.Dependencies {
		if dep == protoImport {
			return true, nil
		}
	}

	return false, nil
}

// EnsureProtoImported checks if the import already exist and return the new import
func EnsureProtoImported(protoImport, protoPath, importPlaceholder string) string {
	if !strings.HasSuffix(protoImport, protoFileSuffix) {
		protoImport += protoFileSuffix
	}

	isImported, err := IsImported(protoImport, protoPath)
	if err != nil {
		return importPlaceholder
	}
	if !isImported {
		templateGogoProtoImport := `%[1]v
import "%[2]v";`
		replacementGogoProtoImport := fmt.Sprintf(
			templateGogoProtoImport,
			importPlaceholder,
			protoImport,
		)
		return replacementGogoProtoImport
	}
	return importPlaceholder
}
