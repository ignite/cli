package starportcmd

import (
	"github.com/spf13/cobra"
)

// NewNetworkChainList returns a new command to list all published chains on Starport Network
func NewNetworkChainList() *cobra.Command {
	c := &cobra.Command{
		Use:   "list",
		Short: "List published chains",
		Args:  cobra.NoArgs,
		RunE:  networkChainListHandler,
	}

	return c
}

func networkChainListHandler(cmd *cobra.Command, args []string) error {
	return nil
}
