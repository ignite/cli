package ignitecmd

import (
	"fmt"

	"github.com/spf13/cobra"
	launchtypes "github.com/tendermint/spn/x/launch/types"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/cosmosutil"
	"github.com/ignite/cli/ignite/services/network"
	"github.com/ignite/cli/ignite/services/network/networktypes"
)

var requestSummaryHeader = []string{"ID", "Status", "Type", "Content"}

// NewNetworkRequestList creates a new request list command to list
// requests for a chain
func NewNetworkRequestList() *cobra.Command {
	c := &cobra.Command{
		Use:   "list [launch-id]",
		Short: "List all pending requests",
		RunE:  networkRequestListHandler,
		Args:  cobra.ExactArgs(1),
	}

	c.Flags().AddFlagSet(flagSetSPNAccountPrefixes())

	return c
}

func networkRequestListHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New()
	defer session.Cleanup()

	nb, err := newNetworkBuilder(cmd, CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}

	addressPrefix := getAddressPrefix(cmd)

	// parse launch ID
	launchID, err := network.ParseID(args[0])
	if err != nil {
		return err
	}

	n, err := nb.Network()
	if err != nil {
		return err
	}

	requests, err := n.Requests(cmd.Context(), launchID)
	if err != nil {
		return err
	}

	session.StopSpinner()

	return renderRequestSummaries(requests, session, addressPrefix)
}

// renderRequestSummaries writes into the provided out, the list of summarized requests
func renderRequestSummaries(
	requests []networktypes.Request,
	session cliui.Session,
	addressPrefix string,
) error {
	requestEntries := make([][]string, 0)
	for _, request := range requests {
		var (
			id          = fmt.Sprintf("%d", request.RequestID)
			requestType = "Unknown"
			content     = ""
		)
		switch req := request.Content.Content.(type) {
		case *launchtypes.RequestContent_GenesisAccount:
			requestType = "Add Genesis Account"

			address, err := cosmosutil.ChangeAddressPrefix(
				req.GenesisAccount.Address,
				addressPrefix,
			)
			if err != nil {
				return err
			}

			content = fmt.Sprintf("%s, %s",
				address,
				req.GenesisAccount.Coins.String())
		case *launchtypes.RequestContent_GenesisValidator:
			requestType = "Add Genesis Validator"
			peer, err := network.PeerAddress(req.GenesisValidator.Peer)
			if err != nil {
				return err
			}

			address, err := cosmosutil.ChangeAddressPrefix(
				req.GenesisValidator.Address,
				addressPrefix,
			)
			if err != nil {
				return err
			}

			content = fmt.Sprintf("%s, %s, %s",
				peer,
				address,
				req.GenesisValidator.SelfDelegation.String())
		case *launchtypes.RequestContent_VestingAccount:
			requestType = "Add Vesting Account"

			// parse vesting options
			var vestingCoins string
			dv := req.VestingAccount.VestingOptions.GetDelayedVesting()
			if dv == nil {
				vestingCoins = "unrecognized vesting option"
			} else {
				vestingCoins = fmt.Sprintf("%s (vesting: %s)", dv.TotalBalance, dv.Vesting)
			}

			address, err := cosmosutil.ChangeAddressPrefix(
				req.VestingAccount.Address,
				addressPrefix,
			)
			if err != nil {
				return err
			}

			content = fmt.Sprintf("%s, %s",
				address,
				vestingCoins,
			)
		case *launchtypes.RequestContent_ValidatorRemoval:
			requestType = "Remove Validator"

			address, err := cosmosutil.ChangeAddressPrefix(
				req.ValidatorRemoval.ValAddress,
				addressPrefix,
			)
			if err != nil {
				return err
			}

			content = address
		case *launchtypes.RequestContent_AccountRemoval:
			requestType = "Remove Account"

			address, err := cosmosutil.ChangeAddressPrefix(
				req.AccountRemoval.Address,
				addressPrefix,
			)
			if err != nil {
				return err
			}

			content = address
		}

		requestEntries = append(requestEntries, []string{
			id,
			request.Status,
			requestType,
			content,
		})
	}
	return session.PrintTable(requestSummaryHeader, requestEntries...)
}
