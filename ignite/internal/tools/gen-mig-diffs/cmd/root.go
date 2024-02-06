package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	semver "github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/hexops/gotextdiff"
	"github.com/ignite/cli/v28/ignite/internal/tools/gen-mig-diffs/diff"
	"github.com/ignite/cli/v28/ignite/internal/tools/gen-mig-diffs/scaffold"
	"github.com/ignite/cli/v28/ignite/pkg/cmdrunner/exec"
	"github.com/ignite/cli/v28/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	fromFlag = "from"
	toFlag   = "to"
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
		RunE:  generateMigrationDiffs,
	}

	cmd.Flags().StringP(fromFlag,"f", "", "Version of ignite or path to ignite source code to generate the diff from")
	cmd.Flags().StringP(toFlag, "t", "", "Version of ignite or path to ignite source code to generate the diff to")
	cmd.Flags().StringP(outputFlag, "o", ".", "Output directory to save the migration diff files")

	return cmd
}

func generateMigrationDiffs(cmd *cobra.Command, args []string) error {
	from, _ := cmd.Flags().GetString(fromFlag)
	to, _ := cmd.Flags().GetString(toFlag)
	output, _ := cmd.Flags().GetString(outputFlag)

	// A temporary directory is created to clone ignite cli repository and build it
	tmpdir, err := os.MkdirTemp("", ".migdoc")
	defer os.RemoveAll(tmpdir)
	if err != nil {
		return err
	}
	log.Println("Created temporary directory:", tmpdir)

	var fromVer, toVer *semver.Version
	if ver, err := semver.NewVersion(from); err == nil {
		fromVer = ver
	}
	if ver, err := semver.NewVersion(to); err == nil {
		toVer = ver
	}

	repoDir, err := cloneIgniteRepo(tmpdir)
	if err != nil {
		return err
	}

	versions, err := getRepositoryVersionTags(repoDir)
	if err != nil {
		return err
	}

	fromVer, toVer, err = validateVersionRange(fromVer, toVer, versions)
	if err != nil {
		return err
	}

	log.Printf("Generating migration document for %s->%s\n", fromVer, toVer)

	// Run scaffolds for fromVer and toVer
	fromVerDir := filepath.Join(tmpdir, fromVer.Original())
	err = runScaffoldsForVersion(repoDir, fromVerDir, fromVer)
	if err != nil {
		return errors.Wrapf(err, "failed to run scaffolds for tag %s", fromVer)
	}
	toVerDir := filepath.Join(tmpdir, toVer.Original())
	err = runScaffoldsForVersion(repoDir, toVerDir, toVer)
	if err != nil {
		return errors.Wrapf(err, "failed to run scaffolds for tag %s", toVer)
	}

	// Run diff between two directories
	log.Println("Generating diff...")
	diffMap, err := calculateDiff(fromVerDir, toVerDir)
	if err != nil {
		return err
	}

	outputDir := filepath.Join(output, fmt.Sprintf("migdoc-%s-%s", fromVer, toVer))
	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return err
	}
	err = saveDiffMap(diffMap, outputDir)
	if err != nil {
		return err
	}
	log.Println("Migration document generated successfully at", outputDir)

	return nil
}

func isStrVersion(s string) bool {
	_, err := semver.NewVersion(s)
	return err == nil
}

func cloneIgniteRepo(tmpdir string) (string, error) {
	log.Println("Cloning", igniteCliRepository)

	repoDir := filepath.Join(tmpdir, igniteRepoPath)
	_, err := git.PlainClone(repoDir, false, &git.CloneOptions{
		URL: igniteCliRepository,
	})
	return repoDir, err
}

// getRepositoryVersionTags returns a sorted collection of semver tags from the ignite cli repository
func getRepositoryVersionTags(repoDir string) (semver.Collection, error) {
	repo, err := git.PlainOpen(repoDir)
	if err != nil {
		return nil, err
	}

	tags, err := repo.Tags()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tags")
	}

	// Iterate over all tags in the repository and pick valid semver tags
	var versions semver.Collection
	err = tags.ForEach(func(ref *plumbing.Reference) error {
		name := ref.Name()
		if name.IsTag() {
			ver, err := semver.NewVersion(name.Short())
			if err != nil {
				// Do nothing as it's not a semver tag
				return nil
			}
			versions = append(versions, ver)
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to iterate over tags")
	}

	sort.Sort(versions)

	return versions, nil
}

// validateVersionRange checks if the provided fromVer and toVer exist in the versions and if any of them is nil, then it picks default values.
func validateVersionRange(fromVer, toVer *semver.Version, versions semver.Collection) (*semver.Version, *semver.Version, error) {
	// Unable to generate migration document if there are less than two releases!
	if versions.Len() < 2 {
		return nil, nil, errors.New("At least two semver tags are required")
	}

	// Replace fromVer and toVer with equivalent semver tags from versions
	if fromVer != nil {
		found := false
		for _, ver := range versions {
			if ver.Equal(fromVer) {
				fromVer = ver
				found = true
				break
			}
		}
		if !found {
			return nil, nil, errors.Errorf("tag %s not found", fromVer)
		}
	}
	if toVer != nil {
		found := false
		for _, ver := range versions {
			if ver.Equal(toVer) {
				toVer = ver
				found = true
				break
			}
		}
		if !found {
			return nil, nil, errors.Errorf("tag %s not found", toVer)
		}
	}

	// Picking default values for fromVer and toVer such that:
	// If both fromVer and toVer are not provided, then generate migration document for second last and last semver tags
	// If only fromVer is not provided, then use the tag before toVer as fromVer
	// If only toVer is not provided, then use the last tag as toVer
	if fromVer == nil {
		if toVer != nil {
			sort.Search(versions.Len(), func(i int) bool {
				if versions[i].LessThan(toVer) {
					fromVer = versions[i]
					return false
				}
				return true
			})
		} else {
			fromVer = versions[versions.Len()-2]
		}
	}
	if toVer == nil {
		toVer = versions[versions.Len()-1]
	}

	// Unable to generate migration document if fromVer is greater or equal to toVer
	if fromVer.GreaterThan(toVer) || fromVer.Equal(toVer) {
		return nil, nil, errors.Errorf("from version %s should be less than to version %s", fromVer, toVer)
	}

	return fromVer, toVer, nil
}

// Run scaffolds commands one by one with the given version of ignite cli and save the output in the output directory
func runScaffoldsForVersion(repoDir, outputDir string, ver *semver.Version) error {
	repo, err := git.PlainOpen(repoDir)
	if err != nil {
		return errors.Wrap(err, "failed to open git repository")
	}

	wt, err := repo.Worktree()
	if err != nil {
		return errors.Wrap(err, "failed to get worktree")
	}

	err = checkoutAndBuildIgniteCli(wt, ver.Original(), repoDir)
	if err != nil {
		return err
	}

	binPath := filepath.Join(repoDir, igniteBinaryPath)
	scaffolder := scaffold.NewScaffolder(binPath, scaffold.DefaultScaffoldCommands)
	err = scaffolder.RunScaffolds(ver, outputDir)
	if err != nil {
		return err
	}

	return nil
}

func checkoutAndBuildIgniteCli(wt *git.Worktree, tag, repoDir string) error {
	err := wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewTagReferenceName(tag),
	})
	if err != nil {
		return errors.Wrapf(err, "failed to checkout tag %s", tag)
	}

	err = exec.Exec(context.Background(), []string{"make", "build"}, exec.StepOption(step.Workdir(repoDir)))
	if err != nil {
		return errors.Wrap(err, "failed to build ignite cli using make build")
	}

	return nil
}

func calculateDiff(fromVerDir, toVerDir string) (map[string][]gotextdiff.Unified, error) {
	diffMap := make(map[string]*diff.Diff)
	for _, s := range scaffold.DefaultScaffoldCommands {
		diff, err := diff.ComputeDiff(filepath.Join(fromVerDir, s.Name), filepath.Join(toVerDir, s.Name))
		if err != nil {
			return nil, err
		}
		diffMap[s.Name] = diff
	}

	subtractBaseDiffs(diffMap)

	unifiedDiffMap := make(map[string][]gotextdiff.Unified)
	for name, diff := range diffMap {
		unifiedDiffMap[name] = diff.ToUnified()
	}

	return unifiedDiffMap, nil
}

// subtractBaseDiffs removes chain and module diffs from other diffs
func subtractBaseDiffs(diffMap map[string]*diff.Diff) error {
	chainDiff := diffMap["chain"]
	moduleDiff := diffMap["module"]
	for name, diff := range diffMap {
		if name != "chain" && name != "module" {
			err := diff.Subtract(moduleDiff)
			if err != nil {
				return errors.Wrapf(err, "failed to subtract module diff from %s diff", name)
			}
		}
		diffMap[name] = diff
	}

	if err := diffMap["module"].Subtract(chainDiff); err != nil {
		return errors.Wrap(err, "failed to subtract chain diff from module diff")
	}

	return nil
}

func saveDiffMap(diffMap map[string][]gotextdiff.Unified, outputPath string) error {
	for name, diffs := range diffMap {
		outf, err := os.Create(filepath.Join(outputPath, name+".diff"))
		if err != nil {
			return err
		}
		defer outf.Close()
		for _, diff := range diffs {
			outf.WriteString(fmt.Sprint(diff))
			outf.WriteString("\n")
		}
	}

	return nil
}
