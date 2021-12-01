package starportcmd

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/numbers"
	"github.com/tendermint/starport/starport/services/network"
)

// NewNetworkRequestReject creates a new request reject
// command to reject requests for a chain.
func NewNetworkRequestReject() *cobra.Command {
	c := &cobra.Command{
		Use:     "reject [launch-id] [number<,...>]",
		Aliases: []string{"accept"},
		Short:   "Reject requests",
		RunE:    networkRequestRejectHandler,
		Args:    cobra.ExactArgs(2),
	}
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	return c
}

func networkRequestRejectHandler(cmd *cobra.Command, args []string) error {
	// initialize network common methods
	nb, err := newNetworkBuilder(cmd)
	if err != nil {
		return err
	}
	defer nb.Cleanup()

	// parse launch ID
	launchID, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return errors.Wrap(err, "error parsing launchID")
	}
	if launchID == 0 {
		return errors.New("launch ID must be greater than 0")
	}

	// Get the list of request ids
	ids, err := numbers.ParseList(args[1])
	if err != nil {
		return err
	}

	n, err := nb.Network()
	if err != nil {
		return err
	}

	// Submit the rejected requests
	reviewals := make([]network.Reviewal, 0)
	for _, id := range ids {
		reviewals = append(reviewals, network.RejectRequest(id))
	}
	if err := n.SubmitRequest(launchID, reviewals...); err != nil {
		return err
	}

	nb.Spinner.Stop()
	fmt.Printf("%s Request(s) %s rejected\n", clispinner.OK, numbers.List(ids, "#"))
	return nil
}
