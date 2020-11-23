package starportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewNetworkAccount() *cobra.Command {
	c := &cobra.Command{
		Use:   "account",
		Short: "Show the underlying SPN account",
		Long: `Show the underlying SPN account.
To pick another account see "starport network account use -h"
If no account is picked, default "spn" account is used.
`,
		RunE: networkAccountGetHandler,
	}
	c.AddCommand(NewNetworkAccountCreate())
	c.AddCommand(NewNetworkAccountImport())
	c.AddCommand(NewNetworkAccountExport())
	c.AddCommand(NewNetworkAccountUse())
	return c
}

func networkAccountGetHandler(cmd *cobra.Command, args []string) error {
	account, err := nb.AccountInUse()
	if err == nil {
		fmt.Printf("🗿 Your spn account is: %s: %s\n", infoColor(account.Name), infoColor(account.Address))
	}
	return nil
}
