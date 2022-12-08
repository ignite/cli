package ignitecmd

import "github.com/spf13/cobra"

// NewNetworkRequest creates a new approval request command that holds some other
// sub commands related to handle request for a chain.
func NewNetworkRequest() *cobra.Command {
	c := &cobra.Command{
		Use:   "request",
		Short: "Create, show, reject and approve requests",
		Long: `The "request" namespace contains commands for creating, showing, approving, and
rejecting requests.

A request is mechanism in Ignite that allows changes to be made to the genesis
file like adding accounts with token balances and validators. Anyone can submit
a request, but only the coordinator of a chain can approve or reject a request.

Each request has a status:

* Pending: waiting for the approval of the coordinator
* Approved: approved by the coordinator, its content has been applied to the
  launch information
* Rejected: rejected by the coordinator or the request creator
`,
	}

	c.AddCommand(
		NewNetworkRequestShow(),
		NewNetworkRequestList(),
		NewNetworkRequestApprove(),
		NewNetworkRequestReject(),
		NewNetworkRequestVerify(),
		NewNetworkRequestAddAccount(),
		NewNetworkRequestRemoveAccount(),
		NewNetworkRequestRemoveValidator(),
		NewNetworkRequestChangeParam(),
	)

	return c
}
