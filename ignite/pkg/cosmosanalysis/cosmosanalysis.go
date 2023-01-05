// Package cosmosanalysis provides a toolset for statically analysing Cosmos SDK's
// source code and blockchain source codes based on the Cosmos SDK
package cosmosanalysis

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"golang.org/x/mod/modfile"

	"github.com/ignite/cli/ignite/pkg/gomodule"
)

const (
	cosmosModulePath     = "github.com/cosmos/cosmos-sdk"
	tendermintModulePath = "github.com/tendermint/tendermint"
	appFileName          = "app.go"
	defaultAppFilePath   = "app/" + appFileName
)

var appImplementation = []string{
	"Name",
	"BeginBlocker",
	"EndBlocker",
}

// implementation tracks the implementation of an interface for a given struct.
type implementation map[string]bool

// DeepFindImplementation functions the same as FindImplementation, but walks recursively through the folder structure
// Useful if implementations might be in sub folders.
func DeepFindImplementation(modulePath string, interfaceList []string) (found []string, err error) {
	err = filepath.Walk(modulePath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				return nil
			}

			currFound, err := FindImplementation(path, interfaceList)
			if err != nil {
				return err
			}

			found = append(found, currFound...)
			return nil
		})

	if err != nil {
		return nil, err
	}

	return found, nil
}

// FindImplementation finds the name of all types that implement the provided interface.
func FindImplementation(modulePath string, interfaceList []string) (found []string, err error) {
	// parse go packages/files under path
	fset := token.NewFileSet()

	pkgs, err := parser.ParseDir(fset, modulePath, nil, 0)
	if err != nil {
		return nil, err
	}
	for _, pkg := range pkgs {
		var files []*ast.File
		for _, f := range pkg.Files {
			files = append(files, f)
		}
		found = append(found, findImplementationInFiles(files, interfaceList)...)
	}

	return found, nil
}

func findImplementationInFiles(files []*ast.File, interfaceList []string) (found []string) {
	// collect all structs under path to find out the ones that satisfies the implementation
	structImplementations := make(map[string]implementation)

	for _, f := range files {
		ast.Inspect(f, func(n ast.Node) bool {
			// look for struct methods.
			methodDecl, ok := n.(*ast.FuncDecl)
			if !ok {
				return true
			}

			// not a method.
			if methodDecl.Recv == nil {
				return true
			}

			methodName := methodDecl.Name.Name

			// find the struct name that method belongs to.
			t := methodDecl.Recv.List[0].Type
			var ident *ast.Ident
			switch t := t.(type) {
			case *ast.Ident:
				// method with a value receiver
				ident = t
			case *ast.IndexExpr:
				// generic method with a value receiver
				ident = t.X.(*ast.Ident)
			case *ast.StarExpr:
				switch t := t.X.(type) {
				case *ast.Ident:
					// method with a pointer receiver
					ident = t
				case *ast.IndexExpr:
					// generic method with a pointer receiver
					ident = t.X.(*ast.Ident)
				default:
					return true
				}
			default:
				return true
			}
			structName := ident.Name

			// mark the implementation that this struct satisfies.
			if _, ok := structImplementations[structName]; !ok {
				structImplementations[structName] = newImplementation(interfaceList)
			}

			structImplementations[structName][methodName] = true

			return true
		})
	}

	for name, impl := range structImplementations {
		if checkImplementation(impl) {
			found = append(found, name)
		}
	}

	return found
}

// newImplementation returns a new object to parse implementation of an interface.
func newImplementation(interfaceList []string) implementation {
	impl := make(implementation)
	for _, m := range interfaceList {
		impl[m] = false
	}
	return impl
}

// checkImplementation checks if the entire implementation is satisfied.
func checkImplementation(r implementation) bool {
	for _, ok := range r {
		if !ok {
			return false
		}
	}
	return true
}

// ErrPathNotChain is returned by IsChainPath() when path is not a chain path.
type ErrPathNotChain struct {
	path string
	err  error
}

func (e ErrPathNotChain) Error() string {
	return fmt.Sprintf("%s not a chain path: %v", e.path, e.err)
}

// IsChainPath returns nil if path contains a cosmos chain.
func IsChainPath(path string) error {
	errf := func(err error) error {
		return ErrPathNotChain{path: path, err: err}
	}
	modFile, err := gomodule.ParseAt(path)
	if err != nil {
		return errf(err)
	}
	if err := ValidateGoMod(modFile); err != nil {
		return errf(err)
	}
	return nil
}

// ValidateGoMod check if the cosmos-sdk and the tendermint packages are imported.
func ValidateGoMod(module *modfile.File) error {
	moduleCheck := map[string]bool{
		cosmosModulePath:     true,
		tendermintModulePath: true,
	}
	for _, r := range module.Require {
		delete(moduleCheck, r.Mod.Path)
	}
	for m := range moduleCheck {
		return fmt.Errorf("invalid go module, missing %s package dependency", m)
	}
	return nil
}

// FindAppFilePath looks for the app file that implements the interfaces listed in appImplementation.
func FindAppFilePath(chainRoot string) (path string, err error) {
	var found []string

	err = filepath.Walk(chainRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(info.Name()) != ".go" {
			return nil
		}

		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return err
		}

		currFound := findImplementationInFiles([]*ast.File{f}, appImplementation)

		if len(currFound) > 0 {
			found = append(found, path)
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	numFound := len(found)
	if numFound == 0 {
		return "", errors.New("app.go file cannot be found")
	}

	if numFound == 1 {
		return found[0], nil
	}

	appFilePath := ""
	for _, p := range found {
		if filepath.Base(p) == appFileName {
			if appFilePath != "" {
				// multiple app.go found, fallback to app/app.go
				return getDefaultAppFile(chainRoot)
			}

			appFilePath = p
		}
	}

	if appFilePath != "" {
		return appFilePath, nil
	}

	return getDefaultAppFile(chainRoot)
}

// getDefaultAppFile returns the default app.go file path for a chain.
func getDefaultAppFile(chainRoot string) (string, error) {
	path := filepath.Join(chainRoot, defaultAppFilePath)
	_, err := os.Stat(path)
	return path, errors.Wrap(err, "cannot locate your app.go")
}
