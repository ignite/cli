package starportserve

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go/build"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/gookit/color"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/fswatcher"
	"github.com/tendermint/starport/starport/pkg/xexec"
	"github.com/tendermint/starport/starport/pkg/xos"
	starportconf "github.com/tendermint/starport/starport/services/serve/conf"
	"golang.org/x/sync/errgroup"
)

var (
	appBackendWatchPaths = append([]string{
		"app",
		"cmd",
		"x",
	}, starportconf.FileNames...)

	vuePath = "vue"

	errorColor = color.Red.Render
	infoColor  = color.Yellow.Render
)

type version struct {
	tag  string
	hash string
}

type starportServe struct {
	app            App
	plugin         Plugin
	version        version
	verbose        bool
	serveCancel    context.CancelFunc
	serveRefresher chan struct{}
	stdout, stderr io.Writer
}

// Serve serves user apps.
func Serve(ctx context.Context, app App, verbose bool) error {
	s := &starportServe{
		app:            app,
		verbose:        verbose,
		serveRefresher: make(chan struct{}, 1),
		stdout:         ioutil.Discard,
		stderr:         ioutil.Discard,
	}

	if verbose {
		s.stdout = os.Stdout
		s.stderr = os.Stderr
	}

	var err error

	s.version, err = s.appVersion()
	if err != nil && err != git.ErrRepositoryNotExists {
		return err
	}

	s.plugin, err = s.pickPlugin()
	if err != nil {
		return err
	}
	fmt.Fprintf(s.stdLog(logStarport).out, "Cosmos' version is: %s\n", infoColor(s.plugin.Name()))

	if err := s.checkSystem(); err != nil {
		return err
	}
	if err := s.plugin.Migrate(ctx); err != nil {
		return err
	}

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return s.watchAppFrontend(ctx)
	})
	g.Go(func() error {
		return s.runDevServer(ctx)
	})
	g.Go(func() error {
		s.refreshServe()
		for {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			select {
			case <-ctx.Done():
				return ctx.Err()

			case <-s.serveRefresher:
				var (
					serveCtx context.Context
					buildErr *CannotBuildAppError
				)
				serveCtx, s.serveCancel = context.WithCancel(ctx)
				err := s.serve(serveCtx)
				switch {
				case err == nil:
				case errors.Is(err, context.Canceled):
				case errors.As(err, &buildErr):
					fmt.Fprintf(s.stdLog(logStarport).err, "%s\n", errorColor(err.Error()))

					var validationErr *starportconf.ValidationError
					if errors.As(err, &validationErr) {
						fmt.Fprintln(s.stdLog(logStarport).out, "see: https://github.com/tendermint/starport#configure")
					}

					fmt.Fprintf(s.stdLog(logStarport).out, "%s\n", infoColor("waiting for a fix before retrying..."))
				default:
					return err
				}
			}
		}
	})
	g.Go(func() error {
		return s.watchAppBackend(ctx)
	})
	return g.Wait()
}

// checkSystem checks if developer's work environment comply must to have
// dependencies and pre-conditions.
func (s *starportServe) checkSystem() error {
	// check if Go has installed.
	if !xexec.IsCommandAvailable("go") {
		return errors.New("Please, check that Go language is installed correctly in $PATH. See https://golang.org/doc/install")
	}
	// check if Go's bin added to System's path.
	gobinpath := path.Join(build.Default.GOPATH, "bin")
	if err := xos.IsInPath(gobinpath); err != nil {
		return errors.New("$(go env GOPATH)/bin must be added to your $PATH. See https://golang.org/doc/gopath_code.html#GOPATH")
	}
	return nil
}

func (s *starportServe) refreshServe() {
	if s.serveCancel != nil {
		s.serveCancel()
	}
	s.serveRefresher <- struct{}{}
}

func (s *starportServe) watchAppBackend(ctx context.Context) error {
	return fswatcher.Watch(
		ctx,
		appBackendWatchPaths,
		fswatcher.Workdir(s.app.Path),
		fswatcher.OnChange(s.refreshServe),
		fswatcher.IgnoreHidden(),
	)
}

func (s *starportServe) serve(ctx context.Context) error {
	opts := []cmdrunner.Option{
		cmdrunner.DefaultWorkdir(s.app.Path),
	}

	conf, err := s.config()
	if err != nil {
		return &CannotBuildAppError{err}
	}

	if err := cmdrunner.
		New(opts...).
		Run(ctx, s.buildSteps(ctx, conf)...); err != nil {
		return err
	}
	return cmdrunner.
		New(append(opts, cmdrunner.RunParallel())...).
		Run(ctx, s.serverSteps()...)
}

func (s *starportServe) buildSteps(ctx context.Context, conf starportconf.Config) (
	steps step.Steps) {
	ldflags := fmt.Sprintf(`'-X github.com/cosmos/cosmos-sdk/version.Name=NewApp 
	-X github.com/cosmos/cosmos-sdk/version.ServerName=%sd 
	-X github.com/cosmos/cosmos-sdk/version.ClientName=%scli 
	-X github.com/cosmos/cosmos-sdk/version.Version=%s 
	-X github.com/cosmos/cosmos-sdk/version.Commit=%s'`, s.app.Name, s.app.Name, s.version.tag, s.version.hash)
	var (
		buildErr = &bytes.Buffer{}
	)
	captureBuildErr := func(err error) error {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return &CannotBuildAppError{errors.New(buildErr.String())}
		}
		return err
	}
	steps.Add(step.New(step.NewOptions().
		Add(
			step.Exec(
				"go",
				"mod",
				"tidy",
			),
			step.PreExec(func() error {
				fmt.Fprintln(s.stdLog(logStarport).out, "\n📦 Installing dependencies...")
				return nil
			}),
			step.PostExec(captureBuildErr),
		).
		Add(s.stdSteps(logStarport)...).
		Add(step.Stderr(buildErr))...,
	))
	steps.Add(step.New(step.NewOptions().
		Add(
			step.Exec(
				"go",
				"mod",
				"verify",
			),
			step.PostExec(captureBuildErr),
		).
		Add(s.stdSteps(logBuild)...).
		Add(step.Stderr(buildErr))...,
	))

	// install the app.
	steps.Add(step.New(
		step.PreExec(func() error {
			fmt.Fprintln(s.stdLog(logStarport).out, "🛠️  Building the app...")
			return nil
		}),
	))
	for _, execOption := range s.plugin.InstallCommands(ldflags) {
		steps.Add(step.New(step.NewOptions().
			Add(
				execOption,
				step.PostExec(captureBuildErr),
			).
			Add(s.stdSteps(logStarport)...).
			Add(step.Stderr(buildErr))...,
		))
	}

	// cleanup persistent data from previous `serve`.
	steps.Add(step.New(
		step.PreExec(func() error {
			for _, path := range s.plugin.StoragePaths() {
				if err := xos.RemoveAllUnderHome(path); err != nil {
					return err
				}
			}
			return nil
		}),
	))

	// init node.
	steps.Add(step.New(step.NewOptions().
		Add(
			step.Exec(
				s.app.d(),
				"init",
				"mynode",
				"--chain-id", s.app.nd(),
			),
			step.PostExec(func(err error) error {
				// overwrite Genesis with user configs.
				if err != nil {
					return err
				}
				if conf.Genesis == nil {
					return nil
				}
				home, err := os.UserHomeDir()
				if err != nil {
					return err
				}
				path := filepath.Join(home, s.plugin.GenesisPath())
				file, err := os.OpenFile(path, os.O_RDWR, 644)
				if err != nil {
					return err
				}
				defer file.Close()
				var genesis map[string]interface{}
				if err := json.NewDecoder(file).Decode(&genesis); err != nil {
					return err
				}
				if err := mergo.Merge(&genesis, conf.Genesis, mergo.WithOverride); err != nil {
					return err
				}
				if err := file.Truncate(0); err != nil {
					return err
				}
				if _, err := file.Seek(0, 0); err != nil {
					return err
				}
				return json.NewEncoder(file).Encode(&genesis)
			}),
			step.PostExec(func(err error) error {
				if err != nil {
					return err
				}
				return s.plugin.PostInit()
			}),
		).
		Add(s.stdSteps(logAppd)...)...,
	))

	for _, account := range conf.Accounts {
		account := account
		var (
			key      = &bytes.Buffer{}
			mnemonic = &bytes.Buffer{}
		)
		steps.Add(step.New(step.NewOptions().
			Add(
				s.plugin.AddUserCommand(account.Name),
				step.PostExec(func(exitErr error) error {
					if exitErr != nil {
						return errors.Wrapf(exitErr, "cannot create %s account", account.Name)
					}
					var user struct {
						Mnemonic string `json:"mnemonic"`
					}
					if err := json.NewDecoder(mnemonic).Decode(&user); err != nil {
						return errors.Wrap(err, "cannot decode mnemonic")
					}
					fmt.Fprintf(s.stdLog(logStarport).out, "🙂 Created an account. Password (mnemonic): %[1]v\n", user.Mnemonic)
					return nil
				}),
			).
			Add(s.stdSteps(logAppcli)...).
			// Stargate pipes from stdout, Launchpad pipes from stderr.
			Add(step.Stderr(mnemonic), step.Stdout(mnemonic))...,
		))
		steps.Add(step.New(step.NewOptions().
			Add(
				s.plugin.ShowAccountCommand(account.Name),
				step.PostExec(func(err error) error {
					if err != nil {
						return err
					}
					coins := strings.Join(account.Coins, ",")
					key := strings.TrimSpace(key.String())
					return cmdrunner.
						New().
						Run(ctx, step.New(step.NewOptions().
							Add(step.Exec(
								s.app.d(),
								"add-genesis-account",
								key,
								coins,
							)).
							Add(s.stdSteps(logAppd)...)...,
						))
				}),
			).
			Add(s.stdSteps(logAppcli)...).
			Add(step.Stdout(key))...,
		))
	}

	for _, execOption := range s.plugin.ConfigCommands() {
		steps.Add(step.New(step.NewOptions().
			Add(execOption).
			Add(s.stdSteps(logAppcli)...)...,
		))
	}

	steps.Add(step.New(step.NewOptions().
		Add(s.plugin.GentxCommand(conf)).
		Add(s.stdSteps(logAppd)...)...,
	))
	steps.Add(step.New(step.NewOptions().
		Add(step.Exec(
			s.app.d(),
			"collect-gentxs",
		)).
		Add(s.stdSteps(logAppd)...)...,
	))
	return
}

func (s *starportServe) serverSteps() (steps step.Steps) {
	var wg sync.WaitGroup
	wg.Add(len(s.plugin.StartCommands()))
	go func() {
		wg.Wait()
		fmt.Fprintf(s.stdLog(logStarport).out, "🌍 Running a Cosmos '%[1]v' app with Tendermint at http://localhost:26657.\n", s.app.Name)
		fmt.Fprintln(s.stdLog(logStarport).out, "🌍 Running a server at http://localhost:1317 (LCD)")
		fmt.Fprintf(s.stdLog(logStarport).out, "\n🚀 Get started: http://localhost:12345/\n\n")
	}()
	for _, execOption := range s.plugin.StartCommands() {
		steps.Add(step.New(step.NewOptions().
			Add(execOption...).
			Add(step.InExec(func() error {
				wg.Done()
				return nil
			})).
			Add(s.stdSteps(logAppd)...)...,
		))
	}
	return
}

func (s *starportServe) watchAppFrontend(ctx context.Context) error {
	vueFullPath := filepath.Join(s.app.Path, vuePath)
	if _, err := os.Stat(vueFullPath); os.IsNotExist(err) {
		return nil
	}
	frontendErr := &bytes.Buffer{}
	postExec := func(err error) error {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() > 0 {
			fmt.Fprintf(s.stdLog(logStarport).err, "%s\n%s",
				infoColor("skipping serving Vue frontend due to following errors:"), errorColor(frontendErr.String()))
		}
		return nil // ignore errors.
	}
	return cmdrunner.
		New(
			cmdrunner.DefaultWorkdir(vueFullPath),
			cmdrunner.DefaultStderr(frontendErr),
		).
		Run(ctx,
			step.New(
				step.Exec("npm", "i"),
				step.PostExec(postExec),
			),
			step.New(
				step.Exec("npm", "run", "serve"),
				step.PostExec(postExec),
			),
		)
}

func (s *starportServe) runDevServer(ctx context.Context) error {
	conf := Config{
		SdkVersion:      s.plugin.Name(),
		EngineAddr:      "http://localhost:26657",
		AppBackendAddr:  "http://localhost:1317",
		AppFrontendAddr: "http://localhost:8080",
	} // TODO get vals from const
	handler, err := newDevHandler(s.app, conf)
	if err != nil {
		return err
	}
	sv := &http.Server{
		Addr:    ":12345",
		Handler: handler,
	}
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		sv.Shutdown(shutdownCtx)
	}()
	err = sv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func (s *starportServe) appVersion() (v version, err error) {
	repo, err := git.PlainOpen(s.app.Path)
	if err != nil {
		return version{}, err
	}
	iter, err := repo.Tags()
	if err != nil {
		return version{}, err
	}
	ref, err := iter.Next()
	if err != nil {
		return version{}, nil
	}
	v.tag = strings.TrimPrefix(ref.Name().Short(), "v")
	v.hash = ref.Hash().String()
	return v, nil
}

func (s *starportServe) config() (starportconf.Config, error) {
	var paths []string
	for _, name := range starportconf.FileNames {
		paths = append(paths, filepath.Join(s.app.Path, name))
	}
	confFile, err := xos.OpenFirst(paths...)
	if err != nil {
		return starportconf.Config{}, errors.Wrap(err, "config file cannot be found")
	}
	defer confFile.Close()
	return starportconf.Parse(confFile)
}

type CannotBuildAppError struct {
	Err error
}

func (e *CannotBuildAppError) Error() string {
	return fmt.Sprintf("cannot build app:\n\n\t%s", e.Err)
}

func (e *CannotBuildAppError) Unwrap() error {
	return e.Err
}
