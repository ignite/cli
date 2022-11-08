package plushhelpers

import (
	"strings"

	"github.com/gobuffalo/plush/v4"

	"github.com/ignite/cli/ignite/pkg/xstrings"
	"github.com/ignite/cli/ignite/templates/field"
	"github.com/ignite/cli/ignite/templates/field/datatype"
)

// ExtendPlushContext sets available field helpers on the provided context.
func ExtendPlushContext(ctx *plush.Context) {
	ctx.Set("mergeGoImports", mergeGoImports)
	ctx.Set("mergeProtoImports", mergeProtoImports)
	ctx.Set("mergeCustomImports", mergeCustomImports)
	ctx.Set("title", xstrings.Title)
	ctx.Set("toLower", strings.ToLower)
}

func mergeCustomImports(fields ...field.Fields) []string {
	allImports := make([]string, 0)
	exist := make(map[string]struct{})
	for _, fields := range fields {
		for _, customImport := range fields.Custom() {
			if _, ok := exist[customImport]; ok {
				continue
			}
			exist[customImport] = struct{}{}
			allImports = append(allImports, customImport)
		}
	}
	return allImports
}

func mergeGoImports(fields ...field.Fields) []datatype.GoImport {
	allImports := make([]datatype.GoImport, 0)
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
