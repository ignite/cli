package starportcmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	launchtypes "github.com/tendermint/spn/x/launch/types"
)

// NewNetworkRequestList creates a new request list command to list
// requests for a chain by filtering options.
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
	nb, s, shutdown, err := initializeNetwork(cmd)
	if err != nil {
		return err
	}
	defer shutdown()

	// parse launch ID
	launchID, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return errors.Wrap(err, "error parsing launchID")
	}
	if launchID == 0 {
		return errors.New("launch ID must be greater than 0")
	}

	requests, err := nb.FetchRequests(cmd.Context(), launchID)
	if err != nil {
		return err
	}

	// Rendering
	requestTable := tablewriter.NewWriter(os.Stdout)
	requestTable.SetHeader([]string{"ID", "Type", "Content"})

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
			content = fmt.Sprintf("%s, %s, %s",
				req.GenesisValidator.Peer,
				req.GenesisValidator.Address,
				req.GenesisValidator.SelfDelegation.String())
		case *launchtypes.RequestContent_VestingAccount:
			requestType = "Add Vesting Account"
			content = fmt.Sprintf("%s, %s",
				req.VestingAccount.Address,
				req.VestingAccount.StartingBalance.String())
		case *launchtypes.RequestContent_ValidatorRemoval:
			requestType = "Remove Validator"
			content = req.ValidatorRemoval.ValAddress
		case *launchtypes.RequestContent_AccountRemoval:
			requestType = "Remove Account"
			content = req.AccountRemoval.Address
		}

		requestTable.Append([]string{
			id,
			requestType,
			content,
		})
	}

	s.Stop()
	requestTable.Render()
	return nil
}
