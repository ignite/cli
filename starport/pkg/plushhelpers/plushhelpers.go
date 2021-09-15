package plushhelpers

import (
	"fmt"
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

func mergeGoImports(fields ...field.Fields) string {
	allImports := ""
	exist := make(map[string]struct{})
	for _, fields := range fields {
		for _, goImport := range fields.GoCLIImports() {
			if _, ok := exist[goImport]; ok {
				continue
			}
			exist[goImport] = struct{}{}
			allImports += fmt.Sprintf("\"%s\"\n", goImport)
		}
	}
	return allImports
}

func mergeProtoImports(fields ...field.Fields) string {
	allImports := ""
	exist := make(map[string]struct{})
	for _, fields := range fields {
		for _, goImport := range fields.ProtoImports() {
			if _, ok := exist[goImport]; ok {
				continue
			}
			exist[goImport] = struct{}{}
			allImports += fmt.Sprintf("\"%s\"\n", goImport)
		}
	}
	return allImports
}
