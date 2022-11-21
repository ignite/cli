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
		Long: `The "approve" command is used by a chain's coordinator to approve requests.
Multiple requests can be approved using a comma-separated list and/or using a
dash syntax.

	ignite network request approve 42 1,2,3-6,7,8

The command above approves requests with IDs from 1 to 8 included on a chain
with a launch ID 42.

When requests are approved Ignite applies the requested changes and simulates
initializing and launching the chain locally. If the chain starts successfully,
requests are considered to be "verified" and are approved. If one or more
requested changes stop the chain from launching locally, the verification
process fails and the approval of all requests is canceled. To skip the
verification process use the "--no-verification" flag.

Note that Ignite will try to approve requests in the same order as request IDs
are submitted to the "approve" command.`,
		RunE: networkRequestApproveHandler,
		Args: cobra.ExactArgs(2),
	}

	flagSetClearCache(c)
	c.Flags().Bool(flagNoVerification, false, "approve the requests without verifying them")
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetKeyringDir())
	return c
}

func networkRequestApproveHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New(cliui.StartSpinner())
	defer session.End()

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
	if err := n.SubmitRequest(cmd.Context(), launchID, reviewals...); err != nil {
		return err
	}

	return session.Printf("%s Request(s) %s approved\n", icons.OK, numbers.List(ids, "#"))
}
