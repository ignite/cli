package starportcmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/tendermint/starport/starport/pkg/cliquiz"
	"github.com/tendermint/starport/starport/pkg/cosmosaccount"
)

const (
	flagAddressPrefix  = "address-prefix"
	flagPassphrase     = "passphrase"
	flagNonInteractive = "non-interactive"
)

func NewAccount() *cobra.Command {
	c := &cobra.Command{
		Use:     "account [command]",
		Short:   "Manage your cosmos chain accounts",
		Aliases: []string{"a"},
		Args:    cobra.ExactArgs(1),
	}

	c.AddCommand(NewAccountCreate())
	c.AddCommand(NewAccountDelete())
	c.AddCommand(NewAccountShow())
	c.AddCommand(NewAccountList())
	c.AddCommand(NewAccountImport())
	c.AddCommand(NewAccountExport())

	return c
}

func printAccounts(cmd *cobra.Command, accounts ...cosmosaccount.Account) {
	w := &tabwriter.Writer{}
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)

	if len(accounts) == 0 {
		return
	}

	fmt.Fprintln(w, "name\taddress\tpublic key")

	for _, acc := range accounts {
		fmt.Fprintf(w, "%s\t%s\t%s\n",
			acc.Name,
			acc.Address(getAddressPrefix(cmd)),
			acc.PubKey(getAddressPrefix(cmd)),
		)
	}

	fmt.Fprintln(w)
	w.Flush()
}

func flagSetKeyringBackend() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.String(flagKeyringBackend, "test", "Keyring backend to store your account keys")
	return fs
}

func getKeyringBackend(cmd *cobra.Command) string {
	backend, _ := cmd.Flags().GetString(flagKeyringBackend)
	return backend
}

func flagSetAccountPrefixes() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.String(flagAddressPrefix, "cosmos", "Account address prefix")
	return fs
}

func getAddressPrefix(cmd *cobra.Command) string {
	prefix, _ := cmd.Flags().GetString(flagAddressPrefix)
	return prefix
}

func flagSetAccountImportExport() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.Bool(flagNonInteractive, false, "Do not enter into interactive mode")
	fs.String(flagPassphrase, "", "Account passphrase")
	return fs
}

func getIsNonInteractive(cmd *cobra.Command) bool {
	is, _ := cmd.Flags().GetBool(flagNonInteractive)
	return is
}

func getPassphrase(cmd *cobra.Command) (string, error) {
	pass, _ := cmd.Flags().GetString(flagPassphrase)

	if pass == "" && !getIsNonInteractive(cmd) {
		if err := cliquiz.Ask(
			cliquiz.NewQuestion("Passphrase",
				&pass,
				cliquiz.HideAnswer(),
				cliquiz.GetConfirmation(),
			)); err != nil {
			return "", err
		}
	}

	return pass, nil
}
