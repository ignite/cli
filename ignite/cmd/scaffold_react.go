package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/services/scaffolder"
)

// NewScaffoldVue scaffolds a Vue.js app for a chain.
func NewScaffoldReact() *cobra.Command {
	c := &cobra.Command{
		Use:     "react",
		Short:   "React web app template",
		Args:    cobra.NoArgs,
		PreRunE: gitChangesConfirmPreRunHandler,
		RunE:    scaffoldReactHandler,
	}

	c.Flags().AddFlagSet(flagSetYes())
	c.Flags().StringP(flagPath, "p", "./react", "path to scaffold content of the ReactJS app")

	return c
}

func scaffoldReactHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New(cliui.StartSpinner())
	defer session.End()

	session.StartSpinner("Scaffolding...")

	path := flagGetPath(cmd)
	if err := scaffolder.React(path); err != nil {
		return err
	}

	return session.Printf("\n🎉 Scaffold a ReactJS app.\n\n")
}
