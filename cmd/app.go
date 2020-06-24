package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/templates/app"
)

var appCmd = &cobra.Command{
	Use:   "app [github.com/org/repo]",
	Short: "Generates an empty application",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		denom, _ := cmd.Flags().GetString("denom")
		var appName string
		if t := strings.Split(args[0], "/"); len(t) > 0 {
			appName = t[len(t)-1]
		}
		g, _ := app.New(&app.Options{
			ModulePath: args[0],
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
	},
}
