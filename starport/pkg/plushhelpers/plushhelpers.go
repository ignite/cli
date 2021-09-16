package plushhelpers

import (
	"strings"

	"github.com/gobuffalo/plush"
	"github.com/tendermint/starport/starport/pkg/field"
)

// ExtendPlushContext sets available helpers on the provided context.
func ExtendPlushContext(ctx *plush.Context) {
	ctx.Set("mergeGoImports", mergeGoImports)
	ctx.Set("mergeProtoImports", mergeProtoImports)
	ctx.Set("title", strings.Title)
}

func mergeGoImports(fields ...field.Fields) []field.GoImport {
	allImports := make([]field.GoImport, 0)
	exist := make(map[string]struct{})
	for _, fields := range fields {
		for _, goImport := range fields.GoCLIImports() {
			if _, ok := exist[goImport.Name]; ok {
				continue
			}
			exist[goImport.Name] = struct{}{}
			allImports = append(allImports, goImport)
		}
	}
	return allImports
}

func mergeProtoImports(fields ...field.Fields) []string {
	allImports := make([]string, 0)
	exist := make(map[string]struct{})
	for _, fields := range fields {
		for _, protoImport := range fields.ProtoImports() {
			if _, ok := exist[protoImport]; ok {
				continue
			}
			exist[protoImport] = struct{}{}
			allImports = append(allImports, protoImport)
		}
	}
	return allImports
}
