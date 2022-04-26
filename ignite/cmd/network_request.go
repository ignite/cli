package ignitecmd

import "github.com/spf13/cobra"

// NewNetworkRequest creates a new approval request command that holds some other
// sub commands related to handle request for a chain.
func NewNetworkRequest() *cobra.Command {
	c := &cobra.Command{
		Use:   "request",
		Short: "Handle requests",
	}

	c.AddCommand(
		NewNetworkRequestShow(),
		NewNetworkRequestList(),
		NewNetworkRequestApprove(),
		NewNetworkRequestReject(),
		NewNetworkRequestVerify(),
	)

	return c
}
