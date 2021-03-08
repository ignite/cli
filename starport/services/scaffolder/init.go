package scaffolder

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/gobuffalo/genny"
	conf "github.com/tendermint/starport/starport/chainconf"
	"github.com/tendermint/starport/starport/errors"
	"github.com/tendermint/starport/starport/pkg/cosmosanalysis/module"
	"github.com/tendermint/starport/starport/pkg/cosmosgen"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
	"github.com/tendermint/starport/starport/pkg/giturl"
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
	if err := s.protoc(absRoot, pathInfo.RawPath, s.options.sdkVersion); err != nil {
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

func (s *Scaffolder) protoc(projectPath, gomodPath string, version cosmosver.MajorVersion) error {
	if version != cosmosver.Stargate {
		return nil
	}

	if err := cosmosgen.InstallDependencies(context.Background(), projectPath); err != nil {
		if err == cosmosgen.ErrProtocNotInstalled {
			return errors.ErrStarportRequiresProtoc
		}
		return err
	}

	confpath, err := conf.LocateDefault(projectPath)
	if err != nil {
		return err
	}
	conf, err := conf.ParseFile(confpath)
	if err != nil {
		return err
	}

	options := []cosmosgen.Option{
		cosmosgen.WithGoGeneration(gomodPath),
		cosmosgen.IncludeDirs(conf.Build.Proto.ThirdPartyPaths),
	}

	// generate Vuex code as well if it is enabled.
	if conf.Client.Vuex.Path != "" {
		storeRootPath := filepath.Join(projectPath, conf.Client.Vuex.Path, "generated")
		options = append(options,
			cosmosgen.WithVuexGeneration(
				false,
				func(m module.Module) string {
					return filepath.Join(storeRootPath, giturl.UserAndRepo(m.Pkg.GoImportName), m.Pkg.Name, "module")
				},
				storeRootPath,
			),
		)
	}

	return cosmosgen.Generate(context.Background(), projectPath, conf.Build.Proto.Path, options...)
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
