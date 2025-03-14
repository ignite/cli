package ignitecmd

import (
	"github.com/spf13/cobra"

	chainconfig "github.com/ignite/cli/v29/ignite/config/chain"
	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosgen"
)

// NewScaffoldReact scaffolds a React app for a chain.
func NewScaffoldReact() *cobra.Command {
	c := &cobra.Command{
		Hidden:  true, // hidden util we have a better ts-client.
		Use:     "react",
		Short:   "React web app template",
		Args:    cobra.NoArgs,
		PreRunE: migrationPreRunHandler,
		RunE:    scaffoldReactHandler,
	}

	c.Flags().AddFlagSet(flagSetYes())
	c.Flags().StringP(flagPath, "p", "./"+chainconfig.DefaultReactPath, "path to scaffold content of the React app")

	return c
}

func scaffoldReactHandler(cmd *cobra.Command, _ []string) error {
	session := cliui.New(cliui.StartSpinnerWithText(statusScaffolding))
	defer session.End()

	path := flagGetPath(cmd)
	if err := cosmosgen.React(path); err != nil {
		return err
	}

	return session.Printf("\n🎉 Scaffolded a React app in %s.\n\n", path)
}
