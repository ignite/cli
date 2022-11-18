package ignitecmd

import (
	"github.com/spf13/cobra"
)

// NewNetworkChain creates a new chain command that holds some other
// sub commands related to launching a network for a chain.
func NewNetworkChain() *cobra.Command {
	c := &cobra.Command{
		Use:   "chain",
		Short: "Commands to launch chains",
		Long: `The "chain" namespace features the most commonly used commands for launching
blockchains with Ignite.

As a coordinator you "publish" your blockchain to Ignite. When enough validators
are approved for the genesis and no changes are excepted to be made to the
genesis, a coordinator announces that the chain is ready for launch with the
"launch" command. In the case of an unsuccessful launch, the coordinator can revert it
using the "revert-launch" command.

As a validator, you "init" your node and apply to become a validator for a
blockchain with the "join" command. After the launch of the chain is announced,
validators can generate the finalized genesis and download the list of peers with the
"prepare" command.

The "install" command can be used to download, compile the source code and
install the chain's binary locally. The binary can be used, for example, to
initialize a validator node or to interact with the chain after it has been
launched.

All chains published to Ignite can be listed by using the "list" command.
`,
	}

	c.AddCommand(
		NewNetworkChainList(),
		NewNetworkChainPublish(),
		NewNetworkChainInit(),
		NewNetworkChainInstall(),
		NewNetworkChainJoin(),
		NewNetworkChainPrepare(),
		NewNetworkChainShow(),
		NewNetworkChainLaunch(),
		NewNetworkChainRevertLaunch(),
	)

	return c
}
