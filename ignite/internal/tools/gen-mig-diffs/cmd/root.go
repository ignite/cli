package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/spf13/cobra"

	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"

	"github.com/ignite/cli/ignite/internal/tools/gen-mig-diffs/pkg/diff"
	"github.com/ignite/cli/ignite/internal/tools/gen-mig-diffs/pkg/repo"
	"github.com/ignite/cli/ignite/internal/tools/gen-mig-diffs/pkg/scaffold"
	"github.com/ignite/cli/ignite/internal/tools/gen-mig-diffs/pkg/url"
	"github.com/ignite/cli/ignite/internal/tools/gen-mig-diffs/templates/doc"
)

const (
	flagFrom           = "from"
	flagTo             = "to"
	flagOutput         = "output"
	flagSource         = "repo-source"
	flagRepoURL        = "repo-url"
	flagRepoOutput     = "repo-output"
	flagScaffoldOutput = "scaffold-output"
	flagScaffoldCache  = "scaffold-cache"
	flagYes            = "yes"

	defaultDocPath = "docs/docs/06-migration"
)

// NewRootCmd creates a new root command.
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen-mig-diffs",
		Short: "generate migration diffs from two different version",
		Long:  "This tool is used to generate migration diff files for each of ignites scaffold commands",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				from, _           = cmd.Flags().GetString(flagFrom)
				to, _             = cmd.Flags().GetString(flagTo)
				repoSource, _     = cmd.Flags().GetString(flagSource)
				output, _         = cmd.Flags().GetString(flagOutput)
				repoURLStr, _     = cmd.Flags().GetString(flagRepoURL)
				repoOutput, _     = cmd.Flags().GetString(flagRepoOutput)
				scaffoldOutput, _ = cmd.Flags().GetString(flagScaffoldOutput)
				scaffoldCache, _  = cmd.Flags().GetString(flagScaffoldCache)
				yes, _            = cmd.Flags().GetBool(flagYes)
			)
			session := cliui.New(cliui.WithoutUserInteraction(yes))
			defer session.End()

			fromVer, err := semver.NewVersion(from)
			if err != nil && from != "" {
				return errors.Wrapf(err, "failed to parse from version %s", from)
			}
			toVer, err := semver.NewVersion(to)
			if err != nil && to != "" {
				return errors.Wrapf(err, "failed to parse to version %s", to)
			}

			// Check or download the source and generate the binaries for each version.
			repoOptions := []repo.Options{repo.WithStdOutput(cmd.OutOrStdout())}
			if repoSource != "" {
				repoOptions = append(repoOptions, repo.WithSource(repoSource))
			}
			if repoURLStr != "" {
				repoURL, err := url.New(repoURLStr)
				if err != nil {
					return err
				}
				repoOptions = append(repoOptions, repo.WithRepoURL(repoURL))
			}
			if repoOutput != "" {
				repoOptions = append(repoOptions, repo.WithRepoOutput(repoOutput))
			}

			igniteRepo, err := repo.New(cmd.Context(), fromVer, toVer, session, repoOptions...)
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
			scaffoldOptions := []scaffold.Option{
				scaffold.WithStderr(os.Stderr),
				scaffold.WithStdout(os.Stdout),
				scaffold.WithStdin(os.Stdin),
			}
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
			if err != nil {
				return errors.Wrap(err, "failed to create the doc generator object")
			}

			runner := xgenny.NewRunner(cmd.Context(), output)
			sm, err := runner.RunAndApply(g, xgenny.ApplyPreRun(func(_, _, duplicated []string) error {
				if len(duplicated) == 0 {
					return nil
				}
				question := fmt.Sprintf("Do you want to overwrite the existing files? \n%s", strings.Join(duplicated, "\n"))
				return session.AskConfirm(question)
			}))
			if err != nil {
				return err
			}

			files := append(sm.CreatedFiles(), sm.ModifiedFiles()...)
			if len(files) == 0 {
				return errors.Errorf("migration doc not created at %s", output)
			}
			session.EventBus().SendInfo(
				fmt.Sprintf("Migration doc generated successfully at %s", files[0]),
			)

			return nil
		},
	}

	defaultRepoURL := repo.DefaultRepoURL.String()
	cmd.Flags().StringP(flagFrom, "f", "", "Version of Ignite or path to Ignite source code to generate the diff from")
	cmd.Flags().StringP(flagTo, "t", "", "Version of Ignite or path to Ignite source code to generate the diff to")
	cmd.Flags().StringP(flagOutput, "o", defaultDocPath, "Output directory to save the migration document")
	cmd.Flags().StringP(flagSource, "s", "", "Path to Ignite source code repository. Set the source automatically set the cleanup to false")
	cmd.Flags().String(flagRepoURL, defaultRepoURL, "Git URL for the Ignite repository")
	cmd.Flags().String(flagRepoOutput, "", "Output path to clone the Ignite repository")
	cmd.Flags().String(flagScaffoldOutput, "", "Output path to clone the Ignite repository")
	cmd.Flags().String(flagScaffoldCache, "", "Path to cache directory")
	cmd.Flags().BoolP(flagYes, "y", false, "answers interactive yes/no questions with yes")

	return cmd
}
