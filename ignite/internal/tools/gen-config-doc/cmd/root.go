package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	v0 "github.com/ignite/cli/v29/ignite/config/chain/v0"
	v1 "github.com/ignite/cli/v29/ignite/config/chain/v1"
	"github.com/ignite/cli/v29/ignite/pkg/clidoc"
	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"

	"github.com/ignite/cli/ignite/internal/tools/gen-config-doc/templates/doc"
)

const (
	flagVersion = "version"
	flagOutput  = "output"

	defaultDocPath = "docs/docs/08-references"
)

// NewRootCmd creates a new root command.
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen-config-doc",
		Short: "generate configuration file documentation",
		Long:  "This tool is used to generate the chain configuration file documentation",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			session := cliui.New()
			defer session.End()

			var (
				version, _ = cmd.Flags().GetString(flagVersion)
				output, _  = cmd.Flags().GetString(flagOutput)
			)

			output, err = filepath.Abs(output)
			if err != nil {
				return errors.Wrap(err, "failed to find the abs path")
			}

			var docs clidoc.Docs
			switch version {
			case "v0":
				docs, err = clidoc.GenDoc(v0.Config{})
			case "v1":
				docs, err = clidoc.GenDoc(v1.Config{})
			default:
				return errors.Errorf("unknown version: %s", version)
			}
			if err != nil {
				return errors.Wrapf(err, "failed to generate migration doc %s", version)
			}

			// Generate the docs file.
			g, err := doc.NewGenerator(doc.Options{
				Path:   output,
				Config: docs.String(),
			})
			if err != nil {
				return errors.Wrap(err, "failed to create the doc generator object")
			}

			sm, err := xgenny.NewRunner(cmd.Context(), output).RunAndApply(g)
			if err != nil {
				return err
			}

			files := append(sm.CreatedFiles(), sm.ModifiedFiles()...)
			if len(files) == 0 {
				return errors.Errorf("config doc not created at %s", output)
			}
			session.EventBus().SendInfo(
				fmt.Sprintf("Config doc generated successfully at %s", files[0]),
			)

			return nil
		},
	}

	cmd.Flags().StringP(flagVersion, "v", "v1", "Version of Ignite config file")
	cmd.Flags().StringP(flagOutput, "o", defaultDocPath, "Output directory to save the config document")

	return cmd
}
