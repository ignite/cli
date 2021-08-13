package typed

import (
	"fmt"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/pkg/protoanalysis"
)

const gogoProtoFile = "gogoproto/gogo.proto"
var gogoProtoImport = fmt.Sprintf(`import "%s";`, gogoProtoFile)

// gogoProtoImported returns true if gogo.proto is imported in the provided proto file
func gogoProtoImported(path string) (bool, error) {
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

// AddGogoProtoImport add the gogo.proto import in the proto file content in case it's not defined
func AddGogoProtoImport(content, protoFile, importPlaceholder string, replacer placeholder.Replacer) string {
	gogoproto, err := gogoProtoImported(protoFile)
	if err != nil {
		replacer.AppendMiscError(fmt.Sprintf("failed to check gogoproto dependency %s", err.Error()))
		return content
	}
	if !gogoproto {
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