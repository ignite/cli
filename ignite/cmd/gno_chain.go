package ignitecmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/ignite/cli/v29/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/v29/ignite/services/gno"
)

// NewGnoChain returns the `ignite chain` command for gno.land chains.
func NewGnoChain() *cobra.Command {
	c := &cobra.Command{
		Use:     "chain [command]",
		Aliases: []string{"c"},
		Short:   "Run, deploy and interact with gno.land chains",
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
		newGnoChainServe(),
		newGnoChainDeploy(),
		newGnoChainCall(),
		newGnoChainQuery(),
		newGnoChainSend(),
		newGnoChainTest(),
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
			dir := dirArg(args)
			tb, err := gnoTxBaseFrom(cmd)
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
			deployedPath, res, err := gno.Deploy(dir, gno.DeployOptions{
				TxBaseOptions: tb,
				PkgPath:       pkgPath,
				MaxDeposit:    maxDeposit,
			})
			if err != nil {
				return err
			}
			return printTxDone(cmd, fmt.Sprintf("Deployed %s", deployedPath), res)
		},
	}

	gnoTxBaseFlags(c)
	c.Flags().String(flagGnoPkgPath, "", "override the module path from gnomod.toml")
	c.Flags().String(flagGnoMaxDeposit, "", "max storage deposit coins")
	return c
}

func newGnoChainCall() *cobra.Command {
	c := &cobra.Command{
		Use:   "call <pkgpath> <func> [args...]",
		Short: "Call a function of a deployed realm",
		Long: `Call a function of a deployed realm and wait for the result.

Example: fund and increment a counter realm:

$ ignite chain call gno.land/r/counter Increment
$ ignite chain call gno.land/r/counter Set 42 --send 1000000ugnot`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			send, err := cmd.Flags().GetString(flagGnoSend)
			if err != nil {
				return err
			}
			tb, err := gnoTxBaseFrom(cmd)
			if err != nil {
				return err
			}
			res, err := gno.Call(gno.CallOptions{
				TxBaseOptions: tb,
				PkgPath:       args[0],
				Func:          args[1],
				Args:          args[2:],
				Send:          send,
			})
			if err != nil {
				return err
			}
			return printTxDone(cmd, fmt.Sprintf("Called %s.%s", args[0], args[1]), res)
		},
	}

	gnoTxBaseFlags(c)
	c.Flags().String(flagGnoSend, "", "coins to send along with the call")
	return c
}

func newGnoChainQuery() *cobra.Command {
	c := &cobra.Command{
		Use:   "query <expression>",
		Short: "Evaluate a read-only expression on a chain",
		Long: `Evaluate a read-only gno expression on a gno.land chain and print the result.

An expression that names a function without parentheses is evaluated as a
call, e.g. "gno.land/r/helloworld.Get" is queried as "gno.land/r/helloworld.Get()".

Example:

$ ignite chain query "gno.land/r/counter.Get()"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			res, err := gno.Query(flagGetGnoRemote(cmd), args[0])
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), res)
			return nil
		},
	}

	c.Flags().String(flagGnoRemote, gnoDefaultRemote, "chain RPC address")
	return c
}

func newGnoChainSend() *cobra.Command {
	c := &cobra.Command{
		Use:   "send <to> <amount>",
		Short: "Send coins to an account",
		Long: `Send coins from the account given by --from to another account.

<to> is a bech32 address or a key name from the keybase. Handy on dev chains
to fund accounts created with ` + "`ignite account create`" + `.

Example:

$ ignite chain send g1... 10000000ugnot
$ ignite chain send alice 10000000ugnot --from test1`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			tb, err := gnoTxBaseFrom(cmd)
			if err != nil {
				return err
			}
			res, err := gno.Send(gno.SendOptions{
				TxBaseOptions: tb,
				To:            args[0],
				Amount:        args[1],
			})
			if err != nil {
				return err
			}
			return printTxDone(cmd, fmt.Sprintf("Sent %s to %s", args[1], args[0]), res)
		},
	}

	gnoTxBaseFlags(c)
	return c
}

const flagGnoSend = "send"

// printTxDone reports a broadcast tx on the command output.
func printTxDone(cmd *cobra.Command, action string, res *gno.BroadcastResult) error {
	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s %s (gas used: %d)\n", icons.OK, action, res.GasUsed); err != nil {
		return err
	}
	_, err := fmt.Fprintf(cmd.OutOrStdout(), "%s Tx hash: %s\n", icons.Info, res.TxHash)
	return err
}

func newGnoChainTest() *cobra.Command {
	return &cobra.Command{
		Use:   "test [dir]",
		Short: "Run the gno tests of a package",
		Long:  "Run the gno tests (_test.gno files) of the package at dir (default: current directory).",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return gno.Test(args, cmd.OutOrStdout(), cmd.OutOrStderr())
		},
	}
}
