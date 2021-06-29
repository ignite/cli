package scaffolder

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/gobuffalo/genny"
	"github.com/tendermint/starport/starport/pkg/giturl"
	"github.com/tendermint/starport/starport/pkg/gomodulepath"
	"github.com/tendermint/starport/starport/pkg/localfs"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/templates/app"
	modulecreate "github.com/tendermint/starport/starport/templates/module/create"
	"github.com/tendermint/vue"
)

var (
	commitMessage = "Initialized with Starport"
	devXAuthor    = &object.Signature{
		Name:  "Developer Experience team at Tendermint",
		Email: "hello@tendermint.com",
		When:  time.Now(),
	}
)

// Init initializes a new app with name and given options.
// path is the relative path to the scaffoled app.
func (s *Scaffolder) Init(tracer *placeholder.Tracer, name string, noDefaultModule bool) (path string, err error) {
	pathInfo, err := gomodulepath.Parse(name)
	if err != nil {
		return "", err
	}
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	absRoot := filepath.Join(pwd, pathInfo.Root)

	// create the project
	if err := s.generate(tracer, pathInfo, absRoot, noDefaultModule); err != nil {
		return "", err
	}

	if err := s.finish(absRoot, pathInfo.RawPath); err != nil {
		return "", err
	}

	// initialize git repository and perform the first commit
	if err := initGit(pathInfo.Root); err != nil {
		return "", err
	}
	return pathInfo.Root, nil
}

//nolint:interfacer
func (s *Scaffolder) generate(
	tracer *placeholder.Tracer,
	pathInfo gomodulepath.Path,
	absRoot string,
	noDefaultModule bool,
) error {
	gu, err := giturl.Parse(pathInfo.RawPath)
	if err != nil {
		return err
	}

	g, err := app.New(&app.Options{
		// generate application template
		ModulePath:       pathInfo.RawPath,
		AppName:          pathInfo.Package,
		OwnerName:        owner(pathInfo.RawPath),
		OwnerAndRepoName: gu.UserAndRepo(),
		BinaryNamePrefix: pathInfo.Root,
		AddressPrefix:    s.options.addressPrefix,
	})
	if err != nil {
		return err
	}

	run := func(runner *genny.Runner, gen *genny.Generator) error {
		runner.With(gen)
		runner.Root = absRoot
		return runner.Run()
	}
	if err := run(genny.WetRunner(context.Background()), g); err != nil {
		return err
	}

	// generate module template
	if !noDefaultModule {
		opts := &modulecreate.CreateOptions{
			ModuleName: pathInfo.Package, // App name
			ModulePath: pathInfo.RawPath,
			AppName:    pathInfo.Package,
			OwnerName:  owner(pathInfo.RawPath),
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

	// generate the vue app.
	vuepath := filepath.Join(absRoot, "vue")
	return localfs.Save(vue.Boilerplate(), vuepath)
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
