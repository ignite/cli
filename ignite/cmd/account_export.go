package ignitecmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

func NewAccountExport() *cobra.Command {
	c := &cobra.Command{
		Use:   "export [name]",
		Short: "Export an account as a private key",
		Args:  cobra.ExactArgs(1),
		RunE:  accountExportHandler,
	}

	c.Flags().AddFlagSet(flagSetAccountExport())
	c.Flags().String(flagPath, "", "path to export private key. default: ./key_[name]")

	return c
}

func accountExportHandler(cmd *cobra.Command, args []string) error {
	var (
		name    = args[0]
		path    = flagGetPath(cmd)
		session = cliui.New(cliui.StartSpinnerWithText(statusExporting))
	)
	defer session.End()

	passphrase, err := getPassphrase(cmd)
	if err != nil {
		return err
	}
	const minPassLength = 8
	if len(passphrase) < minPassLength {
		return errors.Errorf("passphrase must be at least %d characters", minPassLength)
	}

	ca, err := cosmosaccount.New(
		cosmosaccount.WithKeyringBackend(getKeyringBackend(cmd)),
		cosmosaccount.WithHome(getKeyringDir(cmd)),
	)
	if err != nil {
		return err
	}

	armored, err := ca.Export(name, passphrase)
	if err != nil {
		return err
	}

	if path == "" {
		path = fmt.Sprintf("./key_%s", name)
	}
	path, err = filepath.Abs(path)
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, []byte(armored), 0o600); err != nil {
		return err
	}

	return session.Printf("Account %q exported to file: %s\n", name, path)
}
