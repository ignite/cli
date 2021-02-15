package starportcmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/services/scaffolder"
)

// NewIBCPacket creates a new packet in the module
func NewIBCPacket() *cobra.Command {
	c := &cobra.Command{
		Use:   "packet [packetName] [field1] [field2] ... --module [module_name]",
		Short: "Creates a new interpretable IBC packet.",
		Long:  "Use starport ibc packet to create a new packet in your IBC module.",
		Args:  cobra.MinimumNArgs(1),
		RunE:  createPacketHandler,
	}

	c.Flags().String(moduleFlag, "", "IBC Module to add the packet into")

	return c
}

func createPacketHandler(cmd *cobra.Command, args []string) error {
	// Get the module
	module, err := cmd.Flags().GetString(moduleFlag)
	if err != nil {
		return err
	}
	if module == "" {
		return errors.New("please specify a module to create the packet into: --module <module_name>")
	}

	sc := scaffolder.New(appPath)
	if err := sc.AddPacket(module, args[0], args[1:]...); err != nil {
		return err
	}
	fmt.Printf("\n🎉 Created a packet `%[1]v`.\n\n", args[0])
	return nil
}
