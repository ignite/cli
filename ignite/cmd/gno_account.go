package ignitecmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/cliui/bubbleconfirm"
	"github.com/ignite/cli/v29/ignite/pkg/cliui/entrywriter"
	"github.com/ignite/cli/v29/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/v29/ignite/services/gno"
)

// NewGnoAccount returns the `ignite account` command backed by the gno keybase.
func NewGnoAccount() *cobra.Command {
	c := &cobra.Command{
		Use:     "account [command]",
		Aliases: []string{"a"},
		Short:   "Create, delete, and show gno.land accounts",
		Long: `Commands for managing gno.land accounts. An account is a private/public keypair
stored in the gno keybase (~/.config/gno), also used by gnokey. Accounts sign
realm and package deployments.`,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// point the gno keybase at --home when set (mirrors --keyring-dir)
			if home, _ := cmd.Flags().GetString(flagGnoHome); home != "" {
				if err := os.Setenv("GNOHOME", home); err != nil {
					return err
				}
			}
			return nil
		},
		Args: cobra.ExactArgs(0),
	}

	c.PersistentFlags().String(flagGnoHome, "", "directory for the gno keybase (default: ~/.config/gno)")

	c.AddCommand(
		newGnoAccountCreate(),
		newGnoAccountList(),
		newGnoAccountShow(),
		newGnoAccountDelete(),
		newGnoAccountExport(),
		newGnoAccountImport(),
	)

	return c
}

// getGnoPassphraseInteractive returns the --passphrase flag value, prompting
// interactively when unset and stdin is a terminal. Non-interactive sessions
// (scripts, CI) fall back to an empty passphrase (dev friendly).
func getGnoPassphraseInteractive(cmd *cobra.Command) (string, error) {
	pass := flagGetGnoPassphrase(cmd)
	if pass != "" {
		return pass, nil
	}
	stdinFD := int(os.Stdin.Fd())
	if !isTerminal(stdinFD) {
		return "", nil
	}
	if err := bubbleconfirm.Ask(
		bubbleconfirm.NewQuestion("Passphrase (leave empty for no passphrase)",
			&pass,
			bubbleconfirm.HideAnswer(),
			bubbleconfirm.GetConfirmation(),
		)); err != nil {
		return "", err
	}
	return pass, nil
}

func newGnoAccountCreate() *cobra.Command {
	c := &cobra.Command{
		Use:   "create <name>",
		Short: "Create a new gno account",
		Long: `Create a new account in the gno keybase (~/.config/gno).

The account is created with an empty passphrase unless --passphrase is set
(dev friendly). Use --recover to restore an account from a mnemonic.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			session := cliui.New(cliui.StartSpinnerWithText(statusCreating))
			defer session.End()

			recover, err := cmd.Flags().GetBool(flagGnoRecover)
			if err != nil {
				return err
			}
			passphrase, err := getGnoPassphraseInteractive(cmd)
			if err != nil {
				return err
			}

			if recover {
				var mnemonic string
				if err := bubbleconfirm.Ask(
					bubbleconfirm.NewQuestion("Your mnemonic", &mnemonic, bubbleconfirm.Required())); err != nil {
					return err
				}
				info, err := gno.RecoverKey(args[0], strings.TrimSpace(mnemonic), passphrase, 0, 0)
				if err != nil {
					return err
				}
				return session.Printf("%s Account %q recovered: %s\n", icons.OK, info.Name, info.Address)
			}

			info, mnemonic, err := gno.CreateKey(args[0], "", passphrase, 0, 0)
			if err != nil {
				return err
			}
			return session.Printf("%s Account %q created, keep your mnemonic in a secret place:\n\n%s\n\nAddress: %s\n",
				icons.OK, info.Name, mnemonic, info.Address)
		},
	}
	c.Flags().Bool(flagGnoRecover, false, "recover an account from a mnemonic")
	c.Flags().String(flagGnoPassphrase, "", "passphrase to encrypt the key (default: empty)")
	return c
}

// printGnoAccounts prints accounts as a table (name, address, type).
func printGnoAccounts(cmd *cobra.Command, accounts ...gno.KeyInfo) error {
	entries := make([][]string, 0, len(accounts))
	for _, acc := range accounts {
		entries = append(entries, []string{acc.Name, acc.Address, acc.Type})
	}
	return entrywriter.MustWrite(cmd.OutOrStdout(), []string{"name", "address", "type"}, entries...)
}

func newGnoAccountList() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Show a list of all gno accounts",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			keys, err := gno.ListKeys()
			if err != nil {
				return err
			}
			if len(keys) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No accounts found. Create one with `ignite account create <name>`.")
				return nil
			}
			return printGnoAccounts(cmd, keys...)
		},
	}
}

func newGnoAccountShow() *cobra.Command {
	return &cobra.Command{
		Use:   "show <name>",
		Short: "Show a gno account by name or address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			info, err := gno.ShowKey(args[0])
			if err != nil {
				return err
			}
			return printGnoAccounts(cmd, info)
		},
	}
}

func newGnoAccountDelete() *cobra.Command {
	c := &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a gno account by name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			if !getYes(cmd) {
				stdinFD := int(os.Stdin.Fd())
				if isTerminal(stdinFD) {
					var confirmed bool
					if err := bubbleconfirm.Ask(
						bubbleconfirm.NewQuestion(fmt.Sprintf("Are you sure you want to delete account %q?", name),
							&confirmed,
							bubbleconfirm.GetConfirmation(),
						)); err != nil {
						return err
					}
					if !confirmed {
						return nil
					}
				}
			}

			session := cliui.New(cliui.StartSpinnerWithText(statusDeleting))
			defer session.End()

			if err := gno.DeleteKey(name, flagGetGnoPassphrase(cmd)); err != nil {
				return err
			}
			return session.Printf("%s Account %q deleted.\n", icons.OK, name)
		},
	}
	c.Flags().BoolP(flagYes, "y", false, "skips confirmation prompt")
	return c
}

func newGnoAccountExport() *cobra.Command {
	c := &cobra.Command{
		Use:   "export <name>",
		Short: "Export a gno account as an armored private key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			session := cliui.New(cliui.StartSpinnerWithText(statusExporting))
			defer session.End()

			output, err := cmd.Flags().GetString(flagGnoOutput)
			if err != nil {
				return err
			}
			armor, err := gno.ExportKey(args[0], flagGetGnoPassphrase(cmd), output)
			if err != nil {
				return err
			}
			if armor != "" {
				session.End()
				fmt.Println(armor)
				return nil
			}
			return session.Printf("%s Key exported to %s\n", icons.OK, output)
		},
	}
	c.Flags().String(flagGnoOutput, "", "output path for the armored key (default: stdout)")
	c.Flags().String(flagGnoPassphrase, "", "passphrase used to decrypt and re-encrypt the key")
	return c
}

func newGnoAccountImport() *cobra.Command {
	c := &cobra.Command{
		Use:   "import <name> <armor-file>",
		Short: "Import a gno account from an armored private key",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			session := cliui.New(cliui.StartSpinnerWithText(statusImporting))
			defer session.End()

			info, err := gno.ImportKey(args[0], args[1], flagGetGnoPassphrase(cmd))
			if err != nil {
				return err
			}
			return session.Printf("%s Account %q imported: %s\n", icons.OK, info.Name, info.Address)
		},
	}
	c.Flags().String(flagGnoPassphrase, "", "passphrase used to decrypt the armor")
	return c
}
