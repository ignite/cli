package cmd

import (
	"context"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/templates/app"
	"golang.org/x/mod/module"
)

var appCmd = &cobra.Command{
	Use:   "app [github.com/org/repo]",
	Short: "Generates an empty application",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fullName := args[0]
		var appName string
		if t := strings.Split(fullName, "/"); len(t) > 0 {
			appName = t[len(t)-1]
			// app name cannot contain "-" so gracefully remove them
			// if they present.
			appName = strings.ReplaceAll(appName, "-", "")
		}
		if err := validateGoModuleName(fullName); err != nil {
			return err
		}
		if err := validateGoPkgName(appName); err != nil {
			return err
		}
		denom, _ := cmd.Flags().GetString("denom")
		g, _ := app.New(&app.Options{
			ModulePath: fullName,
			AppName:    appName,
			Denom:      denom,
		})
		run := genny.WetRunner(context.Background())
		run.With(g)
		pwd, _ := os.Getwd()
		run.Root = pwd + "/" + appName
		run.Run()
		message := `
⭐️ Successfully created a Cosmos app '%[1]v'.
👉 Get started with the following commands:

 %% cd %[1]v
 %% starport serve

NOTE: add -v flag for advanced use.
`
		fmt.Printf(message, appName)
		return nil
	},
}

func validateGoPkgName(name string) error {
	fset := token.NewFileSet()
	src := fmt.Sprintf("package %s", name)
	if _, err := parser.ParseFile(fset, "", src, parser.PackageClauseOnly); err != nil {
		// parser error is very low level here so let's hide it from the user
		// completely.
		return errors.New("app name is an invalid go package name")
	}
	return nil
}

func validateGoModuleName(name string) error {
	err := module.CheckPath(name)
	return errors.Wrap(err, "app name is an invalid go module name")
}
