package starportcmd

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/numbers"
	"github.com/tendermint/starport/starport/services/network"
)

const (
	flagNoVerification = "no-verification"
)

// NewNetworkRequestApprove creates a new request approve
// command to approve requests for a chain.
func NewNetworkRequestApprove() *cobra.Command {
	c := &cobra.Command{
		Use:     "approve [launch-id] [number<,...>]",
		Aliases: []string{"accept"},
		Short:   "Approve requests",
		RunE:    networkRequestApproveHandler,
		Args:    cobra.ExactArgs(2),
	}
	c.Flags().Bool(flagNoVerification, false, "approve the requests without verifying them")
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	return c
}

func networkRequestApproveHandler(cmd *cobra.Command, args []string) error {
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

	// Get the list of request ids
	ids, err := numbers.ParseListRange(args[1])
	if err != nil {
		return err
	}

	// Verify the requests are valid
	noVerification, err := cmd.Flags().GetBool(flagNoVerification)
	if err != nil {
		return err
	}

	if !noVerification {
		// Verify the request. This operation generate the genesis
		// in a temporary directory and verify this genesis is valid
		err := nb.VerifyRequests(cmd.Context(), launchID, ids)
		if err != nil {
			return err
		}
	}
	// Submit the approved requests
	reviewals := make([]network.Reviewal, 0)
	for _, id := range ids {
		reviewals = append(reviewals, network.ApproveProposal(id))
	}
	if err := nb.SubmitRequest(launchID, reviewals...); err != nil {
		return err
	}

	s.Stop()
	fmt.Printf("Request(s) %s approved ✅\n", numbers.List(ids, "#"))
	return nil
}
