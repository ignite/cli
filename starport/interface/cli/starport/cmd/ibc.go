package starportcmd

import "github.com/spf13/cobra"

// NewIBC creates a new ibc command that holds some other sub commands
// related to scaffolding IBC related features
func NewIBC() *cobra.Command {
	c := &cobra.Command{
		Use:   "ibc",
		Short: "Scaffolding for IBC functionalities",
		Args:  cobra.ExactArgs(1),
	}

	// add sub commands.
	c.AddCommand(NewIBCPacket())
	return c
}
