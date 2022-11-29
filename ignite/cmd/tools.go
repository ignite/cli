package ignitecmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cmdrunner"
	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/ignite/pkg/nodetime"
	"github.com/ignite/cli/ignite/pkg/protoc"
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
	c.AddCommand(NewToolsProtoc())
	return c
}

func NewToolsIBCSetup() *cobra.Command {
	return &cobra.Command{
		Use:   "ibc-setup [--] [...]",
		Short: "Collection of commands to quickly setup a relayer",
		RunE:  toolsNodetimeProxy(nodetime.CommandIBCSetup),
		Example: `ignite tools ibc-setup -- -h
ignite tools ibc-setup -- init --src relayer_test_1 --dest relayer_test_2`,
	}
}

func NewToolsIBCRelayer() *cobra.Command {
	return &cobra.Command{
		Use:     "ibc-relayer [--] [...]",
		Short:   "TypeScript implementation of an IBC relayer",
		RunE:    toolsNodetimeProxy(nodetime.CommandIBCRelayer),
		Example: `ignite tools ibc-relayer -- -h`,
	}
}

func NewToolsProtoc() *cobra.Command {
	return &cobra.Command{
		Use:     "protoc [--] [...]",
		Short:   "Execute the protoc command",
		Long:    "The protoc command. You don't need to setup the global protoc include folder with -I, it's automatically handled",
		RunE:    toolsProtocProxy,
		Example: `ignite tools protoc -- --version`,
	}
}

func toolsNodetimeProxy(c nodetime.CommandName) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		command, cleanup, err := nodetime.Command(c)
		if err != nil {
			return err
		}
		defer cleanup()

		return toolsProxy(cmd.Context(), append(command, args...))
	}
}

func toolsProtocProxy(cmd *cobra.Command, args []string) error {
	command, cleanup, err := protoc.Command()
	if err != nil {
		return err
	}
	defer cleanup()

	return toolsProxy(cmd.Context(), append(command.Command(), args...))
}

func toolsProxy(ctx context.Context, command []string) error {
	return cmdrunner.New().Run(
		ctx,
		step.New(
			step.Exec(command[0], command[1:]...),
			step.Stdout(os.Stdout),
			step.Stderr(os.Stderr),
		),
	)
}
