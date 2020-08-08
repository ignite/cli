package starportserve

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/gookit/color"
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
	appBackendWatchPaths = []string{
		"app",
		"cmd",
		"x",
	}

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
	conf           starportconf.Config
	version        version
	verbose        bool
	serveCancel    context.CancelFunc
	serveRefresher chan struct{}
	stdout, stderr io.Writer
}

// Serve serves user apps.
func Serve(ctx context.Context, app App, conf starportconf.Config, verbose bool) error {
	s := &starportServe{
		app:            app,
		conf:           conf,
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
	if err := cmdrunner.
		New(opts...).
		Run(ctx, s.buildSteps()...); err != nil {
		return err
	}
	return cmdrunner.
		New(append(opts, cmdrunner.RunParallel())...).
		Run(ctx, s.serverSteps()...)
}

func (s *starportServe) buildSteps() (steps step.Steps) {
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
			return &CannotBuildAppError{Log: buildErr.String()}
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
				if !xexec.IsCommandAvailable("go") {
					return errors.New("go must be avaiable in your path")
				}
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

	cwd, _ := os.Getwd()

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
	for _, account := range s.conf.Accounts {
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
						Run(context.Background(), step.New(step.NewOptions().
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
			"--name", s.conf.Accounts[0].Name,
			"--keyring-backend", "test",
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
			step.Exec(fmt.Sprintf("%[1]vd", s.app.Name), "start"),
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
			step.Exec(fmt.Sprintf("%[1]vcli", s.app.Name), "rest-server"),
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
	return cmdrunner.
		New().
		Run(ctx, step.New(
			step.Exec("npm", "run", "dev"),
			step.Workdir(filepath.Join(s.app.Path, "frontend")),
		))
}

func (s *starportServe) runDevServer(ctx context.Context) error {
	conf := Config{
		EngineAddr:            "http://localhost:26657",
		AppBackendAddr:        "http://localhost:1317",
		AppFrontendAddr:       "http://localhost:8080",
		DevFrontendAssetsPath: "./ui/dist",
	} // TODO get vals from const
	sv := &http.Server{
		Addr:    ":12345",
		Handler: newDevHandler(s.app, conf),
	}
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		sv.Shutdown(shutdownCtx)
	}()
	err := sv.ListenAndServe()
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
		return version{}, err
	}
	v.tag = strings.TrimPrefix(ref.Name().Short(), "v")
	v.hash = ref.Hash().String()
	return v, nil
}

type CannotBuildAppError struct {
	Log string
}

func (e *CannotBuildAppError) Error() string {
	return fmt.Sprintf("cannot build app:\n\n\t%s", e.Log)
}
