package ignitecmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/v29/ignite/services/gno"
)

// NewGnoGenerate returns the `ignite generate` command for gno.land packages.
func NewGnoGenerate() *cobra.Command {
	c := &cobra.Command{
		Use:     "generate [command]",
		Aliases: []string{"g"},
		Short:   "Generate clients for gno.land packages",
		Args:    cobra.ExactArgs(0),
	}

	c.AddCommand(
		newGnoGenerateTSClient(),
	)

	return c
}

func newGnoGenerateTSClient() *cobra.Command {
	c := &cobra.Command{
		Use:   "ts-client [dir]",
		Short: "Generate a typed TypeScript client for a gno package",
		Long: `Generate a typed TypeScript client for the gno package at dir (default:
current directory).

The client is a realm module in the GnoWallet.addRealm shape of
@gnolang/gno-js-client. Read-only functions evaluate expressions through the
wallet's provider; realm functions (state-mutating) broadcast transactions
with wallet.callMethod.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			session := cliui.New(cliui.StartSpinnerWithText(statusGenerating))
			defer session.End()

			dir := dirArg(args)
			sig, err := gno.ExtractSignatures(dir)
			if err != nil {
				return err
			}
			out, err := cmd.Flags().GetString(flagGnoOutput)
			if err != nil {
				return err
			}
			if out == "" {
				out = filepath.Join(dir, sig.Name+".client.ts")
			}

			src := gno.GenerateRealmClient(sig)
			if err := os.WriteFile(out, []byte(src), 0o644); err != nil { //nolint:gosec // generated file is world-readable like project sources
				return err
			}
			return session.Printf("%s TypeScript client written to %s\n", icons.OK, out)
		},
	}
	c.Flags().String(flagGnoOutput, "", "output file (default: <dir>/<name>.client.ts)")
	return c
}

// dirArg returns args[0] or the current directory.
func dirArg(args []string) string {
	if len(args) == 1 {
		return args[0]
	}
	return "."
}
