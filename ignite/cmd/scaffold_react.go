package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/chainconfig"
	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/cosmosgen"
)

// NewScaffoldReact scaffolds a React app for a chain.
func NewScaffoldReact() *cobra.Command {
	c := &cobra.Command{
		Use:     "react",
		Short:   "Generate React web app template",
		Args:    cobra.NoArgs,
		PreRunE: gitChangesConfirmPreRunHandler,
		RunE:    scaffoldReactHandler,
	}

	c.Flags().AddFlagSet(flagSetYes())
	c.Flags().StringP(flagPath, "p", "./"+chainconfig.DefaultReactPath, "path to scaffold content of the React app")

	return c
}

func scaffoldReactHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New(cliui.StartSpinnerWithText(statusScaffolding))
	defer session.End()

	path := flagGetPath(cmd)
	if err := cosmosgen.React(path); err != nil {
		return err
	}

	return session.Printf("\n🎉 Scaffolded a React app in %s.\n\n", path)
}
