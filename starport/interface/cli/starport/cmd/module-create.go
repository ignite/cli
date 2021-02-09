package starportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/services/scaffolder"
)

const (
	flagIBC = "ibc"
	flagIBCOrdering = "ibc-ordering"
)

// NewModuleCreate creates a new module create command to scaffold an
// sdk module.
func NewModuleCreate() *cobra.Command {
	c := &cobra.Command{
		Use:   "create [name]",
		Short: "Creates a new empty module to app.",
		Long:  "Use starport module create to create a new empty module to your blockchain.",
		Args:  cobra.MinimumNArgs(1),
		RunE:  createModuleHandler,
	}
	c.Flags().Bool(flagIBC, false, "scaffold an IBC module")
	c.Flags().String(flagIBCOrdering, "none", "channel ordering of the IBC module [none|ordered|unordered]")
	return c
}

func createModuleHandler(cmd *cobra.Command, args []string) error {
	var options []scaffolder.ModuleCreationOption

	// Check if the module must be an IBC module
	ibcModule, err := cmd.Flags().GetBool(flagIBC)
	if err != nil {
		return err
	}

	if ibcModule {
		options = append(options, scaffolder.WithIBC())

		// Get channel ordering
		ibcOrdering, err := cmd.Flags().GetString(flagIBCOrdering)
		if err != nil {
			return err
		}
		options = append(options, scaffolder.WithIBCChannelOrdering(ibcOrdering))
	}

	name := args[0]
	sc := scaffolder.New(appPath)
	if err := sc.CreateModule(name, options...); err != nil {
		return err
	}
	fmt.Printf("\n🎉 Module created %s.\n\n", name)
	return nil
}
