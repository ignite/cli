package starportcmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/services/scaffolder"
)

// NewScaffoldVue scaffolds a Vue.js app for a chain.
func NewScaffoldVue() *cobra.Command {
	c := &cobra.Command{
		Use:   "vue",
		Short: "Vue 3 web app template",
		Args:  cobra.NoArgs,
		RunE:  scaffoldVueHandler,
	}
	
	return c
}

func scaffoldVueHandler(cmd *cobra.Command, args []string) error {
	s := clispinner.New().SetText("Scaffolding...")
	defer s.Stop()

	path := filepath.Join(flagGetAppPath(cmd), "vue")
	if err := scaffolder.Vue(path); err != nil {
		return err
	}

	s.Stop()
	fmt.Printf("\n🎉 Scaffold a Vue.js app.\n\n")

	return nil
}
