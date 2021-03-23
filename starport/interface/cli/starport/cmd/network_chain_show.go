package starportcmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

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
	c.Flags().AddFlagSet(flagSetHomes())
	c.Flags().Bool(genesisFlag, false, "Show exclusively the genesis of the chain")
	c.Flags().Bool(peersFlag, false, "Show exclusively the peers of the chain")
	return c
}

func networkChainShowHandler(cmd *cobra.Command, args []string) error {
	nb, err := newNetworkBuilder()
	if err != nil {
		return err
	}

	chainID := args[0]

	// Get flags
	home, _, err := getHomeFlags(cmd)
	if err != nil {
		return err
	}
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
		// Generate the genesis in a temporary directory and show the content
		tmpHome, err := nb.GenerateTemporaryGenesis(cmd.Context(), chainID, home, info)
		defer os.RemoveAll(string(tmpHome))
		if err != nil {
			return err
		}
		genesis, err := ioutil.ReadFile(tmpHome.GenesisPath())
		if err != nil {
			return err
		}
		fmt.Print(string(genesis))
	case showPeers:
		// Show the peers in the config.toml format
		fmt.Printf(`persistent_peers = "%s"`, strings.Join(info.Peers, ","))
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
