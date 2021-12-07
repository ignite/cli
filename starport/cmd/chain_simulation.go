package starportcmd

import (
	"github.com/spf13/cobra"
)

// NewChainSimulation creates a new simulation command to run the blockchain simulation.
func NewChainSimulation() *cobra.Command {
	c := &cobra.Command{
		Use:   "simulation",
		Short: "Run the blockchain simulation node in development",
		Long:  "Run the blockchain simulation for all chain modules",
		Args:  cobra.ExactArgs(0),
		RunE:  chainSimulationHandler,
	}

	c.Flags().AddFlagSet(flagSetHome())

	return c
}

func chainSimulationHandler(cmd *cobra.Command, args []string) error {
	// create the chain
	c, err := newChainWithHomeFlags(cmd)
	if err != nil {
		return err
	}
	return c.Simulate(cmd.Context())
}
