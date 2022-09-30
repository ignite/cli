package typed

import (
	"errors"
	"fmt"

	"github.com/ignite/cli/ignite/pkg/protoanalysis"
)

const gogoProtoFile = "gogoproto/gogo.proto"

// EnsureGogoProtoImported add the gogo.proto import in the proto file content in case it's not defined
func EnsureGogoProtoImported(protoFile, importPlaceholder string) string {
	err := protoanalysis.IsImported(protoFile, gogoProtoFile)
	if errors.Is(err, protoanalysis.ErrImportNotFound) {
		templateGogoProtoImport := `%[1]v
import "%[2]v";`
		replacementGogoProtoImport := fmt.Sprintf(
			templateGogoProtoImport,
			importPlaceholder,
			gogoProtoFile,
		)
		return replacementGogoProtoImport
	}
	return importPlaceholder
}
