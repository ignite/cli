package diff

import (
	"path/filepath"

	"github.com/gobwas/glob"
)

// List of files that should be ignored when calculating diff of two directories
var exceptionFiles = []glob.Glob{
	mustCompilePathGlob("**/.git/**"),
	mustCompilePathGlob("**.md"),
	mustCompilePathGlob("**/go.sum"),
	mustCompilePathGlob("**_test.go"),
	mustCompilePathGlob("**.pb.go"),
	mustCompilePathGlob("**.pb.gw.go"),
	mustCompilePathGlob("**.pulsar.go"),
	mustCompilePathGlob("**/node_modules/**"),
	mustCompilePathGlob("**/openapi.yml"),
	mustCompilePathGlob("**/.gitignore"),
	mustCompilePathGlob("**.html"),
	mustCompilePathGlob("**.css"),
	mustCompilePathGlob("**.js"),
	mustCompilePathGlob("**.ts"),
}

func mustCompilePathGlob(pattern string) glob.Glob {
	return glob.MustCompile(pattern, filepath.Separator)
}

// Checks if the given path matches any of the exception file patterns
func isException(path string) bool {
	for _, glob := range exceptionFiles {
		if glob.Match(path) {
			return true
		}
	}
	return false
}
