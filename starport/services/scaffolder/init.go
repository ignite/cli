package scaffolder

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/gobuffalo/genny"
	"github.com/tendermint/starport/starport/errors"
	"github.com/tendermint/starport/starport/pkg/cosmosprotoc"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
	"github.com/tendermint/starport/starport/pkg/gomodulepath"
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

// Init initializes a new app with name and given options.
// path is the relative path to the scaffoled app.
func (s *Scaffolder) Init(name string) (path string, err error) {
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
	if err := s.generate(pathInfo, absRoot); err != nil {
		return "", err
	}

	// generate protobuf types
	if err := s.protoc(absRoot, s.options.sdkVersion); err != nil {
		return "", err
	}

	// format the source
	if err := fmtProject(absRoot); err != nil {
		return "", err
	}

	// initialize git repository and perform the first commit
	if err := initGit(pathInfo.Root); err != nil {
		return "", err
	}
	return pathInfo.Root, nil
}

func (s *Scaffolder) generate(pathInfo gomodulepath.Path, absRoot string) error {
	g, err := app.New(s.options.sdkVersion, &app.Options{
		ModulePath:       pathInfo.RawPath,
		AppName:          pathInfo.Package,
		OwnerName:        owner(pathInfo.RawPath),
		BinaryNamePrefix: pathInfo.Root,
		AddressPrefix:    s.options.addressPrefix,
	})
	if err != nil {
		return err
	}
	run := genny.WetRunner(context.Background())
	run.With(g)
	run.Root = absRoot
	return run.Run()
}

func (s *Scaffolder) protoc(absRoot string, version cosmosver.MajorVersion) error {
	if version != cosmosver.Stargate {
		return nil
	}
	if err := cosmosprotoc.InstallDependencies(context.Background(), absRoot); err != nil {
		if err == cosmosprotoc.ErrProtocNotInstalled {
			return errors.ErrStarportRequiresProtoc
		}
		return err
	}
	return cosmosprotoc.Generate(context.Background(),
		filepath.Join(absRoot, "proto"),
		filepath.Join(absRoot, "third_party/proto"),
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
