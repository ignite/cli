package typed

import (
	"fmt"

	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/pkg/protoanalysis"
)

const gogoProtoFile = "gogoproto/gogo.proto"

var gogoProtoImport = fmt.Sprintf(`import "%s";`, gogoProtoFile)

// isGogoProtoImported returns true if gogo.proto is imported in the provided proto file
func isGogoProtoImported(path string) (bool, error) {
	f, err := protoanalysis.ParseFile(path)
	if err != nil {
		return false, err
	}

	for _, dep := range f.Dependencies {
		if dep == gogoProtoFile {
			return true, nil
		}
	}

	return false, nil
}

// EnsureGogoProtoImported add the gogo.proto import in the proto file content in case it's not defined
func EnsureGogoProtoImported(content, protoFile, importPlaceholder string, replacer placeholder.Replacer) string {
	isImported, err := isGogoProtoImported(protoFile)
	if err != nil {
		replacer.AppendMiscError(fmt.Sprintf("failed to check gogoproto dependency %s", err.Error()))
		return content
	}
	if !isImported {
		templateGogoProtoImport := `%[1]v
%[2]v`
		replacementGogoProtoImport := fmt.Sprintf(
			templateGogoProtoImport,
			importPlaceholder,
			gogoProtoImport,
		)
		content = replacer.Replace(content, importPlaceholder, replacementGogoProtoImport)
	}

	return content
}
