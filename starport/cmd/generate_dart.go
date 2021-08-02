package starportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/services/chain"
)

func NewGenerateDart() *cobra.Command {
	c := &cobra.Command{
		Use:   "dart",
		Short: "Generate a Dart client",
		RunE:  generateDartHandler,
	}
	c.Flags().AddFlagSet(flagSetProto3rdParty(""))
	return c
}

func generateDartHandler(cmd *cobra.Command, args []string) error {
	s := clispinner.New().SetText("Generating...")
	defer s.Stop()

	var chainOption []chain.Option

	if flagProto3rdParty(cmd) {
		chainOption = append(chainOption, chain.EnableThirdPartyModuleCodegen())
	}

	c, err := newChainWithHomeFlags(cmd, appPath, chainOption...)
	if err != nil {
		return err
	}

	if err := c.Generate(cmd.Context(), chain.GenerateDart()); err != nil {
		return err
	}

	s.Stop()
	fmt.Println("⛏️  Generated Dart client.")

	return nil
}
