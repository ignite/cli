package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/Masterminds/semver/v3"
	"github.com/spf13/cobra"

	"github.com/ignite/cli/v28/ignite/internal/tools/gen-mig-diffs/pkg/diff"
	"github.com/ignite/cli/v28/ignite/internal/tools/gen-mig-diffs/pkg/repo"
	"github.com/ignite/cli/v28/ignite/internal/tools/gen-mig-diffs/pkg/scaffold"
	"github.com/ignite/cli/v28/ignite/internal/tools/gen-mig-diffs/templates/doc"
	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/pkg/placeholder"
	"github.com/ignite/cli/v28/ignite/pkg/xgenny"
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

	defaultDocPath = "docs/docs/06-migration"
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

			// Check or download the source and generate the binaries for each version.
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

			releaseDescription, err := igniteRepo.ReleaseDescription()
			if err != nil {
				return errors.Wrapf(err, "failed to fetch the release tag %s description", igniteRepo.To.Original())
			}

			fromBin, toBin, err := igniteRepo.GenerateBinaries(cmd.Context())
			if err != nil {
				return err
			}

			// Scaffold the default commands for each version.
			scaffoldOptions := make([]scaffold.Options, 0)
			if scaffoldOutput != "" {
				scaffoldOptions = append(scaffoldOptions, scaffold.WithOutput(scaffoldOutput))
			}
			if scaffoldCache != "" {
				scaffoldOptions = append(scaffoldOptions, scaffold.WithCachePath(scaffoldCache))
			}

			session.StartSpinner(fmt.Sprintf("Running scaffold commands for %s...", igniteRepo.From.Original()))
			sFrom, err := scaffold.New(fromBin, igniteRepo.From, scaffoldOptions...)
			if err != nil {
				return err
			}
			defer sFrom.Cleanup()

			if err := sFrom.Run(cmd.Context()); err != nil {
				return err
			}
			session.StopSpinner()
			session.EventBus().SendInfo(fmt.Sprintf("Scaffolded code for %s at %s", igniteRepo.From.Original(), sFrom.Output))

			session.StartSpinner(fmt.Sprintf("Running scaffold commands for %s...", igniteRepo.To.Original()))
			sTo, err := scaffold.New(toBin, igniteRepo.To, scaffoldOptions...)
			if err != nil {
				return err
			}
			defer sTo.Cleanup()

			if err := sTo.Run(cmd.Context()); err != nil {
				return err
			}
			session.StopSpinner()
			session.EventBus().SendInfo(fmt.Sprintf("Scaffolded code for %s at %s", igniteRepo.To.Original(), sTo.Output))

			// Calculate and save the diffs from the scaffolded code.
			session.StartSpinner("Calculating diff...")
			diffs, err := diff.CalculateDiffs(sFrom.Output, sTo.Output)
			if err != nil {
				return errors.Wrap(err, "failed to calculate diff")
			}

			formatedDiffs, err := diff.FormatDiffs(diffs)
			if err != nil {
				return errors.Wrap(err, "failed to save diff map")
			}
			session.StopSpinner()
			session.EventBus().SendInfo("Diff calculated successfully")

			output, err = filepath.Abs(output)
			if err != nil {
				return errors.Wrap(err, "failed to find the abs path")
			}

			// Generate the docs file.
			g, err := doc.NewGenerator(doc.Options{
				Path:        output,
				FromVersion: igniteRepo.From,
				ToVersion:   igniteRepo.To,
				Diffs:       string(formatedDiffs),
				Description: releaseDescription,
			})
			if _, err := xgenny.RunWithValidation(placeholder.New(), g); err != nil {
				return err
			}

			session.Printf("Migration doc generated successfully at %s\n", output)

			return nil
		},
	}

	cmd.Flags().StringP(flagFrom, "f", "", "Version of ignite or path to ignite source code to generate the diff from")
	cmd.Flags().StringP(flagTo, "t", "", "Version of ignite or path to ignite source code to generate the diff to")
	cmd.Flags().StringP(flagOutput, "o", defaultDocPath, "Output directory to save the migration document")
	cmd.Flags().StringP(flagSource, "s", "", "Path to ignite source code repository. Set the source automatically set the cleanup to false")
	cmd.Flags().String(flagRepoURL, "", "Git URL for the Ignite repository")
	cmd.Flags().String(flagRepoOutput, "", "Output path to clone the ignite repository")
	cmd.Flags().Bool(flagRepoCleanup, true, "Cleanup the repository path after use")
	cmd.Flags().String(flagScaffoldOutput, "", "Output path to clone the ignite repository")
	cmd.Flags().String(flagScaffoldCache, "", "Output path to clone the ignite repository")

	return cmd
}
