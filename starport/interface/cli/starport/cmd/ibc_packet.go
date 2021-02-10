package starportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/services/scaffolder"
)

// NewIBCPacket creates a new packet in the module
func NewIBCPacket() *cobra.Command {
	c := &cobra.Command{
		Use:   "packet [moduleName] [packetName] [field1] [field2] ...",
		Short: "Creates a new interpretable IBC packet.",
		Long:  "Use starport ibc packet to create a new packet in your IBC module.",
		Args:  cobra.MinimumNArgs(2),
		RunE:  createPacketHandler,
	}
	return c
}

func createPacketHandler(cmd *cobra.Command, args []string) error {
	sc := scaffolder.New(appPath)
	if err := sc.AddPacket(args[0], args[1], args[2:]...); err != nil {
		return err
	}
	fmt.Printf("\n🎉 Created a packet `%[1]v`.\n\n", args[1])
	return nil
}
