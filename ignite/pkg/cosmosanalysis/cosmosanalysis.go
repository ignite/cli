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
	"strings"

	"golang.org/x/mod/modfile"

	"github.com/ignite/cli/v29/ignite/pkg/cosmosver"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/gomodule"
)

const (
	tendermintModulePath = "github.com/cometbft/cometbft"
	appFileName          = "app.go"
	defaultAppFilePath   = "app/" + appFileName
)

var AppEmbeddedTypes = []string{
	"github.com/cosmos/cosmos-sdk/runtime.App",
	"github.com/cosmos/cosmos-sdk/baseapp.BaseApp",
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

// FindImplementationInFile find all struct implements the interfaceList into an ast.File.
func FindImplementationInFile(n ast.Node, interfaceList []string) (found []string) {
	// collect all structs under path to find out the ones that satisfies the implementation
	structImplementations := make(map[string]implementation)

	findImplementation(n, func(methodName, structName string) bool {
		// mark the implementation that this struct satisfies.
		if _, ok := structImplementations[structName]; !ok {
			structImplementations[structName] = newImplementation(interfaceList)
		}

		structImplementations[structName][methodName] = true

		return true
	})

	for name, impl := range structImplementations {
		if checkImplementation(impl) {
			found = append(found, name)
		}
	}

	return found
}

// findImplementationInFiles find all struct implements the interfaceList into a list of ast.File.
func findImplementationInFiles(files []*ast.File, interfaceList []string) (found []string) {
	// collect all structs under path to find out the ones that satisfies the implementation
	structImplementations := make(map[string]implementation)

	for _, f := range files {
		findImplementation(f, func(methodName, structName string) bool {
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

// findImplementation parse the ast.Node and call the callback if is a struct implementation.
func findImplementation(f ast.Node, endCallback func(methodName, structName string) bool) {
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

		if endCallback != nil {
			return endCallback(methodName, structName)
		}
		return true
	})
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

// FindEmbed finds the name of all types that embed one of the target types in a given module path.
// targetEmbeddedTypes should be a list of fully qualified type names (e.g., "package/path.TypeName").
func FindEmbed(modulePath string, targetEmbeddedTypes []string) (found []string, err error) {
	// parse go packages/files under path
	fset := token.NewFileSet()

	pkgs, err := parser.ParseDir(fset, modulePath, nil, 0)
	if err != nil {
		return nil, err
	}

	for _, pkg := range pkgs {
		for _, fileNode := range pkg.Files {
			foundStructs := findStructsEmbeddingInFile(fileNode, targetEmbeddedTypes)
			found = append(found, foundStructs...)
		}
	}

	// Deduplicate results as a struct might be found in multiple files of the same package (though unlikely for structs)
	// or if the same struct name exists in different packages (FindEmbed currently doesn't qualify by package).
	if len(found) > 0 {
		uniqueNamesMap := make(map[string]struct{})
		var uniqueResult []string
		for _, name := range found {
			if _, exists := uniqueNamesMap[name]; !exists {
				uniqueNamesMap[name] = struct{}{}
				uniqueResult = append(uniqueResult, name)
			}
		}
		return uniqueResult, nil
	}

	return found, nil
}

// FindEmbedInFile finds all struct names in a given AST node that embed one of the target types.
// The AST node is expected to be an *ast.File.
// targetEmbeddedTypes should be a list of fully qualified type names (e.g., "package/path.TypeName").
func FindEmbedInFile(n ast.Node, targetEmbeddedTypes []string) (found []string) {
	fileNode, ok := n.(*ast.File)
	if !ok {
		return nil
	}

	return findStructsEmbeddingInFile(fileNode, targetEmbeddedTypes)
}

// findStructsEmbeddingInFile checks if any struct in the given AST file embeds one of the target types.
// targetTypes should be fully qualified (e.g., "package/path.TypeName").
func findStructsEmbeddingInFile(fileNode *ast.File, targetEmbeddedTypes []string) (foundStructNames []string) {
	// activeTargets maps local package name to a set of expected TypeNames from that package
	activeTargets := make(map[string]map[string]struct{})

	for _, targetFQN := range targetEmbeddedTypes {
		dotIndex := strings.LastIndex(targetFQN, ".")
		if dotIndex == -1 || dotIndex == 0 || dotIndex == len(targetFQN)-1 {
			continue // invalid format
		}
		expectedImportPath := targetFQN[:dotIndex]
		expectedTypeName := targetFQN[dotIndex+1:]

		for _, imp := range fileNode.Imports {
			importPath := strings.Trim(imp.Path.Value, `"`)
			if importPath == expectedImportPath {
				localPkgName := ""
				if imp.Name != nil { // alias used
					localPkgName = imp.Name.Name
				} else {
					// default name (last part of the path)
					// this is a common heuristic, e.g. "github.com/cosmos/cosmos-sdk/runtime" -> "runtime"
					pathParts := strings.Split(importPath, "/")
					localPkgName = pathParts[len(pathParts)-1]
				}

				if _, ok := activeTargets[localPkgName]; !ok {
					activeTargets[localPkgName] = make(map[string]struct{})
				}
				activeTargets[localPkgName][expectedTypeName] = struct{}{}
				break // found the import for this target, move to next targetFQN
			}
		}
	}

	if len(activeTargets) == 0 {
		return nil // none of the target packages are imported in this file
	}

	ast.Inspect(fileNode, func(n ast.Node) bool {
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}

		for _, field := range structType.Fields.List {
			if len(field.Names) == 0 { // embedded field
				var selExpr *ast.SelectorExpr
				fieldType := field.Type

				if starExpr, isStar := fieldType.(*ast.StarExpr); isStar {
					fieldType = starExpr.X // unwrap pointer
				}

				if se, isSel := fieldType.(*ast.SelectorExpr); isSel {
					selExpr = se
				} else {
					continue
				}

				pkgIdent, okIdent := selExpr.X.(*ast.Ident)
				if !okIdent {
					continue
				}

				pkgNameInCode := pkgIdent.Name
				typeNameInCode := selExpr.Sel.Name

				if expectedTypeNamesSet, pkgFound := activeTargets[pkgNameInCode]; pkgFound {
					if _, typeFound := expectedTypeNamesSet[typeNameInCode]; typeFound {
						foundStructNames = append(foundStructNames, typeSpec.Name.Name)
					}
				}
			}
		}
		return true
	})

	// deduplicate if a struct somehow embeds multiple (or the same) target type
	if len(foundStructNames) > 0 {
		uniqueNamesMap := make(map[string]struct{})
		var uniqueResult []string
		for _, name := range foundStructNames {
			if _, exists := uniqueNamesMap[name]; !exists {
				uniqueNamesMap[name] = struct{}{}
				uniqueResult = append(uniqueResult, name)
			}
		}
		return uniqueResult
	}

	return foundStructNames
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
		cosmosver.CosmosModulePath: true,
		tendermintModulePath:       true,
	}

	for _, r := range module.Require {
		delete(moduleCheck, r.Mod.Path)
	}
	for m := range moduleCheck {
		return errors.Errorf("invalid go module, missing %s package dependency", m)
	}
	return nil
}

// FindAppFilePath Looks for the app file that embeds the runtime.App or baseapp.BaseApp types.
func FindAppFilePath(chainRoot string) (path string, err error) {
	var foundAppStructFiles []string
	err = filepath.Walk(chainRoot, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(info.Name()) != ".go" {
			return nil
		}

		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, currentPath, nil, 0)
		if err != nil {
			// log or handle error, e.g. by returning nil to continue walking
			return nil
		}

		structNames := findStructsEmbeddingInFile(f, AppEmbeddedTypes)
		if len(structNames) > 0 {
			foundAppStructFiles = append(foundAppStructFiles, currentPath)
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	numFound := len(foundAppStructFiles)
	if numFound == 0 {
		return "", errors.New("app.go file cannot be found")
	}

	if numFound == 1 {
		return foundAppStructFiles[0], nil
	}

	// multiple files found, prefer one named appFileName ("app.go")
	appFilePath := ""
	for _, p := range foundAppStructFiles {
		if filepath.Base(p) == appFileName {
			if appFilePath != "" {
				// more than one app.go found among candidates, fallback to default
				return getDefaultAppFile(chainRoot)
			}
			appFilePath = p
		}
	}

	if appFilePath != "" {
		return appFilePath, nil
	}

	// no app.go found among the candidates, or multiple candidates and none are app.go,
	// fallback to default app path logic
	return getDefaultAppFile(chainRoot)
}

// getDefaultAppFile returns the default app.go file path for a chain.
func getDefaultAppFile(chainRoot string) (string, error) {
	path := filepath.Join(chainRoot, defaultAppFilePath)
	_, err := os.Stat(path)
	return path, errors.Wrap(err, "cannot locate your app.go")
}
