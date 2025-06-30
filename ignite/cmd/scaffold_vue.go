package ignitecmd

import (
	"path/filepath"

	"github.com/spf13/cobra"

	chainconfig "github.com/ignite/cli/v29/ignite/config/chain"
	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosgen"
)

// NewScaffoldVue scaffolds a Vue.js app for a chain.
func NewScaffoldVue() *cobra.Command {
	c := &cobra.Command{
		Use:     "vue",
		Short:   "Vue 3 web app template",
		Args:    cobra.NoArgs,
		PreRunE: migrationPreRunHandler,
		RunE:    scaffoldVueHandler,
	}

	c.Flags().AddFlagSet(flagSetYes())

	return c
}

func scaffoldVueHandler(cmd *cobra.Command, _ []string) error {
	session := cliui.New(
		cliui.StartSpinnerWithText(statusScaffolding),
		cliui.WithoutUserInteraction(getYes(cmd)),
	)
	defer session.End()

	path := filepath.Join(".", chainconfig.DefaultVuePath)
	if err := cosmosgen.Vue(path); err != nil {
		return err
	}

	return session.Printf("\n🎉 Scaffolded a Vue.js app in %s.\n\n", path)
}
