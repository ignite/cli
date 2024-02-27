package migdiff

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/hexops/gotextdiff"

	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/pkg/cmdrunner/exec"
	"github.com/ignite/cli/v28/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v28/ignite/pkg/diff"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
)

const (
	defaultRepoURL    = "http://github.com/ignite/cli.git"
	defaultRepoPath   = "src/github.com/ignite/cli"
	defaultBinaryPath = "dist/ignite"
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

type (
	// Generator is used to generate migration diffs.
	Generator struct {
		from, to         *semver.Version
		tempDir, repoDir string
		repo             *git.Repository
		session          *cliui.Session
	}

	// genOptions represents configuration for the generator.
	genOptions struct {
		source   string
		repoPath string
		repoURL  string
	}
	// GenOptions configures the generator.
	GenOptions func(*genOptions)
)

// newGenOptions returns a genOptions with default options.
func newGenOptions() genOptions {
	return genOptions{
		source:   "",
		repoPath: defaultRepoPath,
		repoURL:  defaultRepoURL,
	}
}

// WithSource set the repo source GenOptions.
func WithSource(source string) GenOptions {
	return func(m *genOptions) {
		m.source = source
	}
}

// WithRepoPath set the repo path GenOptions.
func WithRepoPath(repoPath string) GenOptions {
	return func(m *genOptions) {
		m.repoPath = repoPath
	}
}

// WithRepoURL set the repo URL GenOptions.
func WithRepoURL(repoURL string) GenOptions {
	return func(m *genOptions) {
		m.repoURL = repoURL
	}
}

// NewGenerator creates a new generator for migration diffs between from and to versions of ignite cli
// If source is empty, then it clones the ignite cli repository to a temporary directory and uses it as the source.
func NewGenerator(from, to *semver.Version, session *cliui.Session, options ...GenOptions) (*Generator, error) {
	opts := newGenOptions()
	for _, apply := range options {
		apply(&opts)
	}

	tempDir, err := os.MkdirTemp("", ".migdoc")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create temporary directory")
	}
	session.EventBus().SendInfo(fmt.Sprintf("Created temporary directory: %s", tempDir))

	var (
		repoDir = opts.source
		repo    *git.Repository
	)
	if repoDir != "" {
		repo, err = git.PlainOpen(repoDir)
		if err != nil {
			return nil, errors.Wrap(err, "failed to open ignite repository")
		}

		session.EventBus().SendInfo(fmt.Sprintf("Using ignite repository at: %s", repoDir))
	} else {
		session.StartSpinner("Cloning ignite repository...")

		repoDir = filepath.Join(tempDir, opts.repoPath)
		repo, err = git.PlainClone(repoDir, false, &git.CloneOptions{URL: opts.repoURL})
		if err != nil {
			return nil, errors.Wrap(err, "failed to clone ignite repository")
		}

		session.StopSpinner()
		session.EventBus().SendInfo(fmt.Sprintf("Cloned ignite repository to: %s", repoDir))
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
		session: session,
	}, nil
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

	versionMap := make(map[string]*semver.Version)
	for _, ver := range versions {
		versionMap[ver.String()] = ver
	}

	// Picking default values for fromVer and toVer such that:
	// If both fromVer and toVer are not provided, then generate migration document for second last and last semver tags
	// If only fromVer is not provided, then use the tag before toVer as fromVer
	// If only toVer is not provided, then use the last tag as toVer
	if toVer != nil {
		if _, found := versionMap[toVer.String()]; !found {
			return nil, nil, errors.Errorf("tag %s not found", toVer)
		}
	} else {
		toVer = versions[versions.Len()-1]
	}

	// Replace fromVer and toVer with equivalent semver tags from versions
	if fromVer != nil {
		if _, found := versionMap[fromVer.String()]; !found {
			return nil, nil, errors.Errorf("tag %s not found", fromVer)
		}
	} else {
		sort.Search(versions.Len(), func(i int) bool {
			if versions[i].LessThan(toVer) {
				fromVer = versions[i]
				return false
			}
			return true
		})
	}

	// Unable to generate migration document if fromVer is greater or equal to toVer
	if fromVer.GreaterThan(toVer) || fromVer.Equal(toVer) {
		return nil, nil, errors.Errorf("from version %s should be less than to version %s", fromVer, toVer)
	}

	return fromVer, toVer, nil
}

func (g *Generator) Cleanup() {
	err := os.RemoveAll(g.tempDir)
	if err != nil {
		g.session.EventBus().SendError(err)
		return
	}

	g.session.EventBus().SendInfo(fmt.Sprintf("Removed temporary directory: %s", g.tempDir))
}

func (g *Generator) Generate(outputPath string) error {
	g.session.Printf("Generating migration diffs for v%s -> v%s\n", g.from, g.to)

	fromDir := filepath.Join(g.tempDir, g.from.Original())
	if err := g.runScaffoldsForVersion(g.from, fromDir); err != nil {
		return errors.Wrapf(err, "failed to run scaffolds for 'FROM' version %s", g.from)
	}
	toDir := filepath.Join(g.tempDir, g.to.Original())
	if err := g.runScaffoldsForVersion(g.to, toDir); err != nil {
		return errors.Wrapf(err, "failed to run scaffolds for 'TO' version %s", g.to)
	}

	g.session.StartSpinner("Calculating diff...")
	diffs, err := calculateDiffs(fromDir, toDir)
	if err != nil {
		return errors.Wrap(err, "failed to calculate diff")
	}
	g.session.StopSpinner()
	g.session.EventBus().SendInfo("Diff calculated successfully")

	if err = saveDiffs(diffs, outputPath); err != nil {
		return errors.Wrap(err, "failed to save diff map")
	}
	g.session.Println("Migration diffs generated successfully at", outputPath)

	return nil
}

// Run scaffolds commands one by one with the given version of ignite cli and save the output in the output directory.
func (g *Generator) runScaffoldsForVersion(ver *semver.Version, outputDir string) error {
	g.session.StartSpinner(fmt.Sprintf("Building ignite cli for v%s...", ver))

	if err := g.checkoutToTag(ver.Original()); err != nil {
		return err
	}

	if err := g.buildIgniteCli(); err != nil {
		return err
	}

	g.session.StopSpinner()
	g.session.EventBus().SendInfo(fmt.Sprintf("Built ignite cli for v%s", ver))

	g.session.StartSpinner(fmt.Sprintf("Running scaffold commands for v%s...", ver))

	binPath := filepath.Join(g.repoDir, defaultBinaryPath)
	s := NewScaffolder(binPath, defaultScaffoldCommands)
	if err := s.Run(ver, outputDir); err != nil {
		return err
	}

	g.session.StopSpinner()
	g.session.EventBus().SendInfo(fmt.Sprintf("Scaffolded code for commands at %s", outputDir))

	return nil
}

func (g *Generator) checkoutToTag(tag string) error {
	wt, err := g.repo.Worktree()
	if err != nil {
		return err
	}

	// Reset and clean the git directory before the checkout to avoid conflicts.
	if err := wt.Reset(&git.ResetOptions{Mode: git.HardReset}); err != nil {
		return errors.Wrapf(err, "failed to reset %s", g.repoDir)
	}
	if err := wt.Clean(&git.CleanOptions{Dir: true}); err != nil {
		return errors.Wrapf(err, "failed to reset %s", g.repoDir)
	}

	if err = wt.Checkout(&git.CheckoutOptions{Branch: plumbing.NewTagReferenceName(tag)}); err != nil {
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
	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		return err
	}

	for name, diffs := range diffs {
		output, err := os.Create(filepath.Join(outputPath, name+".diff"))
		if err != nil {
			return err
		}
		for _, d := range diffs {
			output.WriteString(fmt.Sprint(d))
			output.WriteString("\n")
		}
		if err := output.Close(); err != nil {
			return err
		}
	}

	return nil
}
