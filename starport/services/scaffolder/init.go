package scaffolder

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/gobuffalo/genny"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/gomodulepath"
	"github.com/tendermint/starport/starport/pkg/xexec"
	"github.com/tendermint/starport/starport/templates/app"
)

var (
	commitMessage = "Initialized with Starport"
	devXAuthor    = &object.Signature{
		Name:  "Developer Experience team at Tendermint",
		Email: "hello@tendermint.com",
		When:  time.Now(),
	}
)

// InitOption configures scaffolding.
type InitOption func(*initOptions)

// initOptions keeps set of options to apply scaffolding.
type initOptions struct {
	addressPrefix string
}

// AddressPrefix configures address prefix for the app.
func AddressPrefix(prefix string) InitOption {
	return func(o *initOptions) {
		o.addressPrefix = prefix
	}
}

// Init initializes a new app with name and given options.
// path is the relative path to the scaffoled app.
func (s *Scaffolder) Init(name string, options ...InitOption) (path string, err error) {
	opts := &initOptions{}
	for _, o := range options {
		o(opts)
	}
	pathInfo, err := gomodulepath.Parse(name)
	if err != nil {
		return "", err
	}
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	absRoot := filepath.Join(pwd, pathInfo.Root)
	if err := s.generate(pathInfo, absRoot, opts); err != nil {
		return "", err
	}
	if err := s.protoc(absRoot); err != nil {
		return "", err
	}
	if err := initGit(pathInfo.Root); err != nil {
		return "", err
	}
	return pathInfo.Root, nil
}

func (s *Scaffolder) generate(pathInfo gomodulepath.Path, absRoot string, opts *initOptions) error {
	g, err := app.New(&app.Options{
		ModulePath:       pathInfo.RawPath,
		AppName:          pathInfo.Package,
		BinaryNamePrefix: pathInfo.Root,
		AddressPrefix:    opts.addressPrefix,
	})
	if err != nil {
		return err
	}
	run := genny.WetRunner(context.Background())
	run.With(g)
	run.Root = absRoot
	return run.Run()
}

// TODO warn if protoc isn't installed.
func (s *Scaffolder) protoc(absRoot string) error {
	scriptPath := filepath.Join(absRoot, "scripts/protocgen")
	if err := os.Chmod(scriptPath, 0700); err != nil {
		return err
	}
	return cmdrunner.
		New(
			cmdrunner.DefaultStderr(os.Stderr),
			cmdrunner.DefaultWorkdir(absRoot),
		).
		Run(context.Background(),
			// installs the gocosmos plugin with the version specified under the
			// go.mod of the app.
			step.New(
				step.Exec(
					"go",
					"get",
					"github.com/regen-network/cosmos-proto/protoc-gen-gocosmos",
				),
				step.PreExec(func() error {
					if !xexec.IsCommandAvailable("protoc") {
						return errors.New("Starport requires protoc installed.\nPlease, follow instructions on https://grpc.io/docs/protoc-installation")
					}
					return nil
				}),
			),
			// generate pb files.
			step.New(
				step.Exec(
					"/bin/bash",
					scriptPath,
				),
			),
		)
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
