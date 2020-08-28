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

type App struct {
	Name string
	Path string
}

type version struct {
	tag  string
	hash string
}

type starportServe struct {
	app            App
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
	var err error
	s.version, err = s.appVersion()
	if err != nil && err != git.ErrRepositoryNotExists {
		return err
	}
	if verbose {
		s.stdout = os.Stdout
		s.stderr = os.Stderr
	}
	if err := s.checkSystem(); err != nil {
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

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	conf, err := s.config()
	if err != nil {
		return &CannotBuildAppError{err}
	}

	if err := cmdrunner.
		New(opts...).
		Run(ctx, s.buildSteps(ctx, conf, cwd)...); err != nil {
		return err
	}
	return cmdrunner.
		New(append(opts, cmdrunner.RunParallel())...).
		Run(ctx, s.serverSteps()...)
}

func (s *starportServe) buildSteps(ctx context.Context, conf starportconf.Config, cwd string) (
	steps step.Steps) {
	ldflags := fmt.Sprintf(`'-X github.com/cosmos/cosmos-sdk/version.Name=NewApp 
	-X github.com/cosmos/cosmos-sdk/version.ServerName=%sd 
	-X github.com/cosmos/cosmos-sdk/version.ClientName=%scli 
	-X github.com/cosmos/cosmos-sdk/version.Version=%s 
	-X github.com/cosmos/cosmos-sdk/version.Commit=%s'`, s.app.Name, s.app.Name, s.version.tag, s.version.hash)
	var (
		// no-dash app name.
		ndapp    = strings.ReplaceAll(s.app.Name, "-", "")
		ndappd   = ndapp + "d"
		ndappcli = ndapp + "cli"

		appd   = s.app.Name + "d"
		appcli = s.app.Name + "cli"

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

	steps.Add(step.New(step.NewOptions().
		Add(
			step.Exec(
				"go",
				"install",
				"-mod", "readonly",
				"-ldflags", ldflags,
				filepath.Join(cwd, "cmd", appd),
			),
			step.PreExec(func() error {
				fmt.Fprintln(s.stdLog(logStarport).out, "🛠️  Building the app...")
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
				"install",
				"-mod", "readonly",
				"-ldflags", ldflags,
				filepath.Join(cwd, "cmd", appcli),
			),
			step.PostExec(captureBuildErr),
		).
		Add(s.stdSteps(logStarport)...).
		Add(step.Stderr(buildErr))...,
	))
	steps.Add(step.New(step.NewOptions().
		Add(
			step.Exec(
				appd,
				"init",
				"mynode",
				"--chain-id", ndapp,
			),
			step.PreExec(func() error {
				return xos.RemoveAllUnderHome(fmt.Sprintf(".%s", ndappd))
			}),
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
				path := filepath.Join(home, "."+appd, "config/genesis.json")
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
		).
		Add(s.stdSteps(logAppd)...)...,
	))
	steps.Add(step.New(step.NewOptions().
		Add(
			step.Exec(
				appcli,
				"config",
				"keyring-backend",
				"test",
			),
			step.PreExec(func() error {
				return xos.RemoveAllUnderHome(fmt.Sprintf(".%s", ndappcli))
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
				step.Exec(
					appcli,
					"keys",
					"add",
					account.Name,
					"--output", "json",
				),
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
			Add(step.Stderr(mnemonic))..., // TODO why mnemonic comes from stderr?
		))
		steps.Add(step.New(step.NewOptions().
			Add(
				step.Exec(
					appcli,
					"keys",
					"show",
					account.Name,
					"-a",
				),
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
								appd,
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
	steps.Add(step.New(step.NewOptions().
		Add(step.Exec(
			appcli,
			"config",
			"chain-id",
			ndapp,
		)).
		Add(s.stdSteps(logAppcli)...)...,
	))
	steps.Add(step.New(step.NewOptions().
		Add(step.Exec(
			appcli,
			"config",
			"output",
			"json",
		)).
		Add(s.stdSteps(logAppcli)...)...,
	))
	steps.Add(step.New(step.NewOptions().
		Add(step.Exec(
			appcli,
			"config",
			"indent",
			"true",
		)).
		Add(s.stdSteps(logAppcli)...)...,
	))
	steps.Add(step.New(step.NewOptions().
		Add(step.Exec(
			appcli,
			"config",
			"trust-node",
			"true",
		)).
		Add(s.stdSteps(logAppcli)...)...,
	))
	steps.Add(step.New(step.NewOptions().
		Add(step.Exec(
			appd,
			"gentx",
			"--name", conf.Validator.Name,
			"--keyring-backend", "test",
			"--amount", conf.Validator.Staked,
		)).
		Add(s.stdSteps(logAppd)...)...,
	))
	steps.Add(step.New(step.NewOptions().
		Add(step.Exec(
			appd,
			"collect-gentxs",
		)).
		Add(s.stdSteps(logAppd)...)...,
	))
	return
}

func (s *starportServe) serverSteps() (steps step.Steps) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		wg.Wait()
		fmt.Fprintf(s.stdLog(logStarport).out, "\n🚀 Get started: http://localhost:12345/\n\n")
	}()
	steps.Add(step.New(step.NewOptions().
		Add(
			step.Exec(
				fmt.Sprintf("%[1]vd", s.app.Name),
				"start",
			),
			step.InExec(func() error {
				defer wg.Done()
				fmt.Fprintf(s.stdLog(logStarport).out, "🌍 Running a Cosmos '%[1]v' app with Tendermint at http://localhost:26657.\n", s.app.Name)
				return nil
			}),
			step.PostExec(func(exitErr error) error {
				return errors.Wrapf(exitErr, "cannot run %[1]vd start", s.app.Name)
			}),
		).
		Add(s.stdSteps(logAppd)...)...,
	))
	steps.Add(step.New(step.NewOptions().
		Add(
			step.Exec(
				fmt.Sprintf("%[1]vcli", s.app.Name),
				"rest-server",
			),
			step.InExec(func() error {
				defer wg.Done()
				fmt.Fprintln(s.stdLog(logStarport).out, "🌍 Running a server at http://localhost:1317 (LCD)")
				return nil
			}),
			step.PostExec(func(exitErr error) error {
				return errors.Wrapf(exitErr, "cannot run %[1]vcli rest-server", s.app.Name)
			}),
		).
		Add(s.stdSteps(logAppcli)...)...,
	))
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
