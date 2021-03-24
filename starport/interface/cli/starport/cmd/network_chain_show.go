package starportcmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/tendermint/starport/starport/pkg/spn"
	"github.com/tendermint/starport/starport/services/networkbuilder"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
)

const (
	genesisFlag = "genesis"
	peersFlag   = "peers"
)

// NewNetworkChainShow creates a new chain show command to show
// a chain on SPN.
func NewNetworkChainShow() *cobra.Command {
	c := &cobra.Command{
		Use:   "show [chain-id]",
		Short: "Show details of a chain",
		RunE:  networkChainShowHandler,
		Args:  cobra.ExactArgs(1),
	}
	c.Flags().Bool(genesisFlag, false, "Show exclusively the genesis of the chain")
	c.Flags().Bool(peersFlag, false, "Show exclusively the peers of the chain")
	return c
}

func networkChainShowHandler(cmd *cobra.Command, args []string) error {
	chainID := args[0]

	nb, err := newNetworkBuilder()
	if err != nil {
		return err
	}

	// Get flags
	showGenesis, err := cmd.Flags().GetBool(genesisFlag)
	if err != nil {
		return err
	}
	showPeers, err := cmd.Flags().GetBool(peersFlag)
	if err != nil {
		return err
	}
	if showGenesis && showPeers {
		return fmt.Errorf("%s and %s flags cannot be used together", genesisFlag, peersFlag)
	}

	// Fetch launch information
	info, err := nb.LaunchInformation(cmd.Context(), chainID)
	if err != nil {
		return err
	}

	switch {
	case showGenesis:
		if err := generateAndShowGenesis(cmd.Context(), nb, chainID, info); err != nil {
			return err
		}
	case showPeers:
		// Show the peers in the config.toml format
		fmt.Printf("persistent_peers = \"%s\"\n", strings.Join(info.Peers, ","))
	default:
		// No flag, show the chain information and launch information in yaml format
		chain, err := nb.ShowChain(cmd.Context(), chainID)
		if err != nil {
			return err
		}
		chainyaml, err := yaml.Marshal(chain)
		if err != nil {
			return err
		}
		infoyaml, err := yaml.Marshal(info)
		if err != nil {
			return err
		}
		fmt.Printf("\nChain:\n---\n%s\n\nLaunch Information:\n---\n%s", string(chainyaml), string(infoyaml))
	}

	return nil
}

func generateAndShowGenesis(ctx context.Context, nb *networkbuilder.Builder, chainID string, info spn.LaunchInformation) error {
	tmpHome, err := os.MkdirTemp("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpHome)

	// Initialize the blockchain
	blockchain, err := nb.Init(
		ctx,
		chainID,
		networkbuilder.SourceChainID(),
		networkbuilder.InitializationHomePath(tmpHome),
	)
	if err != nil {
		return err
	}
	defer blockchain.Cleanup()

	// Generate the genesis in a temporary directory and show the content
	genesisPath, err := nb.GenerateGenesisWithHome(ctx, chainID, info, tmpHome)
	if err != nil {
		return err
	}
	genesis, err := os.ReadFile(genesisPath)
	if err != nil {
		return err
	}
	fmt.Print(string(genesis))

	return nil
}
