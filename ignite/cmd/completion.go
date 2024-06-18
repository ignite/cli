package ignitecmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewCompletionCmd represents the completion command.
func NewCompletionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generates shell completion script.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				if err := cmd.Help(); err != nil {
					fmt.Fprintln(os.Stderr, "Error displaying help:", err)
					os.Exit(1)
				}
				os.Exit(0)
			}
			var err error
			switch args[0] {
			case "bash":
				err = cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				err = cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				err = cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				err = cmd.Root().GenPowerShellCompletion(os.Stdout)
			default:
				err = cmd.Help()
			}
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error generating completion script:", err)
				os.Exit(1)
			}
		},
	}
}
