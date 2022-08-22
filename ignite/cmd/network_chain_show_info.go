package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/cosmosutil"
	"github.com/ignite/cli/ignite/pkg/yaml"
	"github.com/ignite/cli/ignite/services/network"
	"github.com/ignite/cli/ignite/services/network/networktypes"
)

func newNetworkChainShowInfo() *cobra.Command {
	c := &cobra.Command{
		Use:   "info [launch-id]",
		Short: "Show info details of the chain",
		Args:  cobra.ExactArgs(1),
		RunE:  networkChainShowInfoHandler,
	}
	return c
}

func networkChainShowInfoHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New()
	defer session.Cleanup()

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

	reward, err := n.ChainReward(cmd.Context(), launchID)
	if err != nil && err != network.ErrObjectNotFound {
		return err
	}
	chainLaunch.Reward = reward.RemainingCoins.String()

	var genesis []byte
	if chainLaunch.GenesisURL != "" {
		genesis, _, err = cosmosutil.GenesisAndHashFromURL(cmd.Context(), chainLaunch.GenesisURL)
		if err != nil {
			return err
		}
	}
	chainInfo := struct {
		Chain   networktypes.ChainLaunch `json:"Chain"`
		Genesis []byte                   `json:"Genesis"`
	}{
		Chain:   chainLaunch,
		Genesis: genesis,
	}
	info, err := yaml.Marshal(cmd.Context(), chainInfo, "$.Genesis")
	if err != nil {
		return err
	}

	session.StopSpinner()

	return session.Print(info)
}
