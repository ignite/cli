package typed

import (
	"github.com/tendermint/starport/starport/pkg/protoanalysis"
)

const gogoProtoFile = "gogoproto/gogo"

// EnsureGogoProtoImported add the gogo.proto import in the proto file content in case it's not defined
func EnsureGogoProtoImported(protoFile, importPlaceholder string) string {
	return protoanalysis.EnsureProtoImported(gogoProtoFile, protoFile, importPlaceholder)
}
