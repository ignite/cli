package starportcmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/cosmos/go-bip39"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/cliquiz"
	"github.com/tendermint/starport/starport/pkg/cosmosaccount"
)

func NewAccountImport() *cobra.Command {
	c := &cobra.Command{
		Use:   "import [name] [mnemonic or private_key_path]",
		Short: "Import an account",
		Args:  cobra.ExactArgs(2),
		RunE:  accountImportHandler,
	}

	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetAccountImportExport())

	return c
}

func accountImportHandler(cmd *cobra.Command, args []string) error {
	var (
		name       = args[0]
		secret     = strings.TrimSpace(args[1])
		passphrase = getPassphrase(cmd)
	)

	if passphrase == "" && !getIsNonInteractive(cmd) {
		if err := cliquiz.Ask(cliquiz.NewQuestion("Passphrase", &passphrase, cliquiz.HideAnswer())); err != nil {
			return err
		}
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

	ca, err := cosmosaccount.New(getKeyringBackend(cmd))
	if err != nil {
		return err
	}

	if _, err := ca.Import(name, secret, passphrase); err != nil {
		return err
	}

	fmt.Printf("Account %q imported.\n", name)
	return nil
}
