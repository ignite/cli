package ignitecmd

import "github.com/spf13/cobra"

// NewNetworkCoordinator creates a new coordinator command
// it contains sub commands to manage coordinator profile
func NewNetworkCoordinator() *cobra.Command {
	c := &cobra.Command{
		Use:   "coordinator",
		Short: "Show and update a coordinator profile",
	}
	c.AddCommand(
		NewNetworkCoordinatorShow(),
		NewNetworkCoordinatorSet(),
	)
	return c
}
