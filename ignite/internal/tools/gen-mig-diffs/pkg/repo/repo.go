package repo

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/pkg/cmdrunner/exec"
	"github.com/ignite/cli/v28/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
)

const (
	defaultRepoURL    = "https://github.com/ignite/cli.git"
	defaultBinaryPath = "dist/ignite"
)

type (
	// Generator is used to generate migration diffs.
	Generator struct {
		From, To *semver.Version
		source   string
		repo     *git.Repository
		session  *cliui.Session
		cleanup  bool
	}

	// options represents configuration for the generator.
	options struct {
		source  string
		output  string
		repoURL string
		cleanup bool
	}
	// Options configures the generator.
	Options func(*options)
)

// newOptions returns a options with default options.
func newOptions() options {
	tmpDir := os.TempDir()
	return options{
		source:  "",
		output:  filepath.Join(tmpDir, "migration-source"),
		repoURL: defaultRepoURL,
	}
}

// WithSource set the repo source Options.
func WithSource(source string) Options {
	return func(o *options) {
		o.source = source
		// Do not clean up if set the source.
		o.cleanup = false
	}
}

// WithRepoURL set the repo URL Options.
func WithRepoURL(repoURL string) Options {
	return func(o *options) {
		o.repoURL = repoURL
	}
}

// WithRepoOutput set the repo output Options.
func WithRepoOutput(output string) Options {
	return func(o *options) {
		o.output = output
	}
}

// WithCleanup cleanup folders after use.
func WithCleanup() Options {
	return func(o *options) {
		o.cleanup = true
	}
}

// validate options
func (o options) validate() error {
	if o.source != "" && (o.repoURL != defaultRepoURL) {
		return errors.New("cannot set source and repo URL at the same time")
	}
	if o.source != "" && o.cleanup {
		return errors.New("cannot set source and cleanup the repo")
	}
	return nil
}

// New creates a new generator for migration diffs between from and to versions of ignite cli
// If source is empty, then it clones the ignite cli repository to a temporary directory and uses it as the source.
func New(from, to *semver.Version, session *cliui.Session, options ...Options) (*Generator, error) {
	opts := newOptions()
	for _, apply := range options {
		apply(&opts)
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}

	tempDir, err := os.MkdirTemp("", ".migdoc")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create temporary directory")
	}

	session.StopSpinner()
	session.EventBus().SendInfo(fmt.Sprintf("Created temporary directory: %s", tempDir))

	var (
		source = opts.source
		repo   *git.Repository
	)
	if source != "" {
		repo, err = git.PlainOpen(source)
		if err != nil {
			return nil, errors.Wrap(err, "failed to open ignite repository")
		}

		session.EventBus().SendInfo(fmt.Sprintf("Using ignite repository at: %s", source))
	} else {
		session.StartSpinner("Cloning ignite repository...")

		source = opts.output
		if err := os.RemoveAll(source); err != nil {
			return nil, errors.Wrap(err, "failed to clean the output directory")
		}

		repo, err = git.PlainClone(source, false, &git.CloneOptions{URL: opts.repoURL, Depth: 1})
		if err != nil {
			return nil, errors.Wrap(err, "failed to clone ignite repository")
		}

		session.StopSpinner()
		session.EventBus().SendInfo(fmt.Sprintf("Cloned ignite repository to: %s", source))
	}

	versions, err := getRepoVersionTags(source)
	if err != nil {
		return nil, err
	}

	from, to, err = validateVersionRange(from, to, versions)
	if err != nil {
		return nil, err
	}

	return &Generator{
		From:    from,
		To:      to,
		source:  source,
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

// Cleanup cleanup all temporary directories.
func (g *Generator) Cleanup() {
	if !g.cleanup {
		return
	}
	if err := os.RemoveAll(g.source); err != nil {
		g.session.EventBus().SendError(err)
		return
	}
	g.session.EventBus().SendInfo(fmt.Sprintf("Removed temporary directory: %s", g.source))
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

func (g *Generator) GenerateBinaries() (string, string, error) {
	fromBinPath, err := g.buildIgniteCli(g.From)
	if err != nil {
		return "", "", errors.Wrapf(err, "failed to run scaffolds for 'FROM' version %s", g.From)
	}
	toBinPath, err := g.buildIgniteCli(g.To)
	if err != nil {
		return "", "", errors.Wrapf(err, "failed to run scaffolds for 'TO' version %s", g.To)
	}
	return fromBinPath, toBinPath, nil
}

// buildIgniteCli build the ignite CLI from version.
func (g *Generator) buildIgniteCli(ver *semver.Version) (string, error) {
	g.session.StartSpinner(fmt.Sprintf("Building binary for version v%s...", ver))

	if err := g.checkoutToTag(ver.Original()); err != nil {
		return "", err
	}

	err := exec.Exec(context.Background(), []string{"make", "build"}, exec.StepOption(step.Workdir(g.source)))
	if err != nil {
		return "", errors.Wrap(err, "failed to build ignite cli using make build")
	}

	binPath := filepath.Join(g.source, defaultBinaryPath)

	g.session.StopSpinner()
	g.session.EventBus().SendInfo(fmt.Sprintf("Built ignite cli for v%s", ver))

	return binPath, nil
}

// checkoutToTag checkout the repository from a specific git tag.
func (g *Generator) checkoutToTag(tag string) error {
	wt, err := g.repo.Worktree()
	if err != nil {
		return err
	}

	// Reset and clean the git directory before the checkout to avoid conflicts.
	if err := wt.Reset(&git.ResetOptions{Mode: git.HardReset}); err != nil {
		return errors.Wrapf(err, "failed to reset %s", g.source)
	}
	if err := wt.Clean(&git.CleanOptions{Dir: true}); err != nil {
		return errors.Wrapf(err, "failed to reset %s", g.source)
	}
	if err = wt.Checkout(&git.CheckoutOptions{Branch: plumbing.NewTagReferenceName(tag)}); err != nil {
		return errors.Wrapf(err, "failed to checkout tag %s", tag)
	}
	return nil
}
