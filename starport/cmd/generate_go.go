package starportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/services/chain"
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

	c, err := newChainWithHomeFlags(cmd, appPath)
	if err != nil {
		return err
	}

	if err := c.Generate(cmd.Context(), chain.GenerateGo()); err != nil {
		return err
	}

	s.Stop()
	fmt.Println("⛏️  Generated go code.")

	return nil
}
