package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/services/scaffolder"
)

// NewScaffoldFlutter scaffolds a Flutter app for a chain.
func NewScaffoldFlutter() *cobra.Command {
	c := &cobra.Command{
		Use:     "flutter",
		Short:   "A Flutter app for your chain",
		Args:    cobra.NoArgs,
		PreRunE: gitChangesConfirmPreRunHandler,
		RunE:    scaffoldFlutterHandler,
	}

	c.Flags().AddFlagSet(flagSetYes())
	c.Flags().StringP(flagPath, "p", "./flutter", "path to scaffold content of the Flutter app")

	return c
}

func scaffoldFlutterHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New(cliui.StartSpinner())
	defer session.End()

	session.StartSpinner("Scaffolding...")

	path := flagGetPath(cmd)
	if err := scaffolder.Flutter(path); err != nil {
		return err
	}

	return session.Printf("\n🎉 Scaffold a Flutter app.\n\n")
}
