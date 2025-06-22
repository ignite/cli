package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
	"github.com/ignite/cli/v29/ignite/services/scaffolder"
)

const (
	flagAck = "ack"
)

// NewScaffoldPacket creates a new packet in the module.
func NewScaffoldPacket() *cobra.Command {
	c := &cobra.Command{
		Use:     "packet [packetName] [field1] [field2] ... --module [moduleName]",
		Short:   "Message for sending an IBC packet",
		Long:    "Scaffold an IBC packet in a specific IBC-enabled Cosmos SDK module",
		Args:    cobra.MinimumNArgs(1),
		PreRunE: migrationPreRunHandler,
		RunE:    createPacketHandler,
	}

	flagSetPath(c)
	flagSetClearCache(c)

	c.Flags().AddFlagSet(flagSetYes())
	c.Flags().StringSlice(flagAck, []string{}, "custom acknowledgment type (field1,field2,...)")
	c.Flags().String(flagModule, "", "IBC Module to add the packet into")
	c.Flags().String(flagSigner, "", "label for the message signer (default: creator)")
	c.Flags().Bool(flagNoMessage, false, "disable send message scaffolding")

	return c
}

func createPacketHandler(cmd *cobra.Command, args []string) error {
	var (
		packet       = args[0]
		packetFields = args[1:]
		signer       = flagGetSigner(cmd)
		appPath      = flagGetPath(cmd)
	)

	session := cliui.New(
		cliui.StartSpinnerWithText(statusScaffolding),
		cliui.WithoutUserInteraction(getYes(cmd)),
	)
	defer session.End()

	cfg, _, err := getChainConfig(cmd)
	if err != nil {
		return err
	}

	module, _ := cmd.Flags().GetString(flagModule)
	if module == "" {
		return errors.New("please specify a module to create the packet into: --module <module_name>")
	}

	ackFields, _ := cmd.Flags().GetStringSlice(flagAck)
	noMessage, _ := cmd.Flags().GetBool(flagNoMessage)

	cacheStorage, err := newCache(cmd)
	if err != nil {
		return err
	}

	var options []scaffolder.PacketOption
	if noMessage {
		options = append(options, scaffolder.PacketWithoutMessage())
	} else if signer != "" {
		options = append(options, scaffolder.PacketWithSigner(signer))
	}

	sc, err := scaffolder.New(cmd.Context(), appPath, cfg.Build.Proto.Path)
	if err != nil {
		return err
	}

	err = sc.AddPacket(cmd.Context(), module, packet, packetFields, ackFields, options...)
	if err != nil {
		return err
	}

	sm, err := sc.ApplyModifications(xgenny.ApplyPreRun(scaffolder.AskOverwriteFiles(session)))
	if err != nil {
		return err
	}

	if err := sc.PostScaffold(cmd.Context(), cacheStorage, false); err != nil {
		return err
	}

	modificationsStr, err := sm.String()
	if err != nil {
		return err
	}

	session.Println(modificationsStr)
	session.Printf("\n🎉 Created a packet `%[1]v`.\n\n", args[0])

	return nil
}
