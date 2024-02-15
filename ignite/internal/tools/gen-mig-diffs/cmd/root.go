package cmd

import (
	semver "github.com/Masterminds/semver/v3"
	"github.com/spf13/cobra"

	"github.com/ignite/cli/v28/ignite/internal/tools/gen-mig-diffs/migdiff"
	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
)

const (
	fromFlag   = "from"
	toFlag     = "to"
	sourceFlag = "source"
	outputFlag = "output"
)

// NewRootCmd creates a new root command.
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen-mig-diffs",
		Short: "Generate migration diffs",
		Long:  "This tool is used to generate migration diff files for each of ignites scaffold commands",
		RunE: func(cmd *cobra.Command, args []string) error {
			from, _ := cmd.Flags().GetString(fromFlag)
			to, _ := cmd.Flags().GetString(toFlag)
			source := cmd.Flag(sourceFlag).Value.String()
			output, _ := cmd.Flags().GetString(outputFlag)

			fromVer, err := semver.NewVersion(from)
			if err != nil && from != "" {
				return errors.Wrapf(err, "failed to parse from version %s", from)
			}
			toVer, err := semver.NewVersion(to)
			if err != nil && to != "" {
				return errors.Wrapf(err, "failed to parse to version %s", to)
			}

			session := cliui.New()
			defer session.End()
			mdg, err := migdiff.NewGenerator(fromVer, toVer, source, session)
			if err != nil {
				return err
			}
			defer mdg.Cleanup()

			return mdg.Generate(output)
		},
	}

	cmd.Flags().StringP(fromFlag, "f", "", "Version of ignite or path to ignite source code to generate the diff from")
	cmd.Flags().StringP(toFlag, "t", "", "Version of ignite or path to ignite source code to generate the diff to")
	cmd.Flags().StringP(sourceFlag, "s", "", "Path to ignite source code repository (optional)")
	cmd.Flags().StringP(outputFlag, "o", "./diffs", "Output directory to save the migration diff files")

	return cmd
}
