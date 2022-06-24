package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/ignite/pkg/cosmosutil"
	"github.com/ignite/cli/ignite/services/network"
)

var (
	chainGenesisValSummaryHeader = []string{"Genesis Validator", "Self Delegation", "Peer"}
)

func newNetworkChainShowValidators() *cobra.Command {
	c := &cobra.Command{
		Use:   "validators [launch-id]",
		Short: "Show all validators of the chain",
		Args:  cobra.ExactArgs(1),
		RunE:  networkChainShowValidatorsHandler,
	}

	c.Flags().AddFlagSet(flagSetSPNAccountPrefixes())

	return c
}

func networkChainShowValidatorsHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New()
	defer session.Cleanup()

	addressPrefix := getAddressPrefix(cmd)

	nb, launchID, err := networkChainLaunch(cmd, args, session)
	if err != nil {
		return err
	}
	n, err := nb.Network()
	if err != nil {
		return err
	}

	validators, err := n.GenesisValidators(cmd.Context(), launchID)
	if err != nil {
		return err
	}
	validatorEntries := make([][]string, 0)
	for _, acc := range validators {
		peer, err := network.PeerAddress(acc.Peer)
		if err != nil {
			return err
		}

		address, err := cosmosutil.ChangeAddressPrefix(acc.Address, addressPrefix)
		if err != nil {
			return err
		}

		validatorEntries = append(validatorEntries, []string{
			address,
			acc.SelfDelegation.String(),
			peer,
		})
	}
	if len(validatorEntries) == 0 {
		return session.Printf("%s %s\n", icons.Info, "no account found")
	}

	session.StopSpinner()

	return session.PrintTable(chainGenesisValSummaryHeader, validatorEntries...)
}
