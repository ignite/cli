package starportcmd

import (
	"github.com/spf13/cobra"
)

// NewNetworkChainReward creates a new chain reward command
func NewNetworkChainReward() *cobra.Command {
	c := &cobra.Command{
		Use:   "reward",
		Short: "Manage network rewards",
	}
	c.AddCommand(
		NewNetworkChainRewardSet(),
	)
	return c
}
