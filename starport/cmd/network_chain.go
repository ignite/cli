package starportcmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/services/network"
)

// NewNetworkChain creates a new chain command that holds some other
// sub commands related to launching a network for a chain.
func NewNetworkChain() *cobra.Command {
	c := &cobra.Command{
		Use:   "chain",
		Short: "Build networks",
	}

	c.AddCommand(NewNetworkChainPublish())

	return c
}

// checkChainHomeExist checks if a home with the provided launchID already exist
func checkChainHomeExist(launchID uint64) (string, bool, error) {
	home, err := network.ChainHome(launchID)
	if err != nil {
		return home, false, err
	}

	if _, err := os.Stat(home); os.IsNotExist(err) {
		return home, false, nil
	}
	return home, true, err
}
