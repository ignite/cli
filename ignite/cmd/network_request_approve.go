package ignitecmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/ignite/pkg/numbers"
	"github.com/ignite/cli/ignite/services/network"
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

	flagSetClearCache(c)
	c.Flags().Bool(flagNoVerification, false, "approve the requests without verifying them")
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	return c
}

func networkRequestApproveHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New()
	defer session.Cleanup()

	nb, err := newNetworkBuilder(cmd, CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}

	// parse launch ID
	launchID, err := network.ParseID(args[0])
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

	cacheStorage, err := newCache(cmd)
	if err != nil {
		return err
	}

	// if requests must be verified, we simulate the chain in a temporary directory with the requests
	if !noVerification {
		if err := verifyRequest(cmd.Context(), cacheStorage, nb, launchID, ids...); err != nil {
			return errors.Wrap(err, "request(s) not valid")
		}
		session.Printf("%s Request(s) %s verified\n", icons.OK, numbers.List(ids, "#"))
	}

	// Submit the approved requests
	reviewals := make([]network.Reviewal, 0)
	for _, id := range ids {
		reviewals = append(reviewals, network.ApproveRequest(id))
	}
	if err := n.SubmitRequest(launchID, reviewals...); err != nil {
		return err
	}

	session.StopSpinner()

	return session.Printf("%s Request(s) %s approved\n", icons.OK, numbers.List(ids, "#"))
}
