package ignitecmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/v29/ignite/services/gno"
)

// NewGnoChain returns the `ignite chain` command for gno.land chains.
func NewGnoChain() *cobra.Command {
	c := &cobra.Command{
		Use:     "chain [command]",
		Aliases: []string{"c"},
		Short:   "Run, deploy and interact with gno.land chains",
		Args:    cobra.ExactArgs(0),
	}

	c.AddCommand(
		newGnoChainServe(),
		newGnoChainDeploy(),
		newGnoChainCall(),
		newGnoChainQuery(),
		newGnoChainSend(),
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
			pkgPath, err := cmd.Flags().GetString(flagGnoPkgPath)
			if err != nil {
				return err
			}
			maxDeposit, err := cmd.Flags().GetString(flagGnoMaxDeposit)
			if err != nil {
				return err
			}
			return gno.Deploy(dir, gno.DeployOptions{
				TxBaseOptions: gnoTxBaseFrom(cmd),
				PkgPath:       pkgPath,
				MaxDeposit:    maxDeposit,
			})
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
			return gno.Call(gno.CallOptions{
				TxBaseOptions: gnoTxBaseFrom(cmd),
				PkgPath:       args[0],
				Func:          args[1],
				Args:          args[2:],
				Send:          send,
			})
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

Example:

$ ignite chain query "gno.land/r/counter.Get()"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			session := cliui.New(cliui.StartSpinnerWithText("Querying..."))
			defer session.End()

			res, err := gno.Query(flagGetGnoRemote(cmd), args[0])
			if err != nil {
				return err
			}
			return session.Printf("%s %s\n", icons.OK, res)
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
			return gno.Send(gno.SendOptions{
				TxBaseOptions: gnoTxBaseFrom(cmd),
				To:            args[0],
				Amount:        args[1],
			})
		},
	}

	gnoTxBaseFlags(c)
	return c
}

const flagGnoSend = "send"
