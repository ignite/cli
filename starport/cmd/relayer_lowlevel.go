package starportcmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/nodetime"
)

func NewRelayerLowLevel() *cobra.Command {
	c := &cobra.Command{
		Use:   "lowlevel",
		Short: "Low-level relayer commands from @confio/relayer",
	}
	c.AddCommand(NewRelayerLowLevelIBCSetup())
	c.AddCommand(NewRelayerLowLevelIBCRelayer())
	return c
}

func NewRelayerLowLevelIBCSetup() *cobra.Command {
	c := &cobra.Command{
		Use:   "ibc-setup [--] [...]",
		Short: "Collection of commands to quickly setup a relayer",
		RunE:  relayerLowLevelHandle(nodetime.CommandIBCSetup),
		Example: `starport relayer lowlevel ibc-setup -- -h
starport relayer lowlevel ibc-setup -- init --src relayer_test_1 --dest relayer_test_2`,
	}
	return c
}

func NewRelayerLowLevelIBCRelayer() *cobra.Command {
	c := &cobra.Command{
		Use:     "ibc-relayer [--] [...]",
		Short:   "Typescript implementation of an IBC relayer",
		RunE:    relayerLowLevelHandle(nodetime.CommandIBCRelayer),
		Example: `starport relayer lowlevel ibc-relayer -- -h`,
	}
	return c
}

func relayerLowLevelHandle(c nodetime.CommandName) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		command, cleanup, err := nodetime.Command(c)
		if err != nil {
			return err
		}
		defer cleanup()

		command = append(command, args...)

		return cmdrunner.New().Run(
			cmd.Context(),
			step.New(
				step.Exec(command[0], command[1:]...),
				step.Stdout(os.Stdout),
				step.Stderr(os.Stderr),
			),
		)
	}
}
