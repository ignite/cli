package starportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewNetworkAccountCreate creates a new account create command to create
// an account to be used in SPN.
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
	nb, err := newNetworkBuilder()
	if err != nil {
		return err
	}
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
