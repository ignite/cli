package starportserve

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/fswatcher"
	"github.com/tendermint/starport/starport/pkg/xexec"
)

type App struct {
	Name string
	Path string
}

type starportServe struct {
	ctx     context.Context
	app     App
	verbose bool
}

// Serve serves user apps.
func Serve(ctx context.Context, app App, verbose bool) error {
	s := &starportServe{
		ctx:     ctx,
		app:     app,
		verbose: verbose,
	}

	serveCtx, cancel := context.WithCancel(ctx)
	s.serve(serveCtx) // TODO handle error

	go s.watchAppFrontend()
	go s.runDevServer()

	changeHook := func() {
		cancel()
		serveCtx, cancel = context.WithCancel(ctx)
		s.serve(serveCtx) // TODO handle error
	}
	return fswatcher.Watch(
		ctx,
		[]string{"app", "cmd", "x"},
		fswatcher.Workdir(app.Path),
		fswatcher.OnChange(changeHook),
	)
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
		log.Fatal(err)
	}

	go cmdrunner.
		New(append(opts, cmdrunner.RunParallel())...).
		Run(ctx, s.serverSteps()...) // TODO handle err
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

func (s *starportServe) watchAppFrontend() {
	cmdrunner.
		New().
		Run(s.ctx, step.New(
			step.Exec("npm", "run", "dev"),
			step.Workdir(filepath.Join(s.app.Path, "frontend")),
		))
}

func (s *starportServe) runDevServer() error {
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
	return http.ListenAndServe(":12345", newDevHandler(s.app, conf))
}
