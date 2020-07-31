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
	"time"

	"github.com/pkg/errors"
	"github.com/radovskyb/watcher"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/xexec"
)

type App struct {
	Name string
	Path string
}

func Serve(ctx context.Context, app App, verbose bool) error {
	go cmdrunner.
		New().
		Run(ctx, step.New(
			step.Exec("npm", "run", "dev"),
			step.Workdir(filepath.Join(app.Path, "frontend")),
		))

	serveCtx, cancel := context.WithCancel(ctx)
	startServe(serveCtx, app, verbose) // TODO handle error
	go runDevServer(app, verbose)

	w := watcher.New()
	w.SetMaxEvents(1)
	go func() {
		for {
			select {
			case <-ctx.Done():
				w.Close()
			case <-w.Event:
				cancel()
				serveCtx, cancel = context.WithCancel(ctx)
				startServe(serveCtx, app, verbose) // TODO handle error
			case err := <-w.Error:
				log.Println(err)
			}
		}
	}()
	if err := w.AddRecursive(filepath.Join(app.Path, "./app")); err != nil {
		log.Fatalln(err)
	}
	if err := w.AddRecursive(filepath.Join(app.Path, "./cmd")); err != nil {
		log.Fatalln(err)
	}
	if err := w.AddRecursive(filepath.Join(app.Path, "./x")); err != nil {
		log.Fatalln(err)
	}
	return w.Start(time.Millisecond * 1000)
}

func startServe(ctx context.Context, app App, verbose bool) error {
	var (
		steps step.Steps

		stdout = ioutil.Discard
		stderr = ioutil.Discard

		mnemonic = &bytes.Buffer{}
	)
	if verbose {
		stdout = os.Stdout
		stderr = os.Stderr
	}

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

	if err := cmdrunner.
		New(cmdrunner.DefaultStdout(stdout),
			cmdrunner.DefaultStderr(stderr),
			cmdrunner.DefaultWorkdir(app.Path)).
		Run(ctx, steps...); err != nil {
		log.Fatal(err)
	}

	var servers step.Steps
	servers.Add(step.New(
		step.Exec(fmt.Sprintf("%[1]vd", app.Name), "start"), //nolint:gosec // Subprocess launched with function call as argument or cmd arguments
		step.InExec(func() error {
			if verbose {
				fmt.Println("🌍 Running a server at http://localhost:26657 (Tendermint)")
			} else {
				fmt.Printf("🌍 Running a Cosmos '%[1]v' app with Tendermint.\n", app.Name)
			}
			return nil
		}),
		step.PostExec(func(exitErr error) error {
			return errors.Wrapf(exitErr, "cannot run %[1]vd start", app.Name)
		}),
	))
	servers.Add(step.New(
		step.Exec(fmt.Sprintf("%[1]vcli", app.Name), "rest-server"), //nolint:gosec // Subprocess launched with function call as argument or cmd arguments
		step.InExec(func() error {
			if verbose {
				fmt.Println("🌍 Running a server at http://localhost:1317 (LCD)")
			}
			return nil
		}),
		step.PostExec(func(exitErr error) error {
			return errors.Wrapf(exitErr, "cannot run %[1]vcli rest-server", app.Name)
		}),
	))

	serverRunner := cmdrunner.New(
		cmdrunner.RunParallel(),
		cmdrunner.DefaultStdout(stdout),
		cmdrunner.DefaultStderr(stderr),
		cmdrunner.DefaultWorkdir(app.Path),
	)
	go serverRunner.Run(ctx, servers...) // TODO handle err
	return nil
}

func runDevServer(app App, verbose bool) error {
	if verbose {
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
	return http.ListenAndServe(":12345", newDevHandler(app, conf))
}
