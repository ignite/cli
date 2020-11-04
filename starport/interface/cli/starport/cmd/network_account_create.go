package starportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/services/networkbuilder"
)

func NewNetworkAccountCreate() *cobra.Command {
	c := &cobra.Command{
		Use:   "create [name]",
		Short: "Create an account",
		RunE:  networkAccountCreateHandler,
		Args:  cobra.ExactArgs(1),
	}
	return c
}

func networkAccountCreateHandler(cmd *cobra.Command, args []string) error {
	b, err := networkbuilder.New(spnAddress)
	if err != nil {
		return err
	}
	account, err := b.AccountCreate(args[0])
	if err != nil {
		return err
	}
	fmt.Printf("🗿 Account created.\nPlease save your mnmenonic in a secret place.\n\n%s\n\n", account.Mnemonic)
	return nil
}
