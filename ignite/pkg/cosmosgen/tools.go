package cosmosgen

import (
	"go/ast"

	"github.com/ignite/cli/ignite/pkg/goanalysis"
)

func MissingTools(f *ast.File) (missingTools []string) {
	imports := make(map[string]string)
	for name, imp := range goanalysis.FormatImports(f) {
		imports[imp] = name
	}

	for _, tool := range DepTools() {
		if _, ok := imports[tool]; !ok {
			missingTools = append(missingTools, tool)
		}
	}
	return
}

func UpgradeTools(f *ast.File, addImports, removeImports []string) error {
	return nil
}
