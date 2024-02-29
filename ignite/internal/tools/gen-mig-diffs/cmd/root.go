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
	flagFrom           = "from"
	flagTo             = "to"
	flagOutput         = "output"
	flagSource         = "repo-source"
	flagRepoURL        = "repo-url"
	flagRepoOutput     = "repo-output"
	flagRepoCleanup    = "repo-cleanup"
	flagScaffoldOutput = "scaffold-output"
	flagScaffoldCache  = "scaffold-cache"
)

// NewRootCmd creates a new root command.
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen-mig-diffs",
		Short: "GenerateBinaries migration diffs",
		Long:  "This tool is used to generate migration diff files for each of ignites scaffold commands",
		RunE: func(cmd *cobra.Command, args []string) error {
			session := cliui.New()
			defer session.End()

			var (
				from, _           = cmd.Flags().GetString(flagFrom)
				to, _             = cmd.Flags().GetString(flagTo)
				repoSource, _     = cmd.Flags().GetString(flagSource)
				output, _         = cmd.Flags().GetString(flagOutput)
				repoURL, _        = cmd.Flags().GetString(flagRepoURL)
				repoOutput, _     = cmd.Flags().GetString(flagRepoOutput)
				repoCleanup, _    = cmd.Flags().GetBool(flagRepoCleanup)
				scaffoldOutput, _ = cmd.Flags().GetString(flagScaffoldOutput)
				scaffoldCache, _  = cmd.Flags().GetString(flagScaffoldCache)
			)
			fromVer, err := semver.NewVersion(from)
			if err != nil && from != "" {
				return errors.Wrapf(err, "failed to parse from version %s", from)
			}
			toVer, err := semver.NewVersion(to)
			if err != nil && to != "" {
				return errors.Wrapf(err, "failed to parse to version %s", to)
			}

			repoOptions := make([]repo.Options, 0)
			if repoCleanup {
				repoOptions = append(repoOptions, repo.WithCleanup())
			}
			if repoSource != "" {
				repoOptions = append(repoOptions, repo.WithSource(repoSource))
			}
			if repoURL != "" {
				repoOptions = append(repoOptions, repo.WithRepoURL(repoURL))
			}
			if repoOutput != "" {
				repoOptions = append(repoOptions, repo.WithRepoOutput(repoOutput))
			}

			igniteRepo, err := repo.New(fromVer, toVer, session, repoOptions...)
			if err != nil {
				return err
			}
			defer igniteRepo.Cleanup()

			fromBin, toBin, err := igniteRepo.GenerateBinaries(cmd.Context())
			if err != nil {
				return err
			}

			scaffoldOptions := make([]scaffold.Options, 0)
			if scaffoldOutput != "" {
				scaffoldOptions = append(scaffoldOptions, scaffold.WithOutput(scaffoldOutput))
			}
			if scaffoldCache != "" {
				scaffoldOptions = append(scaffoldOptions, scaffold.WithCachePath(scaffoldCache))
			}

			session.StartSpinner(fmt.Sprintf("Running scaffold commands for %s...", igniteRepo.From.Original()))
			fromDir, err := scaffold.Run(cmd.Context(), fromBin, igniteRepo.From, scaffoldOptions...)
			if err != nil {
				return err
			}

			session.StartSpinner(fmt.Sprintf("Running scaffold commands for %s...", igniteRepo.To.Original()))
			toDir, err := scaffold.Run(cmd.Context(), toBin, igniteRepo.To, scaffoldOptions...)
			if err != nil {
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

	cmd.Flags().StringP(flagFrom, "f", "", "Version of ignite or path to ignite source code to generate the diff from")
	cmd.Flags().StringP(flagTo, "t", "", "Version of ignite or path to ignite source code to generate the diff to")
	cmd.Flags().StringP(flagOutput, "o", "./diffs", "Output directory to save the migration diff files")
	cmd.Flags().StringP(flagSource, "s", "", "Path to ignite source code repository. Set the source automatically set the cleanup to false")
	cmd.Flags().String(flagRepoURL, "", "Git URL for the Ignite repository")
	cmd.Flags().String(flagRepoOutput, "", "Output path to clone the ignite repository")
	cmd.Flags().Bool(flagRepoCleanup, true, "Cleanup the repository path after use")
	cmd.Flags().String(flagScaffoldOutput, "", "Output path to clone the ignite repository")
	cmd.Flags().String(flagScaffoldCache, "", "Output path to clone the ignite repository")

	return cmd
}
