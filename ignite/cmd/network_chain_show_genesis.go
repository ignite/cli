package ignitecmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ignite-hq/cli/ignite/pkg/cliui"
	"github.com/ignite-hq/cli/ignite/pkg/cliui/icons"
	"github.com/ignite-hq/cli/ignite/pkg/cosmosutil"
	"github.com/ignite-hq/cli/ignite/services/network/networkchain"
)

func newNetworkChainShowGenesis() *cobra.Command {
	c := &cobra.Command{
		Use:   "genesis [launch-id]",
		Short: "Show the chain genesis file",
		Args:  cobra.ExactArgs(1),
		RunE:  networkChainShowGenesisHandler,
	}

	c.Flags().String(flagOut, "./genesis.json", "Path to output Genesis file")
	c.Flags().String(flagSPNChainID, cosmosutil.SPNChainID, "Chain ID to use for this network")

	return c
}

func networkChainShowGenesisHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New()
	defer session.Cleanup()

	var (
		out, _     = cmd.Flags().GetString(flagOut)
		chainID, _ = cmd.Flags().GetString(flagSPNChainID)
	)

	nb, launchID, err := networkChainLaunch(cmd, args, session)
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

	genesisPath, err := c.GenesisPath()
	if err != nil {
		return err
	}

	// check if the genesis already exists
	if _, err = os.Stat(genesisPath); os.IsNotExist(err) {
		// fetch the information to construct genesis
		genesisInformation, err := n.GenesisInformation(cmd.Context(), launchID)
		if err != nil {
			return err
		}

		// create the chain in a temp dir
		tmpHome, err := os.MkdirTemp("", "*-spn")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tmpHome)

		c.SetHome(tmpHome)

		rewardInfo, lastBlockHeight, unboundingTime, err := n.RewardsInfo(
			cmd.Context(),
			launchID,
			chainLaunch.ConsumerRevisionHeight,
		)
		if err != nil {
			return err
		}

		if err = c.Prepare(
			cmd.Context(),
			genesisInformation,
			rewardInfo,
			chainID,
			lastBlockHeight,
			unboundingTime,
		); err != nil {
			return err
		}

		// get the new genesis path
		genesisPath, err = c.GenesisPath()
		if err != nil {
			return err
		}
	}

	if err := os.MkdirAll(filepath.Dir(out), 0744); err != nil {
		return err
	}

	if err := os.Rename(genesisPath, out); err != nil {
		return err
	}

	session.StopSpinner()

	return session.Printf("%s Genesis generated: %s\n", icons.Bullet, out)
}
