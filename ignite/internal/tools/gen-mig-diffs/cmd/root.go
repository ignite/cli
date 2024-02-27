package cmd

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/spf13/cobra"

	"github.com/ignite/cli/v28/ignite/internal/tools/gen-mig-diffs/pkg/diff"
	"github.com/ignite/cli/v28/ignite/internal/tools/gen-mig-diffs/pkg/repo"
	"github.com/ignite/cli/v28/ignite/internal/tools/gen-mig-diffs/pkg/scaffold"
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
		Short: "GenerateBinaries migration diffs",
		Long:  "This tool is used to generate migration diff files for each of ignites scaffold commands",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				from, _   = cmd.Flags().GetString(fromFlag)
				to, _     = cmd.Flags().GetString(toFlag)
				source, _ = cmd.Flags().GetString(sourceFlag)
				output, _ = cmd.Flags().GetString(outputFlag)
			)

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

			igniteRepo, err := repo.New(fromVer, toVer, session, repo.WithSource(source))
			if err != nil {
				return err
			}
			defer igniteRepo.Cleanup()

			fromBin, toBin, err := igniteRepo.GenerateBinaries()
			if err != nil {
				return err
			}

			sFrom := scaffold.New(fromBin, scaffold.DefaultCommands)
			session.StartSpinner(fmt.Sprintf("Running scaffold commands for v%s...", fromVer.String()))
			if err := sFrom.Run(fromVer, output); err != nil {
				return err
			}

			sTo := scaffold.New(toBin, scaffold.DefaultCommands)
			session.StartSpinner(fmt.Sprintf("Running scaffold commands for v%s...", toVer.String()))
			if err := sTo.Run(toVer, output); err != nil {
				return err
			}

			session.StopSpinner()
			session.EventBus().SendInfo(fmt.Sprintf("Scaffolded code for commands at %s", output))

			session.StartSpinner("Calculating diff...")
			diffs, err := diff.CalculateDiffs(fromDir, toDir)
			if err != nil {
				return errors.Wrap(err, "failed to calculate diff")
			}
			session.StopSpinner()
			session.EventBus().SendInfo("Diff calculated successfully")

			if err = diff.SaveDiffs(diffs, output); err != nil {
				return errors.Wrap(err, "failed to save diff map")
			}
			session.Println("Migration diffs generated successfully at", output)

			return nil
		},
	}

	cmd.Flags().StringP(fromFlag, "f", "", "Version of ignite or path to ignite source code to generate the diff from")
	cmd.Flags().StringP(toFlag, "t", "", "Version of ignite or path to ignite source code to generate the diff to")
	cmd.Flags().StringP(sourceFlag, "s", "", "Path to ignite source code repository (optional)")
	cmd.Flags().StringP(outputFlag, "o", "./diffs", "Output directory to save the migration diff files")

	return cmd
}
