package starportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
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
	account, err := nb.AccountCreate(args[0], "")
	if err != nil {
		return err
	}
	fmt.Printf("🗿 Account created. \n\nAddress: %s\n\nPlease save your mnmenonic in a secret place:\n%s\n\n",
		account.Address,
		account.Mnemonic,
	)
	return nil
}
