package ignitecmd

import (
	"os"

	"github.com/cosmos/go-bip39"
	"github.com/spf13/cobra"

	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/cliui/bubbleconfirm"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
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
	c.Flags().AddFlagSet(flagSetCoinType())

	return c
}

func accountImportHandler(cmd *cobra.Command, args []string) error {
	var (
		name      = args[0]
		secret, _ = cmd.Flags().GetString(flagSecret)
		session   = cliui.New(cliui.StartSpinnerWithText(statusImporting))
	)
	defer session.End()

	if secret == "" {
		session.StopSpinner()

		if err := bubbleconfirm.Ask(
			bubbleconfirm.NewQuestion("Your mnemonic or path to your private key", &secret, bubbleconfirm.Required())); err != nil {
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
		cosmosaccount.WithCoinType(getCoinType(cmd)),
	)
	if err != nil {
		return err
	}

	if _, err := ca.Import(name, secret, passphrase); err != nil {
		return err
	}

	return session.Printf("Account %q imported.\n", name)
}
