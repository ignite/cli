package starportcmd

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
)

func NewNetworkAccountImport() *cobra.Command {
	c := &cobra.Command{
		Use:   "import [name] [password] [path-to-private-key]",
		Short: "Import an account",
		RunE:  networkAccountImportHandler,
		Args:  cobra.ExactArgs(3),
	}
	return c
}

func networkAccountImportHandler(cmd *cobra.Command, args []string) error {
	var name, password, privateKeyPath = args[0], args[1], args[2]
	privateKey, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return err
	}
	b, err := newNetworkBuilder()
	if err != nil {
		return err
	}
	if err := b.AccountImport(name, string(privateKey), password); err != nil {
		return err
	}
	fmt.Println("🗿 Account imported")
	return nil
}
