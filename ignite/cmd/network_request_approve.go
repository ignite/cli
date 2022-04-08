package starportcmd

import (
	"fmt"

	"github.com/ignite-hq/cli/ignite/pkg/clispinner"
	"github.com/ignite-hq/cli/ignite/pkg/numbers"
	"github.com/ignite-hq/cli/ignite/services/network"
	"github.com/spf13/cobra"
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
	nb, err := newNetworkBuilder(cmd)
	if err != nil {
		return err
	}
	defer nb.Cleanup()

	// parse launch ID
	launchID, err := network.ParseLaunchID(args[0])
	if err != nil {
		return err
	}

	// Get the list of request ids
	ids, err := numbers.ParseList(args[1])
	if err != nil {
		return err
	}

	// Verify the requests are valid
	noVerification, err := cmd.Flags().GetBool(flagNoVerification)
	if err != nil {
		return err
	}

	n, err := nb.Network()
	if err != nil {
		return err
	}

	if !noVerification {
		err := n.VerifyRequests(cmd.Context(), launchID, ids...)
		if err != nil {
			return err
		}
	}
	// Submit the approved requests
	reviewals := make([]network.Reviewal, 0)
	for _, id := range ids {
		reviewals = append(reviewals, network.ApproveRequest(id))
	}
	if err := n.SubmitRequest(launchID, reviewals...); err != nil {
		return err
	}

	nb.Spinner.Stop()
	fmt.Printf("%s Request(s) %s approved\n", clispinner.OK, numbers.List(ids, "#"))
	return nil
}
