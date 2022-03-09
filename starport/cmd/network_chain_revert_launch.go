package starportcmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/tendermint/starport/starport/pkg/cosmosutil"
	"github.com/tendermint/starport/starport/services/network/networkchain"

	"github.com/tendermint/starport/starport/services/network"
)

// NewNetworkChainRevertLaunch creates a new chain revert launch command
// to revert a launched chain.
func NewNetworkChainRevertLaunch() *cobra.Command {
	c := &cobra.Command{
		Use:   "revert-launch [launch-id]",
		Short: "Revert launch a network as a coordinator",
		Args:  cobra.ExactArgs(1),
		RunE:  networkChainRevertLaunchHandler,
	}

	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetKeyringBackend())

	return c
}

func networkChainRevertLaunchHandler(cmd *cobra.Command, args []string) error {
	nb, err := newNetworkBuilder(cmd)
	if err != nil {
		return err
	}
	defer nb.Cleanup()

	// parse launch ID
	launchID, err := network.ParseID(args[0])
	if err != nil {
		return err
	}

	n, err := nb.Network()
	if err != nil {
		return err
	}

	chainLaunch, err := n.ChainLaunch(cmd.Context(), launchID)
	if err != nil {
		return err
	}

	c, err := nb.Chain(networkchain.SourceLaunch(chainLaunch))
	if err != nil {
		return err
	}

	// set the genesis time for the chain
	genesisPath, err := c.GenesisPath()
	if err != nil {
		return errors.Wrap(err, "genesis of the blockchain can't be read")
	}
	if err := cosmosutil.SetGenesisTime(genesisPath, 0); err != nil {
		return errors.Wrap(err, "genesis time can't be set")
	}

	return n.RevertLaunch(cmd.Context(), launchID)
}
