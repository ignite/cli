package typed

import (
	"github.com/tendermint/starport/starport/pkg/protoimport"
)

const gogoProtoFile = "gogoproto/gogo.proto"

// EnsureGogoProtoImported add the gogo.proto import in the proto file content in case it's not defined
func EnsureGogoProtoImported(protoFile, importPlaceholder string) string {
	return protoimport.EnsureProtoImported(gogoProtoFile, protoFile, importPlaceholder)
}
