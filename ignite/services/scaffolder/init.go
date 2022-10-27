package scaffolder

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/gobuffalo/genny"
	vue "github.com/ignite/web"

	"github.com/ignite/cli/ignite/pkg/cache"
	"github.com/ignite/cli/ignite/pkg/cmdrunner/exec"
	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/ignite/pkg/gocmd"
	"github.com/ignite/cli/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/ignite/pkg/localfs"
	"github.com/ignite/cli/ignite/pkg/placeholder"
	"github.com/ignite/cli/ignite/templates/app"
	modulecreate "github.com/ignite/cli/ignite/templates/module/create"
)

var (
	commitMessage = "Initialized with Ignite CLI"
	devXAuthor    = &object.Signature{
		Name:  "Developer Experience team at Tendermint",
		Email: "hello@tendermint.com",
		When:  time.Now(),
	}
)

// Init initializes a new app with name and given options.
func Init(ctx context.Context, cacheStorage cache.Storage, tracer *placeholder.Tracer, root, name, addressPrefix string, noDefaultModule bool) (path string, err error) {
	if root, err = filepath.Abs(root); err != nil {
		return "", err
	}

	pathInfo, err := gomodulepath.Parse(name)
	if err != nil {
		return "", err
	}

	path = filepath.Join(root, pathInfo.Root)

	// create the project
	if err := generate(ctx, tracer, pathInfo, addressPrefix, path, noDefaultModule); err != nil {
		return "", err
	}

	if err := finish(ctx, cacheStorage, path, pathInfo.RawPath); err != nil {
		return "", err
	}

	// initialize git repository and perform the first commit
	if err := initGit(path); err != nil {
		return "", err
	}

	return path, nil
}

//nolint:interfacer
func generate(
	ctx context.Context,
	tracer *placeholder.Tracer,
	pathInfo gomodulepath.Path,
	addressPrefix,
	absRoot string,
	noDefaultModule bool,
) error {
	githubPath := gomodulepath.ExtractAppPath(pathInfo.RawPath)
	if !strings.Contains(githubPath, "/") {
		// A username must be added when the app module path has a single element
		githubPath = fmt.Sprintf("username/%s", githubPath)
	}

	g, err := app.New(&app.Options{
		// generate application template
		ModulePath:       pathInfo.RawPath,
		AppName:          pathInfo.Package,
		AppPath:          absRoot,
		GitHubPath:       githubPath,
		BinaryNamePrefix: pathInfo.Root,
		AddressPrefix:    addressPrefix,
	})
	if err != nil {
		return err
	}

	run := func(runner *genny.Runner, gen *genny.Generator) error {
		runner.With(gen)
		runner.Root = absRoot
		return runner.Run()
	}
	if err := run(genny.WetRunner(ctx), g); err != nil {
		return err
	}

	// generate module template
	if !noDefaultModule {
		opts := &modulecreate.CreateOptions{
			ModuleName: pathInfo.Package, // App name
			ModulePath: pathInfo.RawPath,
			AppName:    pathInfo.Package,
			AppPath:    absRoot,
			IsIBC:      false,
		}
		g, err = modulecreate.NewStargate(opts)
		if err != nil {
			return err
		}
		if err := run(genny.WetRunner(context.Background()), g); err != nil {
			return err
		}
		g = modulecreate.NewStargateAppModify(tracer, opts)
		if err := run(genny.WetRunner(context.Background()), g); err != nil {
			return err
		}

	}

	// FIXME(tb) untagged version of ignite/cli triggers a 404 not found when go
	// mod tidy requests the sumdb, until we understand why, we disable sumdb.
	// related issue:  https://github.com/golang/go/issues/56174
	opt := exec.StepOption(step.Env("GOSUMDB=off"))
	if err := gocmd.ModTidy(ctx, absRoot, opt); err != nil {
		return err
	}

	// generate the vue app.
	return Vue(filepath.Join(absRoot, "vue"))
}

// Vue scaffolds a Vue.js app for a chain.
func Vue(path string) error {
	return localfs.Save(vue.Boilerplate(), path)
}

func initGit(path string) error {
	repo, err := git.PlainInit(path, false)
	if err != nil {
		return err
	}
	wt, err := repo.Worktree()
	if err != nil {
		return err
	}
	if _, err := wt.Add("."); err != nil {
		return err
	}
	_, err = wt.Commit(commitMessage, &git.CommitOptions{
		All:    true,
		Author: devXAuthor,
	})
	return err
}
