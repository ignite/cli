package starportcmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/services/scaffolder"
)

const (
	ackFlag = "ack"
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

	c.Flags().StringSlice(ackFlag, []string{}, "Custom acknowledgment type (field1,field2,...)")
	c.Flags().String(moduleFlag, "", "IBC Module to add the packet into")

	return c
}

func createPacketHandler(cmd *cobra.Command, args []string) error {
	s := clispinner.New().SetText("Scaffolding...")
	defer s.Stop()

	var (
		packet       = args[0]
		packetFields = args[1:]
	)

	module, err := cmd.Flags().GetString(moduleFlag)
	if err != nil {
		return err
	}
	if module == "" {
		return errors.New("please specify a module to create the packet into: --module <module_name>")
	}

	ackFields, err := cmd.Flags().GetStringSlice(ackFlag)
	if err != nil {
		return err
	}

	sc := scaffolder.New(appPath)
	if err := sc.AddPacket(module, packet, packetFields, ackFields); err != nil {
		return err
	}

	s.Stop()

	fmt.Printf("\n🎉 Created a packet `%[1]v`.\n\n", args[0])
	return nil
}
