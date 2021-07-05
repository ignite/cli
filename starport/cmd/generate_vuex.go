package starportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/services/chain"
)

func NewGenerateVuex() *cobra.Command {
	c := &cobra.Command{
		Use:   "vuex",
		Short: "Generate Vuex store for you chain's frontend from your config.yml",
		RunE:  generateVuexHandler,
	}

	c.Flags().Bool(flagRebuildProtoOnce, false, "Enables proto code generation for 3rd party modules.")

	return c
}

func generateVuexHandler(cmd *cobra.Command, args []string) error {
	s := clispinner.New().SetText("Generating...")
	defer s.Stop()

	isRebuildProtoOnce, _ := cmd.Flags().GetBool(flagRebuildProtoOnce)

	chainOption := []chain.Option{}

	if isRebuildProtoOnce {
		chainOption = append(chainOption, chain.EnableThirdPartyModuleCodegen())
	}

	c, err := newChainWithHomeFlags(cmd, appPath, chainOption...)
	if err != nil {
		return err
	}

	if err := c.Generate(cmd.Context(), chain.GenerateVuex()); err != nil {
		return err
	}

	s.Stop()
	fmt.Println("⛏️  Generated vuex stores.")

	return nil
}
