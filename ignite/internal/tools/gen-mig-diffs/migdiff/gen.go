package migdiff

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/hexops/gotextdiff"

	"github.com/ignite/cli/v28/ignite/pkg/cmdrunner/exec"
	"github.com/ignite/cli/v28/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v28/ignite/pkg/diff"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
)

const (
	igniteCliRepository = "http://github.com/ignite/cli.git"
	igniteRepoPath      = "src/github.com/ignite/cli"
	igniteBinaryPath    = "dist/ignite"
)

var diffIgnoreGlobs = []string{
	".git/**",
	"**.md",
	"go.sum",
	"**_test.go",
	"**.pb.go",
	"**.pb.gw.go",
	"**.pulsar.go",
	"**/node_modules/**",
	"**/openapi.yml",
	".gitignore",
	".github/**",
	"**.html",
	"**.css",
	"**.js",
	"**.ts",
	"**.json",
}

// Generator is used to generate migration diffs.
type Generator struct {
	from, to         *semver.Version
	tempDir, repoDir string
	repo             *git.Repository
	logger           *log.Logger
}

// NewGenerator creates a new generator for migration diffs between from and to versions of ignite cli
// If source is empty, then it clones the ignite cli repository to a temporary directory and uses it as the source.
func NewGenerator(from, to *semver.Version, source string) (*Generator, error) {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	tempDir, err := createTempDir()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create temporary directory")
	}
	logger.Println("Created temporary directory:", tempDir)

	var (
		repoDir = source
		repo    *git.Repository
	)
	if source == "" {
		logger.Println("Cloning ignite repository...")
		repoDir = filepath.Join(tempDir, igniteRepoPath)
		repo, err = cloneIgniteRepo(repoDir)
		if err != nil {
			return nil, errors.Wrap(err, "failed to clone ignite repository")
		}
	} else {
		repo, err = git.PlainOpen(source)
		if err != nil {
			return nil, errors.Wrap(err, "failed to open ignite repository")
		}
	}

	versions, err := getRepoVersionTags(repoDir)
	if err != nil {
		return nil, err
	}

	from, to, err = validateVersionRange(from, to, versions)
	if err != nil {
		return nil, err
	}

	return &Generator{
		from:    from,
		to:      to,
		tempDir: tempDir,
		repoDir: repoDir,
		repo:    repo,
		logger:  logger,
	}, nil
}

func createTempDir() (string, error) {
	tmpdir, err := os.MkdirTemp("", ".migdoc")
	defer os.RemoveAll(tmpdir)
	if err != nil {
		return "", err
	}

	return tmpdir, nil
}

func cloneIgniteRepo(path string) (*git.Repository, error) {
	repo, err := git.PlainClone(path, false, &git.CloneOptions{
		URL: igniteCliRepository,
	})
	return repo, err
}

// getRepoVersionTags returns a sorted collection of semver tags from the ignite cli repository.
func getRepoVersionTags(repoDir string) (semver.Collection, error) {
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

func (g *Generator) Cleanup() error {
	g.logger.Println("Cleaning up temporary directory:", g.tempDir)
	return os.RemoveAll(g.tempDir)
}

func (g *Generator) Generate(outputPath string) error {
	g.logger.Printf("Generating migration diffs for %s->%s\n", g.from, g.to)

	fromDir := filepath.Join(g.tempDir, g.from.Original())
	err := g.runScaffoldsForVersion(g.from, fromDir)
	if err != nil {
		return errors.Wrapf(err, "failed to run scaffolds for version %s", g.from)
	}
	toDir := filepath.Join(g.tempDir, g.to.Original())
	err = g.runScaffoldsForVersion(g.to, toDir)
	if err != nil {
		return errors.Wrapf(err, "failed to run scaffolds for version %s", g.to)
	}

	g.logger.Println("Calculating diff...")
	diffs, err := calculateDiffs(fromDir, toDir)
	if err != nil {
		return errors.Wrap(err, "failed to calculate diff")
	}

	err = saveDiffs(diffs, outputPath)
	if err != nil {
		return errors.Wrap(err, "failed to save diff map")
	}
	log.Println("Migration diffs generated successfully at", outputPath)

	return nil
}

// Run scaffolds commands one by one with the given version of ignite cli and save the output in the output directory.
func (g *Generator) runScaffoldsForVersion(ver *semver.Version, outputDir string) error {
	err := g.checkoutToTag(ver.Original())
	if err != nil {
		return err
	}

	err = g.buildIgniteCli()
	if err != nil {
		return err
	}

	binPath := filepath.Join(g.repoDir, igniteBinaryPath)
	scaffolder := NewScaffolder(binPath, defaultScaffoldCommands)
	err = scaffolder.Run(ver, outputDir)
	if err != nil {
		return err
	}

	return nil
}

func (g *Generator) checkoutToTag(tag string) error {
	wt, err := g.repo.Worktree()
	if err != nil {
		return err
	}

	err = wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewTagReferenceName(tag),
	})
	if err != nil {
		return errors.Wrapf(err, "failed to checkout tag %s", tag)
	}

	return nil
}

func (g *Generator) buildIgniteCli() error {
	err := exec.Exec(context.Background(), []string{"make", "build"}, exec.StepOption(step.Workdir(g.repoDir)))
	if err != nil {
		return errors.Wrap(err, "failed to build ignite cli using make build")
	}

	return nil
}

func calculateDiffs(fromDir, toDir string) (map[string][]gotextdiff.Unified, error) {
	diffs := make(map[string][]gotextdiff.Unified)
	for _, s := range defaultScaffoldCommands {
		diff, err := diff.ComputeFS(
			os.DirFS(filepath.Join(fromDir, s.Name)),
			os.DirFS(filepath.Join(toDir, s.Name)),
			diffIgnoreGlobs...,
		)
		if err != nil {
			return nil, err
		}
		diffs[s.Name] = diff
	}

	subtractBaseDiffs(diffs)

	return diffs, nil
}

// subtractBaseDiffs removes chain and module diffs from other diffs.
func subtractBaseDiffs(diffs map[string][]gotextdiff.Unified) {
	chainDiff := diffs["chain"]
	moduleDiff := diffs["module"]
	for name, d := range diffs {
		if name != "chain" && name != "module" {
			diffs[name] = subtractUnifieds(d, moduleDiff)
		}
	}

	diffs["module"] = subtractUnifieds(moduleDiff, chainDiff)
}

func subtractUnifieds(a, b []gotextdiff.Unified) []gotextdiff.Unified {
	for i, ad := range a {
		for _, bd := range b {
			if ad.From == bd.From && ad.To == bd.To {
				a[i] = diff.Subtract(ad, bd)
			}
		}
	}
	return a
}

func saveDiffs(diffs map[string][]gotextdiff.Unified, outputPath string) error {
	err := os.MkdirAll(outputPath, os.ModePerm)
	if err != nil {
		return err
	}

	for name, diffs := range diffs {
		outf, err := os.Create(filepath.Join(outputPath, name+".diff"))
		if err != nil {
			return err
		}
		defer outf.Close()
		for _, d := range diffs {
			outf.WriteString(fmt.Sprint(d))
			outf.WriteString("\n")
		}
	}

	return nil
}
