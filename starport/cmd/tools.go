package starportcmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/nodetime"
)

// NewTools returns a command where various tools (binaries) are attached as sub commands
// for advanced users.
func NewTools() *cobra.Command {
	c := &cobra.Command{
		Use:   "tools",
		Short: "Tools for advanced users",
	}
	c.AddCommand(NewToolsIBCSetup())
	c.AddCommand(NewToolsIBCRelayer())
	return c
}

func NewToolsIBCSetup() *cobra.Command {
	c := &cobra.Command{
		Use:   "ibc-setup [--] [...]",
		Short: "Collection of commands to quickly setup a relayer",
		RunE:  toolsNodetimeProxy(nodetime.CommandIBCSetup),
		Example: `starport tools ibc-setup -- -h
starport relayer lowlevel ibc-setup -- init --src relayer_test_1 --dest relayer_test_2`,
	}
	return c
}

func NewToolsIBCRelayer() *cobra.Command {
	c := &cobra.Command{
		Use:     "ibc-relayer [--] [...]",
		Short:   "Typescript implementation of an IBC relayer",
		RunE:    toolsNodetimeProxy(nodetime.CommandIBCRelayer),
		Example: `starport tools ibc-relayer -- -h`,
	}
	return c
}

func toolsNodetimeProxy(c nodetime.CommandName) func(cmd *cobra.Command, args []string) error {
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
