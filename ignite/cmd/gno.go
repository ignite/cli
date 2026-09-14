package ignitecmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/ignite/cli/v29/ignite/services/gno"
)

const (
	flagGnoRemote     = "remote"
	flagGnoFrom       = "from"
	flagGnoChainID    = "chain-id"
	flagGnoGasWanted  = "gas-wanted"
	flagGnoGasFee     = "gas-fee"
	flagGnoPassphrase = "passphrase"
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
	c := &cobra.Command{
		Use:   "realm <name>",
		Short: "Scaffold a new gno.land realm (stateful smart contract)",
		Long: `Scaffold a new gno.land realm.

A realm is a stateful smart contract: package-level variables are persisted
on-chain. <name> is either a bare name ("counter", deployed as
gno.land/r/counter) or a full path ("gno.land/r/demo/counter").`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, modulePath, err := gno.Scaffold(gno.KindRealm, args[0], gno.ScaffoldOptions{})
			if err != nil {
				return err
			}
			fmt.Printf("⭐ Realm scaffolded.\n\n  Directory: %s\n  Module:   %s\n", dir, modulePath)
			return nil
		},
	}
	return c
}

func newGnoScaffoldPackage() *cobra.Command {
	return &cobra.Command{
		Use:   "package <name>",
		Short: "Scaffold a new gno.land package (stateless library)",
		Long: `Scaffold a new gno.land package.

A package is a stateless library of pure functions (gno.land/p/...).`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, modulePath, err := gno.Scaffold(gno.KindPackage, args[0], gno.ScaffoldOptions{})
			if err != nil {
				return err
			}
			fmt.Printf("⭐ Package scaffolded.\n\n  Directory: %s\n  Module:   %s\n", dir, modulePath)
			return nil
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
			return gno.Deploy(dir, gno.DeployOptions{
				From:       flagGetGnoFrom(cmd),
				PkgPath:    pkgPath,
				Remote:     flagGetGnoRemote(cmd),
				ChainID:    flagGetGnoChainID(cmd),
				GasWanted:  gasWanted,
				GasFee:     flagGetGnoGasFee(cmd),
				Passphrase: flagGetGnoPassphrase(cmd),
				MaxDeposit: maxDeposit,
			})
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
		Short:   "Manage gno.land accounts (gno keybase)",
		Args:    cobra.ExactArgs(0),
	}

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

func newGnoAccountCreate() *cobra.Command {
	c := &cobra.Command{
		Use:   "create <name>",
		Short: "Create a new gno account",
		Long: `Create a new account in the gno keybase (~/.config/gno).

The account is created with an empty passphrase unless --passphrase is set
(dev friendly). Use --recover to restore an account from a mnemonic.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			recover, err := cmd.Flags().GetBool(flagGnoRecover)
			if err != nil {
				return err
			}
			passphrase := flagGetGnoPassphrase(cmd)

			if recover {
				var mnemonic string
				fmt.Print("Enter your mnemonic:\n> ")
				if _, err := fmt.Scanln(&mnemonic); err != nil {
					return fmt.Errorf("reading mnemonic: %w", err)
				}
				info, err := gno.RecoverKey(args[0], strings.TrimSpace(mnemonic), passphrase, 0, 0)
				if err != nil {
					return err
				}
				fmt.Printf("🔑 Account recovered: %s (%s)\n", info.Name, info.Address)
				return nil
			}

			info, mnemonic, err := gno.CreateKey(args[0], "", passphrase, 0, 0)
			if err != nil {
				return err
			}
			fmt.Printf("🔑 Account created: %s (%s)\n\nMnemonic (save it somewhere safe):\n%s\n", info.Name, info.Address, mnemonic)
			return nil
		},
	}
	c.Flags().Bool(flagGnoRecover, false, "recover an account from a mnemonic")
	c.Flags().String(flagGnoPassphrase, "", "passphrase to encrypt the key (default: empty)")
	return c
}

func newGnoAccountList() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List gno accounts",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			keys, err := gno.ListKeys()
			if err != nil {
				return err
			}
			if len(keys) == 0 {
				fmt.Println("No accounts found. Create one with `ignite account create <name>`.")
				return nil
			}
			for _, k := range keys {
				fmt.Printf("%s\t%s\t%s\n", k.Name, k.Address, k.Type)
			}
			return nil
		},
	}
}

func newGnoAccountShow() *cobra.Command {
	return &cobra.Command{
		Use:   "show <name>",
		Short: "Show a gno account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			info, err := gno.ShowKey(args[0])
			if err != nil {
				return err
			}
			fmt.Printf("%s\t%s\t%s\n", info.Name, info.Address, info.Type)
			return nil
		},
	}
}

func newGnoAccountDelete() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a gno account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return gno.DeleteKey(args[0], flagGetGnoPassphrase(cmd))
		},
	}
}

func newGnoAccountExport() *cobra.Command {
	c := &cobra.Command{
		Use:   "export <name>",
		Short: "Export a gno account as an armored private key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			output, err := cmd.Flags().GetString(flagGnoOutput)
			if err != nil {
				return err
			}
			armor, err := gno.ExportKey(args[0], flagGetGnoPassphrase(cmd), output)
			if err != nil {
				return err
			}
			if armor != "" {
				fmt.Println(armor)
			} else {
				fmt.Printf("🔐 Key exported to %s\n", output)
			}
			return nil
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
			info, err := gno.ImportKey(args[0], args[1], flagGetGnoPassphrase(cmd))
			if err != nil {
				return err
			}
			fmt.Printf("🔑 Account imported: %s (%s)\n", info.Name, info.Address)
			return nil
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
