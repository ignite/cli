package starportcmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/cliquiz"
)

const (
	pathPlaceholder    = "[account].key"
	accountPlaceholder = "[account in use]"
)

func NewNetworkAccountExport() *cobra.Command {
	c := &cobra.Command{
		Use:   "export",
		Short: "Export an account",
		RunE:  networkAccountExportHandler,
	}
	c.Flags().StringP("account", "a", accountPlaceholder, "path to save private key")
	c.Flags().StringP("path", "p", pathPlaceholder, "path to save private key")
	return c
}

func networkAccountExportHandler(cmd *cobra.Command, args []string) error {
	b, err := newNetworkBuilder()
	if err != nil {
		return err
	}

	// prep path and name.
	path, _ := cmd.Flags().GetString("path")
	name, _ := cmd.Flags().GetString("account")
	if name == accountPlaceholder {
		account, err := b.AccountInUse()
		if err != nil {
			return err
		}
		name = account.Name
	}
	if path == pathPlaceholder {
		path = name + ".key"
	}

	// check if chosen account exists before asking a password.
	if _, err := b.AccountGet(name); err != nil {
		return errors.New("account does not exist")
	}

	// ask for encryption password.
	var password string
	if err := cliquiz.Ask(cliquiz.NewQuestion("Encytrpion password", &password, cliquiz.HideAnswer())); err != nil {
		return err
	}
	if password == "" {
		return errors.New("password cannot be empty")
	}

	// generate the private key and save.
	privateKey, err := b.AccountExport(name, password)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(path, []byte(privateKey), 0755); err != nil {
		return err
	}
	privateKeyPathAbs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	fmt.Printf(`
📩 Account exported.

Account name: %s
Your private key saved to: %s
Please do not forget your password, it'll be later used to decrypt your private key while importing.

`,
		infoColor(name),
		infoColor(privateKeyPathAbs),
	)
	return nil
}
