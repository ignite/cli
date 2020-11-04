package starportcmd

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/services/networkbuilder"
)

func NewNetworkAccountExport() *cobra.Command {
	c := &cobra.Command{
		Use:   "export [name] [password] [save-path]",
		Short: "Export an account",
		RunE:  networkAccountExportHandler,
		Args:  cobra.ExactArgs(3),
	}
	return c
}

func networkAccountExportHandler(cmd *cobra.Command, args []string) error {
	var name, password, privateKeyPath = args[0], args[1], args[2]
	b, err := networkbuilder.New(spnAddress)
	if err != nil {
		return err
	}
	privateKey, err := b.AccountExport(name, password)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(privateKeyPath, []byte(privateKey), 0755); err != nil {
		return err
	}
	privateKeyPathAbs, err := filepath.Abs(privateKeyPath)
	if err != nil {
		return err
	}
	fmt.Printf(`📩 Account exported.

Your private key saved to: %s
Please do not forget your password, it'll be later used to decrypt your private key while importing.

`, color.New(color.FgYellow).SprintFunc()(privateKeyPathAbs))
	return nil
}
