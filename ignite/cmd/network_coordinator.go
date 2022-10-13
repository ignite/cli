package ignitecmd

import "github.com/spf13/cobra"

// NewNetworkCoordinator creates a new coordinator command
// it contains sub commands to manage coordinator profile
func NewNetworkCoordinator() *cobra.Command {
	c := &cobra.Command{
		Use:   "coordinator",
		Short: "Interact with coordinator profiles",
	}
	c.AddCommand(
		NewNetworkCoordinatorShow(),
		NewNetworkCoordinatorSet(),
	)
	return c
}
