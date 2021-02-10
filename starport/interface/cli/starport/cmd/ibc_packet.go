package starportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/services/scaffolder"
)

// NewModuleCreate creates a new module create command to scaffold an
// sdk module.
func NewIBCPacket() *cobra.Command {
	c := &cobra.Command{
		Use:   "packet [moduleName] [packetName] [field1] [field2] ...",
		Short: "Creates a new empty module to app.",
		Long:  "Use starport module create to create a new empty module to your blockchain.",
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
