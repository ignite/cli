package starportcmd

import (
	"github.com/spf13/cobra"
)

const (
	flagRecover     = "recover"
	flagMnemonic = "mnemomic"
	flagKeyName = "key-name"
	flagOut = "out"
)

// NewNetworkChainInit returns a new command to initialize a chain from a published chain ID
func NewNetworkChainInit() *cobra.Command {
	c := &cobra.Command{
		Use:   "init [chain-id]",
		Short: "Initialize a chain from a published chain ID",
		Args:  cobra.ExactArgs(1),
		RunE:  networkChainInitHandler,
	}

	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().Bool(flagRecover, false, "Recover chain account from a mnemonic")
	c.Flags().String(flagMnemonic, "", "Mnemonic for recovered account")
	c.Flags().String(flagKeyName, "", "key name for the chain account")

	return c
}

func networkChainInitHandler(cmd *cobra.Command, args []string) error {
	var (
		chainID        = args[0]
		recover, _     = cmd.Flags().GetBool(flagRecover)
		mnemonic, _       = cmd.Flags().GetString(flagMnemonic)
		keyName, _ = cmd.Flags().GetString(flagKeyName)
	)

	return nil
}