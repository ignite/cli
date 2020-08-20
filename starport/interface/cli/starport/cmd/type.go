package starportcmd

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/templates/typed"
)

func NewType() *cobra.Command {
	c := &cobra.Command{
		Use:   "type [typeName] [field1] [field2] ...",
		Short: "Generates CRUD actions for type",
		Args:  cobra.MinimumNArgs(1),
		RunE:  typeHandler,
	}
	c.Flags().StringVarP(&appPath, "path", "p", "", "path of the app")
	return c
}

func typeHandler(cmd *cobra.Command, args []string) error {
	appName, modulePath := getAppAndModule(appPath)
	typeName := args[0]
	ok, err := isTypeCreated(appPath, appName, typeName)
	if err != nil {
		return err
	}
	if ok {
		return fmt.Errorf("%s type is already added.", typeName)
	}
	var fields []typed.Field
	for _, f := range args[1:] {
		fs := strings.Split(f, ":")
		name := fs[0]
		var datatype string
		acceptedTypes := map[string]bool{
			"string": true,
			"bool":   true,
			"int":    true,
			"float":  true,
		}
		if len(fs) == 2 && acceptedTypes[fs[1]] {
			datatype = fs[1]
		} else {
			datatype = "string"
		}
		field := typed.Field{Name: name, Datatype: datatype}
		fields = append(fields, field)
	}
	g, _ := typed.New(&typed.Options{
		ModulePath: modulePath,
		AppName:    appName,
		TypeName:   typeName,
		Fields:     fields,
	})
	run := genny.WetRunner(context.Background())
	run.With(g)
	run.Run()
	fmt.Printf("\n🎉 Created a type `%[1]v`.\n\n", args[0])
	return nil
}

func isTypeCreated(appPath, appName, typeName string) (isCreated bool, err error) {
	abspath, err := filepath.Abs(filepath.Join(appPath, "x", appName, "types"))
	if err != nil {
		return false, err
	}
	fset := token.NewFileSet()
	all, err := parser.ParseDir(fset, abspath, func(os.FileInfo) bool { return true }, parser.ParseComments)
	if err != nil {
		return false, err
	}
	for _, pkg := range all {
		for _, f := range pkg.Files {
			ast.Inspect(f, func(x ast.Node) bool {
				typeSpec, ok := x.(*ast.TypeSpec)
				if !ok {
					return true
				}
				if _, ok := typeSpec.Type.(*ast.StructType); !ok {
					return true
				}
				if strings.Title(typeName) != typeSpec.Name.Name {
					return true
				}
				isCreated = true
				return false
			})
		}
	}
	return
}
