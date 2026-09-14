package ignitecmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/cliui/bubbleconfirm"
	"github.com/ignite/cli/v29/ignite/pkg/cliui/entrywriter"
	"github.com/ignite/cli/v29/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/v29/ignite/services/gno"
)

const (
	flagGnoRemote     = "remote"
	flagGnoFrom       = "from"
	flagGnoChainID    = "chain-id"
	flagGnoGasWanted  = "gas-wanted"
	flagGnoGasFee     = "gas-fee"
	flagGnoPassphrase = "passphrase"
	flagGnoHome       = "home"
)

// NewGnoScaffold returns the `ignite scaffold` command for gno.land
// realms and packages.
func NewGnoScaffold() *cobra.Command {
	c := &cobra.Command{
		Use:     "scaffold [command]",
		Aliases: []string{"s"},
		Short:   "Scaffold gno.land realms and packages",
		Long:    `Scaffold new gno.land smart contracts: stateful realms (r/) and stateless packages (p/).`,
		Args:    cobra.ExactArgs(0),
	}

	c.AddCommand(
		newGnoScaffoldRealm(),
		newGnoScaffoldPackage(),
	)

	return c
}

func newGnoScaffoldRealm() *cobra.Command {
	return &cobra.Command{
		Use:   "realm <name>",
		Short: "Scaffold a new gno.land realm (stateful smart contract)",
		Long: `Scaffold a new gno.land realm.

A realm is a stateful smart contract: package-level variables are persisted
on-chain. <name> is either a bare name ("counter", deployed as
gno.land/r/counter) or a full path ("gno.land/r/demo/counter").`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			session := cliui.New(cliui.StartSpinnerWithText(statusCreating))
			defer session.End()

			dir, modulePath, err := gno.Scaffold(gno.KindRealm, args[0], gno.ScaffoldOptions{})
			if err != nil {
				return err
			}
			return session.Printf("%s Realm scaffolded.\n\n  Directory: %s\n  Module:   %s\n", icons.OK, dir, modulePath)
		},
	}
}

func newGnoScaffoldPackage() *cobra.Command {
	return &cobra.Command{
		Use:   "package <name>",
		Short: "Scaffold a new gno.land package (stateless library)",
		Long: `Scaffold a new gno.land package.

A package is a stateless library of pure functions (gno.land/p/...).`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			session := cliui.New(cliui.StartSpinnerWithText(statusCreating))
			defer session.End()

			dir, modulePath, err := gno.Scaffold(gno.KindPackage, args[0], gno.ScaffoldOptions{})
			if err != nil {
				return err
			}
			return session.Printf("%s Package scaffolded.\n\n  Directory: %s\n  Module:   %s\n", icons.OK, dir, modulePath)
		},
	}
}

// NewGnoChain returns the `ignite chain` command for gno.land chains.
func NewGnoChain() *cobra.Command {
	c := &cobra.Command{
		Use:     "chain [command]",
		Aliases: []string{"c"},
		Short:   "Run and deploy to gno.land chains",
		Args:    cobra.ExactArgs(0),
	}

	c.AddCommand(
		newGnoChainServe(),
		newGnoChainDeploy(),
	)

	return c
}

func newGnoChainServe() *cobra.Command {
	c := &cobra.Command{
		Use:   "serve",
		Short: "Start a local gno.land dev chain",
		Long: `Start a local in-memory gno.land dev chain.

The chain pre-funds the well-known dev account (test1), preloads any gno
package found in the current directory and reloads the chain on every .gno
file change. The RPC endpoint listens on tcp://127.0.0.1:26657 by default.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			opts := gno.ServeOptions{
				ChainID:        flagGetGnoChainID(cmd),
				RPCListener:    flagGetGnoRemote(cmd),
				MaxGasPerBlock: flagGetGnoMaxGas(cmd),
			}
			dir, err := os.Getwd()
			if err != nil {
				return err
			}

			ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
			defer cancel()

			return gno.Serve(ctx, dir, opts, cmd.OutOrStdout())
		},
	}

	c.Flags().String(flagGnoChainID, "dev", "chain id")
	c.Flags().String(flagGnoRemote, "tcp://127.0.0.1:26657", "node RPC listen address")
	c.Flags().Int64(flagGnoMaxGas, 10_000_000_000, "max gas per block")
	return c
}

func newGnoChainDeploy() *cobra.Command {
	c := &cobra.Command{
		Use:   "deploy [dir]",
		Short: "Deploy a gno realm or package to a chain",
		Long: `Deploy the gno package at dir (default: current directory) to a gno.land chain.

The module path is read from gnomod.toml. The package is signed with the key
given by --from and broadcast to --remote (default: local dev chain started
with ` + "`ignite chain serve`" + `).`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "."
			if len(args) == 1 {
				dir = args[0]
			}
			gasWanted, err := cmd.Flags().GetInt64(flagGnoGasWanted)
			if err != nil {
				return err
			}
			pkgPath, err := cmd.Flags().GetString(flagGnoPkgPath)
			if err != nil {
				return err
			}
			maxDeposit, err := cmd.Flags().GetString(flagGnoMaxDeposit)
			if err != nil {
				return err
			}
			if err := gno.Deploy(dir, gno.DeployOptions{
				From:       flagGetGnoFrom(cmd),
				PkgPath:    pkgPath,
				Remote:     flagGetGnoRemote(cmd),
				ChainID:    flagGetGnoChainID(cmd),
				GasWanted:  gasWanted,
				GasFee:     flagGetGnoGasFee(cmd),
				Passphrase: flagGetGnoPassphrase(cmd),
				MaxDeposit: maxDeposit,
			}); err != nil {
				return err
			}
			return nil
		},
	}

	c.Flags().String(flagGnoFrom, "", "key name or bech32 address signing the deployment (default: dev account test1)")
	c.Flags().String(flagGnoPkgPath, "", "override the module path from gnomod.toml")
	c.Flags().String(flagGnoRemote, "127.0.0.1:26657", "chain RPC address")
	c.Flags().String(flagGnoChainID, "dev", "chain id")
	c.Flags().Int64(flagGnoGasWanted, gno.DefaultGasWanted, "gas requested for the tx")
	c.Flags().String(flagGnoGasFee, "", fmt.Sprintf("gas payment fee (default: %d%s)", gno.DefaultGasFee, "ugnot"))
	c.Flags().String(flagGnoMaxDeposit, "", "max storage deposit coins")
	c.Flags().String(flagGnoPassphrase, "", "passphrase to unlock the signing key (empty for dev keys)")
	return c
}

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
				stdinFD := int(os.Stdin.Fd())
				if !isTerminal(stdinFD) {
					return fmt.Errorf("--recover requires an interactive terminal; pipe the mnemonic via --recover with a TTY or use gnokey")
				}
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

// flag names used across gno commands.
const (
	flagGnoPkgPath    = "pkg-path"
	flagGnoMaxDeposit = "max-deposit"
	flagGnoMaxGas     = "max-gas"
	flagGnoRecover    = "recover"
	flagGnoOutput     = "output"
)

func flagGetGnoFrom(cmd *cobra.Command) string {
	from, _ := cmd.Flags().GetString(flagGnoFrom)
	return from
}

func flagGetGnoRemote(cmd *cobra.Command) string {
	remote, _ := cmd.Flags().GetString(flagGnoRemote)
	return remote
}

func flagGetGnoChainID(cmd *cobra.Command) string {
	chainID, _ := cmd.Flags().GetString(flagGnoChainID)
	return chainID
}

func flagGetGnoGasFee(cmd *cobra.Command) string {
	gasFee, _ := cmd.Flags().GetString(flagGnoGasFee)
	return gasFee
}

func flagGetGnoPassphrase(cmd *cobra.Command) string {
	passphrase, _ := cmd.Flags().GetString(flagGnoPassphrase)
	return passphrase
}

func flagGetGnoMaxGas(cmd *cobra.Command) int64 {
	maxGas, _ := cmd.Flags().GetInt64(flagGnoMaxGas)
	return maxGas
}
