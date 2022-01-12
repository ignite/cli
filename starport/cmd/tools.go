package starportcmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/nodetime"
	"github.com/tendermint/starport/starport/pkg/protoc"
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
	c.AddCommand(NewToolsCompletions())
	return c
}

func NewToolsIBCSetup() *cobra.Command {
	return &cobra.Command{
		Use:   "ibc-setup [--] [...]",
		Short: "Collection of commands to quickly setup a relayer",
		RunE:  toolsNodetimeProxy(nodetime.CommandIBCSetup),
		Example: `starport tools ibc-setup -- -h
starport tools ibc-setup -- init --src relayer_test_1 --dest relayer_test_2`,
	}
}

func NewToolsIBCRelayer() *cobra.Command {
	return &cobra.Command{
		Use:     "ibc-relayer [--] [...]",
		Short:   "Typescript implementation of an IBC relayer",
		RunE:    toolsNodetimeProxy(nodetime.CommandIBCRelayer),
		Example: `starport tools ibc-relayer -- -h`,
	}
}

func NewToolsProtoc() *cobra.Command {
	return &cobra.Command{
		Use:     "protoc [--] [...]",
		Short:   "Execute the protoc command",
		Long:    "The protoc command. You don't need to setup the global protoc include folder with -I, it's automatically handled",
		RunE:    toolsProtocProxy,
		Example: `starport tools protoc -- --version`,
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

	return toolsProxy(cmd.Context(), append(command.Command, args...))
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

func NewToolsCompletions() *cobra.Command {

	// completionCmd represents the completion command
	c := &cobra.Command{
		Use:   "completions",
		Short: "Generate completions script",
		Long: ` The completions command outputs a completion script you can use in your shell. The output script requires 
				that [bash-completion](https://github.com/scop/bash-completion)	is installed and enabled in your 
				system. Since most Unix-like operating systems come with bash-completion by default, bash-completion 
				is probably already installed and operational.

Bash:

  $ source <(starport  tools completions bash)

  To load completions for every new session, run:

  ** Linux **
  $ starport  tools completions bash > /etc/bash_completion.d/starport

  ** macOS **
  $ starport  tools completions bash > /usr/local/etc/bash_completion.d/starport

Zsh:

  If shell completions is not already enabled in your environment, you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  To load completions for each session, execute once:
  
  $ starport  tools completions zsh > "${fpath[1]}/_starport"

  You will need to start a new shell for this setup to take effect.

fish:

  $ starport  tools completions fish | source

  To load completions for each session, execute once:
  
  $ starport  tools completions fish > ~/.config/fish/completionss/starport.fish

PowerShell:

  PS> starport  tools completions powershell | Out-String | Invoke-Expression

  To load completions for every new session, run:
  
  PS> starport  tools completions powershell > starport.ps1
  
  and source this file from your PowerShell profile.
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			}
		},
	}
	return c
}
