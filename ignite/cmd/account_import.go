package ignitecmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/cosmos/go-bip39"
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui/cliquiz"
	"github.com/ignite/cli/ignite/pkg/cosmosaccount"
)

const flagSecret = "secret"

func NewAccountImport() *cobra.Command {
	c := &cobra.Command{
		Use:   "import [name]",
		Short: "Import an account by using a mnemonic or a private key",
		Args:  cobra.ExactArgs(1),
		RunE:  accountImportHandler,
	}

	c.Flags().String(flagSecret, "", "Your mnemonic or path to your private key (use interactive mode instead to securely pass your mnemonic)")
	c.Flags().AddFlagSet(flagSetAccountImport())

	return c
}

func accountImportHandler(cmd *cobra.Command, args []string) error {
	var (
		name      = args[0]
		secret, _ = cmd.Flags().GetString(flagSecret)
	)

	if secret == "" {
		if err := cliquiz.Ask(
			cliquiz.NewQuestion("Your mnemonic or path to your private key", &secret, cliquiz.Required())); err != nil {
			return err
		}
	}

	var passphrase string
	if !bip39.IsMnemonicValid(secret) {
		var err error
		passphrase, err = getPassphrase(cmd)
		if err != nil {
			return err
		}

		privKey, err := os.ReadFile(secret)
		if os.IsNotExist(err) {
			return errors.New("mnemonic is not valid or private key not found at path")
		}
		if err != nil {
			return err
		}
		secret = string(privKey)
	}

	ca, err := cosmosaccount.New(
		cosmosaccount.WithKeyringBackend(getKeyringBackend(cmd)),
		cosmosaccount.WithHome(getKeyringDir(cmd)),
	)
	if err != nil {
		return err
	}

	if _, err := ca.Import(name, secret, passphrase); err != nil {
		return err
	}

	fmt.Printf("Account %q imported.\n", name)
	return nil
}
