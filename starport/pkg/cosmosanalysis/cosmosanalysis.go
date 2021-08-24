// Package cosmosanalysis provides a toolset for staticly analysing Cosmos SDK's
// source code and blockchain source codes based on the Cosmos SDK
package cosmosanalysis

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

// implementation tracks the implementation of an interface for a given struct
type implementation map[string]bool

// FindImplementation finds the name of all types that implement the provided interface
func FindImplementation(modulePath string, interfaceList []string) (found []string, err error) {
	// parse go packages/files under path
	fset := token.NewFileSet()

	// collect all structs under path to find out the ones that satisfies the implementation
	structImplementations := make(map[string]implementation)
	pkgs, err := parser.ParseDir(fset, modulePath, nil, 0)
	if err != nil {
		return nil, err
	}
	for _, pkg := range pkgs {
		for _, f := range pkg.Files {
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
				ident, ok := t.(*ast.Ident)
				if !ok {
					sexp, ok := t.(*ast.StarExpr)
					if !ok {
						return true
					}
					ident = sexp.X.(*ast.Ident)
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
	}

	// append structs that satisfy the implementation
	for name, impl := range structImplementations {
		if checkImplementation(impl) {
			found = append(found, name)
		}
	}

	return found, nil
}

const (
	TypeString = "string"
	TypeBool   = "bool"
	TypeInt32  = "int32"
	TypeInt    = "int"
	TypeUint64 = "uint64"
	TypeUint   = "uint"
)

var (
	staticTypes = map[string]string{
		TypeString: TypeString,
		TypeBool:   TypeBool,
		TypeInt32:  TypeInt,
		TypeUint64: TypeUint,
	}
)

// FindStructFields finds the fields from a specific struct declared into the app
func FindStructFields(modulePath, structName string) ([]string, error) {
	// Parse go packages/files under path
	fset := token.NewFileSet()
	fields := make([]string, 0)
	pkgs, err := parser.ParseDir(fset, modulePath, nil, 0)
	if err != nil {
		return nil, err
	}

	found := false
	for _, pkg := range pkgs {
		for _, f := range pkg.Files {
			ast.Inspect(f, func(n ast.Node) bool {
				// Look for struct methods.
				structSpec, ok := n.(*ast.TypeSpec)
				if !ok {
					return true
				}
				specName := structSpec.Name.Name

				// Check if it is the struct we want
				structType, ok := structSpec.Type.(*ast.StructType)
				if !ok || structName != specName {
					return true
				}
				found = true

				// Collect all struct fields
				for _, field := range structType.Fields.List {
					i, ok := field.Type.(*ast.Ident)
					if !ok || len(field.Names) == 0 {
						continue
					}
					fieldType := i.Name
					if fieldType, ok := staticTypes[fieldType]; ok {
						for _, fieldName := range field.Names {
							fieldString := fmt.Sprintf("%s:%s", fieldName.Name, fieldType)
							fields = append(fields, fieldString)
						}
						continue
					}

					// Find nested custom fields
					customFields, err := FindStructFields(modulePath, fieldType)
					if err != nil {
						return false
					}
					fields = append(fields, customFields...)
				}
				return true
			})
		}
	}
	if !found {
		return fields, fmt.Errorf("struct '%s' not found", structName)
	}
	return fields, err
}

// newImplementation returns a new object to parse implementation of an interface
func newImplementation(interfaceList []string) implementation {
	impl := make(implementation)
	for _, m := range interfaceList {
		impl[m] = false
	}
	return impl
}

// checkImplementation checks if the entire implementation is satisfied
func checkImplementation(r implementation) bool {
	for _, ok := range r {
		if !ok {
			return false
		}
	}
	return true
}
