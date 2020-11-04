package starportcmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/services/networkbuilder"
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
	b, err := networkbuilder.New(spnAddress)
	if err != nil {
		return err
	}
	account, err := b.AccountInUse()
	if err != nil {
		return err
	}
	fmt.Printf("🗿 Your spn account is: %s\n", color.New(color.FgYellow).SprintFunc()(account.Name))
	return nil
}
