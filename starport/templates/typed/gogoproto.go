package typed

import (
	"fmt"
	"github.com/tendermint/starport/starport/pkg/protoanalysis"
)

const gogoproto = "gogoproto/gogo.proto"

// GogoprotoImport is the import statement to import gogo.proto
var GogoprotoImport = fmt.Sprintf(`import "%s";`, gogoproto)

// GogoprotoImported returns true if gogo.proto is imported in the provided proto file
func GogoprotoImported(path string) (bool, error) {
	f, err := protoanalysis.ParseFile(path)
	if err != nil {
		return false, err
	}

	for _, dep := range f.Dependencies {
		if dep == gogoproto {
			return true, nil
		}
	}

	return false, nil
}