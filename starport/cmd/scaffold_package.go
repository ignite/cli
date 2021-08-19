package starportcmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/services/scaffolder"
)

const (
	flagAck = "ack"
)

// NewScaffoldPacket creates a new packet in the module
func NewScaffoldPacket() *cobra.Command {
	c := &cobra.Command{
		Use:   "packet [packetName] [field1] [field2] ... --module [moduleName]",
		Short: "Message for sending an IBC packet",
		Long:  "Scaffold an IBC packet in a specific IBC-enabled Cosmos SDK module",
		Args:  cobra.MinimumNArgs(1),
		RunE:  createPacketHandler,
	}

	c.Flags().StringSlice(flagAck, []string{}, "Custom acknowledgment type (field1,field2,...)")
	c.Flags().String(flagModule, "", "IBC Module to add the packet into")
	c.Flags().String(flagSigner, "", "Label for the message signer (default: creator)")
	c.Flags().Bool(flagNoMessage, false, "Disable send message scaffolding")

	return c
}

func createPacketHandler(cmd *cobra.Command, args []string) error {
	s := clispinner.New().SetText("Scaffolding...")
	defer s.Stop()

	var (
		packet       = args[0]
		packetFields = args[1:]
	)

	module, err := cmd.Flags().GetString(flagModule)
	if err != nil {
		return err
	}
	if module == "" {
		return errors.New("please specify a module to create the packet into: --module <module_name>")
	}

	ackFields, err := cmd.Flags().GetStringSlice(flagAck)
	if err != nil {
		return err
	}

	noMessage, err := cmd.Flags().GetBool(flagNoMessage)
	if err != nil {
		return err
	}

	signer, err := cmd.Flags().GetString(flagSigner)
	if err != nil {
		return err
	}

	var options []scaffolder.PacketOption
	if noMessage {
		options = append(options, scaffolder.PacketWithoutMessage())
	}
	if signer != "" {
		options = append(options, scaffolder.PacketWithSigner(signer))
	}

	sc, err := scaffolder.New(appPath)
	if err != nil {
		return err
	}
	sm, err := sc.AddPacket(placeholder.New(), module, packet, packetFields, ackFields, options...)
	if err != nil {
		return err
	}

	s.Stop()

	fmt.Println(sourceModificationToString(sm))
	fmt.Printf("\n🎉 Created a packet `%[1]v`.\n\n", args[0])
	return nil
}
