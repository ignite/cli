package ignitecmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/cosmos/go-bip39"
	"github.com/spf13/cobra"

	"github.com/ignite-hq/cli/ignite/pkg/cliui/cliquiz"
	"github.com/ignite-hq/cli/ignite/pkg/cosmosaccount"
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
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetKeyringDir())
	c.Flags().AddFlagSet(flagSetAccountImportExport())

	return c
}

func accountImportHandler(cmd *cobra.Command, args []string) error {
	var (
		name           = args[0]
		secret, _      = cmd.Flags().GetString(flagSecret)
		keyringBackend = getKeyringBackend(cmd)
		keyringDir     = getKeyringDir(cmd)
	)

	if secret == "" {
		if err := cliquiz.Ask(
			cliquiz.NewQuestion("Your mnemonic or path to your private key", &secret, cliquiz.Required())); err != nil {
			return err
		}
	}

	passphrase, err := getPassphrase(cmd)
	if err != nil {
		return err
	}

	if !bip39.IsMnemonicValid(secret) {
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
		cosmosaccount.WithKeyringBackend(keyringBackend),
		cosmosaccount.WithKeyringDir(keyringDir),
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
