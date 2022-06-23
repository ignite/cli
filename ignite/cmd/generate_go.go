package ignitecmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui/clispinner"
	"github.com/ignite/cli/ignite/services/chain"
)

func NewGenerateGo() *cobra.Command {
	return &cobra.Command{
		Use:   "proto-go",
		Short: "Generate proto based Go code needed for the app's source code",
		RunE:  generateGoHandler,
	}
}

func generateGoHandler(cmd *cobra.Command, args []string) error {
	s := clispinner.New().SetText("Generating...")
	defer s.Stop()

	c, err := newChainWithHomeFlags(cmd)
	if err != nil {
		return err
	}

	cacheStorage, err := newCache(cmd)
	if err != nil {
		return err
	}

	if err := c.Generate(cmd.Context(), cacheStorage, chain.GenerateGo()); err != nil {
		return err
	}

	s.Stop()
	fmt.Println("⛏️  Generated go code.")

	return nil
}
