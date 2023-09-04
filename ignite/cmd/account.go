package ignitecmd

import (
	"os"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/ignite/cli/ignite/pkg/cliui/cliquiz"
	"github.com/ignite/cli/ignite/pkg/cliui/entrywriter"
	"github.com/ignite/cli/ignite/pkg/cosmosaccount"
)

const (
	flagAddressPrefix  = "address-prefix"
	flagPassphrase     = "passphrase"
	flagNonInteractive = "non-interactive"
	flagKeyringBackend = "keyring-backend"
	flagKeyringDir     = "keyring-dir"
)

func NewAccount() *cobra.Command {
	c := &cobra.Command{
		Use:   "account [command]",
		Short: "Create, delete, and show Ignite accounts",
		Long: `Commands for managing Ignite accounts. An Ignite account is a private/public
keypair stored in a keyring. Currently Ignite accounts are used when interacting
with Ignite relayer commands and when using "ignite network" commands.

Note: Ignite account commands are not for managing your chain's keys and accounts. Use
you chain's binary to manage accounts from "config.yml". For example, if your
blockchain is called "mychain", use "mychaind keys" to manage keys for the
chain.
`,
		Aliases: []string{"a"},
		Args:    cobra.ExactArgs(1),
	}

	c.PersistentFlags().AddFlagSet(flagSetKeyringBackend())
	c.PersistentFlags().AddFlagSet(flagSetKeyringDir())

	c.AddCommand(
		NewAccountCreate(),
		NewAccountDelete(),
		NewAccountShow(),
		NewAccountList(),
		NewAccountImport(),
		NewAccountExport(),
	)

	return c
}

func printAccounts(cmd *cobra.Command, accounts ...cosmosaccount.Account) error {
	var accEntries [][]string
	for _, acc := range accounts {
		addr, err := acc.Address(getAddressPrefix(cmd))
		if err != nil {
			return err
		}

		pubKey, err := acc.PubKey()
		if err != nil {
			return err
		}

		accEntries = append(accEntries, []string{acc.Name, addr, pubKey})
	}
	return entrywriter.MustWrite(os.Stdout, []string{"name", "address", "public key"}, accEntries...)
}

func flagSetKeyringBackend() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.String(flagKeyringBackend, string(cosmosaccount.KeyringTest), "keyring backend to store your account keys")
	return fs
}

func getKeyringBackend(cmd *cobra.Command) cosmosaccount.KeyringBackend {
	backend, _ := cmd.Flags().GetString(flagKeyringBackend)
	return cosmosaccount.KeyringBackend(backend)
}

func flagSetKeyringDir() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.String(flagKeyringDir, cosmosaccount.KeyringHome, "accounts keyring directory")
	return fs
}

func getKeyringDir(cmd *cobra.Command) string {
	keyringDir, _ := cmd.Flags().GetString(flagKeyringDir)
	return keyringDir
}

func flagSetAccountPrefixes() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.String(flagAddressPrefix, cosmosaccount.AccountPrefixCosmos, "account address prefix")
	return fs
}

func getAddressPrefix(cmd *cobra.Command) string {
	prefix, _ := cmd.Flags().GetString(flagAddressPrefix)
	return prefix
}

func flagSetAccountImport() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.Bool(flagNonInteractive, false, "do not enter into interactive mode")
	fs.String(flagPassphrase, "", "passphrase to decrypt the imported key (ignored when secret is a mnemonic)")
	return fs
}

func flagSetAccountExport() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.Bool(flagNonInteractive, false, "do not enter into interactive mode")
	fs.String(flagPassphrase, "", "passphrase to encrypt the exported key")
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
