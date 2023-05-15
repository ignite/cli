package cosmosgen

import (
	"go/ast"

	"github.com/ignite/cli/ignite/pkg/goanalysis"
)

// MissingTools find missing tools import indo a *ast.File.
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
