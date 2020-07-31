package starportserve

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/fswatcher"
	"github.com/tendermint/starport/starport/pkg/xexec"
	"golang.org/x/sync/errgroup"
)

var (
	appBackendWatchPaths = []string{
		"app",
		"cmd",
		"x",
	}
)

type App struct {
	Name string
	Path string
}

type starportServe struct {
	app     App
	verbose bool
}

// Serve serves user apps.
func Serve(ctx context.Context, app App, verbose bool) error {
	s := &starportServe{
		app:     app,
		verbose: verbose,
	}

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return s.watchAppFrontend(ctx)
	})
	g.Go(func() error {
		return s.runDevServer(ctx)
	})

	var (
		serveCtx    context.Context
		serveCancel context.CancelFunc
		serveErr    = make(chan error, 1)
	)
	serve := func() {
		if serveCancel != nil {
			serveCancel()
		}
		serveCtx, serveCancel = context.WithCancel(ctx)
		if err := s.serve(serveCtx); err != nil && err != context.Canceled {
			serveErr <- err
		}
	}
	go serve()

	g.Go(func() error {
		select {
		case err := <-serveErr:
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	})
	g.Go(func() error {
		return fswatcher.Watch(
			ctx,
			appBackendWatchPaths,
			fswatcher.Workdir(app.Path),
			fswatcher.OnChange(serve),
			fswatcher.IgnoreHidden(),
		)
	})
	return g.Wait()
}

func (s *starportServe) serve(ctx context.Context) error {
	var (
		stdout = ioutil.Discard
		stderr = ioutil.Discard
	)
	if s.verbose {
		stdout = os.Stdout
		stderr = os.Stderr
	}
	opts := []cmdrunner.Option{
		cmdrunner.DefaultStdout(stdout),
		cmdrunner.DefaultStderr(stderr),
		cmdrunner.DefaultWorkdir(s.app.Path),
	}

	if err := cmdrunner.
		New(opts...).
		Run(ctx, s.buildSteps()...); err != nil {
		return err
	}

	if err := cmdrunner.
		New(append(opts, cmdrunner.RunParallel())...).
		Run(ctx, s.serverSteps()...); err != nil {
		if _, ok := errors.Cause(err).(*exec.ExitError); ok {
			return nil
		}
		return err
	}
	return nil
}

func (s *starportServe) buildSteps() (steps step.Steps) {
	mnemonic := &bytes.Buffer{}
	steps.Add(step.New(
		step.Exec("go", "mod", "tidy"),
		step.PreExec(func() error {
			if !xexec.IsCommandAvailable("go") {
				return errors.New("go must be avaiable in your path")
			}
			fmt.Println("\n📦 Installing dependencies...")
			return nil
		}),
		step.PostExec(func(exitErr error) error {
			return errors.Wrap(exitErr, "cannot install go modules")
		}),
	))
	steps.Add(step.New(
		step.Exec("make"),
		step.PreExec(func() error {
			if !xexec.IsCommandAvailable("make") {
				return errors.New("make must be avaiable in your path")
			}
			fmt.Println("🚧 Building the application...")
			return nil
		}),
		step.PostExec(func(exitErr error) error {
			return errors.Wrap(exitErr, "cannot build your app")
		}),
	))
	steps.Add(step.New(
		step.Exec("make", "init-pre"),
		step.PreExec(func() error {
			fmt.Println("💫 Initializing the chain...")
			return nil
		}),
		step.PostExec(func(exitErr error) error {
			return errors.Wrap(exitErr, "cannot initialize the chain")
		}),
	))
	for _, user := range []string{"user1", "user2"} {
		steps.Add(step.New(
			step.Exec("make", fmt.Sprintf("init-%s", user), "-s"),
			step.PostExec(func(exitErr error) error {
				if exitErr != nil {
					return errors.Wrapf(exitErr, "cannot create %s account", user)
				}
				var user struct {
					Mnemonic string `json:"mnemonic"`
				}
				if err := json.Unmarshal(mnemonic.Bytes(), &user); err != nil {
					return errors.Wrap(err, "cannot decode mnemonic")
				}
				mnemonic.Reset()
				fmt.Printf("🙂 Created an account. Password (mnemonic): %[1]v\n", user.Mnemonic)
				return nil
			}),
			step.Stdout(mnemonic),
		))
	}
	steps.Add(step.New(
		step.Exec("make", "init-post"),
	))
	return
}

func (s *starportServe) serverSteps() (steps step.Steps) {
	steps.Add(step.New(
		step.Exec(fmt.Sprintf("%[1]vd", s.app.Name), "start"),
		step.InExec(func() error {
			if s.verbose {
				fmt.Println("🌍 Running a server at http://localhost:26657 (Tendermint)")
			} else {
				fmt.Printf("🌍 Running a Cosmos '%[1]v' app with Tendermint.\n", s.app.Name)
			}
			return nil
		}),
		step.PostExec(func(exitErr error) error {
			return errors.Wrapf(exitErr, "cannot run %[1]vd start", s.app.Name)
		}),
	))
	steps.Add(step.New(
		step.Exec(fmt.Sprintf("%[1]vcli", s.app.Name), "rest-server"),
		step.InExec(func() error {
			if s.verbose {
				fmt.Println("🌍 Running a server at http://localhost:1317 (LCD)")
			}
			return nil
		}),
		step.PostExec(func(exitErr error) error {
			return errors.Wrapf(exitErr, "cannot run %[1]vcli rest-server", s.app.Name)
		}),
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
	if s.verbose {
		fmt.Printf("🔧 Running dev interface at http://localhost:12345\n\n")
	} else {
		fmt.Printf("\n🚀 Get started: http://localhost:12345/\n\n")
	}
	conf := Config{
		EngineAddr:            "http://localhost:26657",
		AppBackendAddr:        "http://localhost:1317",
		AppFrontendAddr:       "http://localhost:8080",
		DevFrontendAssetsPath: "../../ui/dist",
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
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}
