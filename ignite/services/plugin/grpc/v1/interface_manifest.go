package v1

import "github.com/spf13/cobra"

// ImportCobraCommand appends Cobra command definitions to the list of plugin commands.
// This method can be used in cases where a plugin defines the commands using Cobra.
func (m *Manifest) ImportCobraCommand(cmd *cobra.Command, placeCommandUnder string) {
	m.Commands = append(m.Commands, convertCobraCommand(cmd, placeCommandUnder))
}

func convertCobraCommand(c *cobra.Command, placeCommandUnder string) *Command {
	cmd := &Command{
		Use:               c.Use,
		Aliases:           c.Aliases,
		Short:             c.Short,
		Long:              c.Long,
		Hidden:            c.Hidden,
		PlaceCommandUnder: placeCommandUnder,
		Flags:             extractCobraFlags(c),
	}

	for _, c := range c.Commands() {
		cmd.Commands = append(cmd.Commands, convertCobraCommand(c, ""))
	}

	return cmd
}
