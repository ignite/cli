package starportcmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	launchtypes "github.com/tendermint/spn/x/launch/types"

	"github.com/tendermint/starport/starport/pkg/entrywriter"
	"github.com/tendermint/starport/starport/services/network"
)

var requestSummaryHeader = []string{"ID", "Type", "Content"}

// NewNetworkRequestList creates a new request list command to list
// requests for a chain
func NewNetworkRequestList() *cobra.Command {
	c := &cobra.Command{
		Use:   "list [launch-id]",
		Short: "List all pending requests",
		RunE:  networkRequestListHandler,
		Args:  cobra.ExactArgs(1),
	}
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetHome())
	return c
}

func networkRequestListHandler(cmd *cobra.Command, args []string) error {
	// initialize network common methods
	nb, err := newNetworkBuilder(cmd)
	if err != nil {
		return err
	}

	// parse launch ID
	launchID, err := network.ParseLaunchID(args[0])
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

	nb.Cleanup()
	return renderRequestSummaries(requests, os.Stdout)
}

// renderRequestSummaries writes into the provided out, the list of summarized requests
func renderRequestSummaries(requests []launchtypes.Request, out io.Writer) error {
	requestEntries := make([][]string, 0)
	for _, request := range requests {
		id := fmt.Sprintf("%d", request.RequestID)
		requestType := "Unknown"
		content := ""

		switch req := request.Content.Content.(type) {
		case *launchtypes.RequestContent_GenesisAccount:
			requestType = "Add Genesis Account"
			content = fmt.Sprintf("%s, %s",
				req.GenesisAccount.Address,
				req.GenesisAccount.Coins.String())
		case *launchtypes.RequestContent_GenesisValidator:
			requestType = "Add Genesis Validator"
			peer, err := network.PeerAddress(req.GenesisValidator.Peer)
			if err != nil {
				return err
			}
			content = fmt.Sprintf("%s, %s, %s",
				peer,
				req.GenesisValidator.Address,
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
			content = fmt.Sprintf("%s, %s",
				req.VestingAccount.Address,
				vestingCoins,
			)
		case *launchtypes.RequestContent_ValidatorRemoval:
			requestType = "Remove Validator"
			content = req.ValidatorRemoval.ValAddress
		case *launchtypes.RequestContent_AccountRemoval:
			requestType = "Remove Account"
			content = req.AccountRemoval.Address
		}

		requestEntries = append(requestEntries, []string{
			id,
			requestType,
			content,
		})
	}
	return entrywriter.MustWrite(out, requestSummaryHeader, requestEntries...)
}
