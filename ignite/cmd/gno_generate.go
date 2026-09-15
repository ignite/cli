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
		Short:   "Generate clients and interface files for gno.land packages",
		Args:    cobra.ExactArgs(0),
	}

	c.AddCommand(
		newGnoGenerateIDL(),
		newGnoGenerateTSClient(),
	)

	return c
}

func newGnoGenerateIDL() *cobra.Command {
	c := &cobra.Command{
		Use:   "idl [dir]",
		Short: "Generate an IDL (interface definition) file for a gno package",
		Long: `Generate an IDL JSON file for the gno package at dir (default: current
directory), describing its callable functions and their arguments. The IDL
is the input of client generators.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			session := cliui.New(cliui.StartSpinnerWithText(statusGenerating))
			defer session.End()

			dir := dirArg(args)
			out, err := cmd.Flags().GetString(flagGnoOutput)
			if err != nil {
				return err
			}
			if out == "" {
				out = filepath.Join(dir, "idl.json")
			}

			idl, err := gno.GenerateIDL(dir)
			if err != nil {
				return err
			}
			data, err := idl.WriteIDLJSON()
			if err != nil {
				return err
			}
			if err := os.WriteFile(out, data, 0o644); err != nil { //nolint:gosec // generated file is world-readable like project sources
				return err
			}
			return session.Printf("%s IDL written to %s (%d functions)\n", icons.OK, out, len(idl.Instructions))
		},
	}
	c.Flags().String(flagGnoOutput, "", "output file (default: <dir>/idl.json)")
	return c
}

func newGnoGenerateTSClient() *cobra.Command {
	c := &cobra.Command{
		Use:   "ts-client [dir]",
		Short: "Generate a TypeScript client for a gno package",
		Long: `Generate a TypeScript client for the gno package at dir (default:
current directory).

Read-only functions query the chain through the vm/qeval_json endpoint;
realm functions (state-mutating) go through an injected GnoSigner, so any
wallet or gnokey-based signer can be plugged in.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			session := cliui.New(cliui.StartSpinnerWithText(statusGenerating))
			defer session.End()

			dir := dirArg(args)
			idlg, err := gno.GenerateIDL(dir)
			if err != nil {
				return err
			}
			out, err := cmd.Flags().GetString(flagGnoOutput)
			if err != nil {
				return err
			}
			if out == "" {
				out = filepath.Join(dir, idlg.Name+".client.ts")
			}

			src := gno.GenerateTSClient(idlg)
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
