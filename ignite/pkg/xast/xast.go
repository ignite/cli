package xast

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/pkg/errors"
)

var ErrStop = errors.New("ast stop")

// Inspect is like ast.Inspect but with error handling.
// Unlike ast.Inspect the function parameter f returns an error and not a bool.
// The returned error is propagated to the caller, unless it is equal to
// ErrStop, which in that case indicates the child nodes shouldn't not be
// inspected (like returning false in the function of ast.Inspect).
func Inspect(n ast.Node, f func(n ast.Node) error) (err error) {
	ast.Inspect(n, func(n ast.Node) bool {
		err = f(n)
		if err == nil {
			return true
		}
		if errors.Is(err, ErrStop) {
			err = nil
		}
		return false
	})
	return
}

// ParseDir invokes ast.ParseDir and returns the first package found that is
// doesn't has the "_test" suffix.
func ParseDir(dir string) (*ast.Package, *token.FileSet, error) {
	fileSet := token.NewFileSet()
	pkgs, err := parser.ParseDir(fileSet, dir, nil, 0)
	if err != nil {
		return nil, nil, err
	}
	for name, pkg := range pkgs {
		if strings.HasSuffix(name, "_test") {
			continue
		}
		return pkg, fileSet, nil
	}
	return nil, nil, errors.Errorf("no valid package found in %s", dir)
}

// ParseFile invokes ast.ParseFile and returns the *ast.File.
func ParseFile(filepath string) (*ast.File, *token.FileSet, error) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, filepath, nil, 0)
	return file, fileSet, err
}
