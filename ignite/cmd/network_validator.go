package ignitecmd

import "github.com/spf13/cobra"

// NewNetworkValidator creates a new validator command
// it contains sub commands to manage validator profile
func NewNetworkValidator() *cobra.Command {
	c := &cobra.Command{
		Use:   "validator",
		Short: "Interact with validator profiles",
	}
	c.AddCommand(
		NewNetworkValidatorShow(),
		NewNetworkValidatorSet(),
	)
	return c
}
