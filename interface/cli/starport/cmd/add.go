package starportcmd

import (
	"context"

	"github.com/gobuffalo/genny"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/templates/add"
)

var addCmd = &cobra.Command{
	Use:   "add [feature]",
	Short: "Adds a feature to a project.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appName, _ := getAppAndModule(appPath)
		g, _ := add.New(&add.Options{
			Feature: args[0],
			AppName: appName,
		})
		run := genny.WetRunner(context.Background())
		run.With(g)
		run.Run()
		// fmt.Printf("\n🎉 Created a type `%[1]v`.\n\n", args[0])
	},
}
