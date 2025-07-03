package ignitecmd

import (
	"github.com/spf13/cobra"
<<<<<<< HEAD

	chainconfig "github.com/ignite/cli/v28/ignite/config/chain"
	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosgen"
=======
>>>>>>> d1bf508a (refactor!: remove react frontend + re-enable disabled integration tests (#4744))
)

// NewScaffoldReact scaffolds a React app for a chain.
func NewScaffoldReact() *cobra.Command {
	c := &cobra.Command{
		Use:        "react",
		Deprecated: "the React scaffolding feature is removed from Ignite CLI.\nPlease use the Ignite CCA app to create a React app.\nFor more information, visit: https://ignite.com/marketplace/CCA",
	}

	return c
}
<<<<<<< HEAD

func scaffoldReactHandler(cmd *cobra.Command, _ []string) error {
	session := cliui.New(cliui.StartSpinnerWithText(statusScaffolding))
	defer session.End()

	path := flagGetPath(cmd)
	if err := cosmosgen.React(path); err != nil {
		return err
	}

	return session.Printf("\n🎉 Scaffolded a React app in %s.\n\n", path)
}
=======
>>>>>>> d1bf508a (refactor!: remove react frontend + re-enable disabled integration tests (#4744))
