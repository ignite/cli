package typed

import (
	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/pkg/protoanalysis"
)

const gogoProtoFile = "gogoproto/gogo"

// EnsureGogoProtoImported add the gogo.proto import in the proto file content in case it's not defined
func EnsureGogoProtoImported(content, protoFile, importPlaceholder string, replacer placeholder.Replacer) string {
	return protoanalysis.EnsureProtoImported(content, gogoProtoFile, protoFile, importPlaceholder, replacer)
}
