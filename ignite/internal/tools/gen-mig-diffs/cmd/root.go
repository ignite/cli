package cmd

import (
	semver "github.com/Masterminds/semver/v3"
	"github.com/ignite/cli/v28/ignite/internal/tools/gen-mig-diffs/migdiff"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	fromFlag   = "from"
	toFlag     = "to"
	outputFlag = "output"

	igniteCliRepository = "http://github.com/ignite/cli.git"
	igniteRepoPath      = "src/github.com/ignite/cli"
	igniteBinaryPath    = "dist/ignite"
)

// NewRootCmd creates a new root command
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen-mig-diffs",
		Short: "Generate migration diffs",
		Long:  "This tool is used to generate migration diff files for each of ignites scaffold commands",
		RunE: func(cmd *cobra.Command, args []string) error {
			from, _ := cmd.Flags().GetString(fromFlag)
			to, _ := cmd.Flags().GetString(toFlag)
			output, _ := cmd.Flags().GetString(outputFlag)

			fromVer, err := semver.NewVersion(from)
			if err != nil && from != "" {
				return errors.Wrapf(err, "failed to parse from version %s", from)
			}
			toVer, err := semver.NewVersion(to)
			if err != nil && to != "" {
				return errors.Wrapf(err, "failed to parse to version %s", to)
			}

			mdg, err := migdiff.NewMigDiffGenerator(fromVer, toVer)
			if err != nil {
				return err
			}

			return mdg.Generate(output)
		},
	}

	cmd.Flags().StringP(fromFlag, "f", "", "Version of ignite or path to ignite source code to generate the diff from")
	cmd.Flags().StringP(toFlag, "t", "", "Version of ignite or path to ignite source code to generate the diff to")
	cmd.Flags().StringP(outputFlag, "o", ".", "Output directory to save the migration diff files")

	return cmd
}
