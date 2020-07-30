package starportcmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"time"

	gocmd "github.com/go-cmd/cmd"
	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/radovskyb/watcher"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/pkg/cmdrunner/step"
)

var appPath string

func NewServe() *cobra.Command {
	c := &cobra.Command{
		Use:   "serve",
		Short: "Launches a reloading server",
		Args:  cobra.ExactArgs(0),
		Run:   serveHandler,
	}
	c.Flags().StringVarP(&appPath, "path", "p", "", "path of the app")
	c.Flags().BoolP("verbose", "v", false, "Verbose output")
	return c
}

func serveHandler(cmd *cobra.Command, args []string) {
	verbose, _ := cmd.Flags().GetBool("verbose")
	cmdNpm := gocmd.NewCmd("npm", "run", "dev")
	cmdNpm.Dir = filepath.Join(appPath, "frontend")
	cmdNpm.Start()
	cancel := startServe(appPath, verbose)
	appName, _ := getAppAndModule(appPath)
	go runDevServer(appName, verbose)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		cmdNpm.Stop()
		cancel()
		os.Exit(0)
	}()
	w := watcher.New()
	w.SetMaxEvents(1)
	go func() {
		for {
			select {
			case <-w.Event:
				cancel()
				cancel = startServe(appPath, verbose)
			case err := <-w.Error:
				log.Println(err)
			case <-w.Closed:
				return
			}
		}
	}()
	if err := w.AddRecursive(filepath.Join(appPath, "./app")); err != nil {
		log.Fatalln(err)
	}
	if err := w.AddRecursive(filepath.Join(appPath, "./cmd")); err != nil {
		log.Fatalln(err)
	}
	if err := w.AddRecursive(filepath.Join(appPath, "./x")); err != nil {
		log.Fatalln(err)
	}
	if err := w.Start(time.Millisecond * 1000); err != nil {
		log.Fatalln(err)
	}
}

// Env ...
type Env struct {
	ChainID string `json:"chain_id"`
	NodeJS  bool   `json:"node_js"`
}

func startServe(path string, verbose bool) context.CancelFunc {
	appName, _ := getAppAndModule(path)

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
			if !isCommandAvailable("go") {
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
			if !isCommandAvailable("make") {
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
			cmdrunner.DefaultWorkdir(path)).
		Run(context.Background(), steps...); err != nil {
		log.Fatal(err)
	}

	var servers step.Steps
	servers.Add(step.New(
		step.Exec(fmt.Sprintf("%[1]vd", appName), "start"), //nolint:gosec // Subprocess launched with function call as argument or cmd arguments
		step.InExec(func() error {
			if verbose {
				fmt.Println("🌍 Running a server at http://localhost:26657 (Tendermint)")
			} else {
				fmt.Printf("🌍 Running a Cosmos '%[1]v' app with Tendermint.\n", appName)
			}
			return nil
		}),
		step.PostExec(func(exitErr error) error {
			return errors.Wrapf(exitErr, "cannot run %[1]vd start", appName)
		}),
	))
	servers.Add(step.New(
		step.Exec(fmt.Sprintf("%[1]vcli", appName), "rest-server"), //nolint:gosec // Subprocess launched with function call as argument or cmd arguments
		step.InExec(func() error {
			if verbose {
				fmt.Println("🌍 Running a server at http://localhost:1317 (LCD)")
			}
			return nil
		}),
		step.PostExec(func(exitErr error) error {
			return errors.Wrapf(exitErr, "cannot run %[1]vcli rest-server", appName)
		}),
	))

	serverRunner := cmdrunner.New(
		cmdrunner.RunParallel(),
		cmdrunner.DefaultStdout(stdout),
		cmdrunner.DefaultStderr(stderr),
		cmdrunner.DefaultWorkdir(path),
	)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		if err := serverRunner.Run(ctx, servers...); err != nil {
			if _, ok := errors.Cause(err).(*exec.ExitError); !ok {
				log.Fatal(err)
			}
		}
	}()
	return cancel
}

func runDevServer(appName string, verbose bool) error {
	if verbose {
		fmt.Printf("🔧 Running dev interface at http://localhost:12345\n\n")
	} else {
		fmt.Printf("\n🚀 Get started: http://localhost:12345/\n\n")
	}
	router := mux.NewRouter()
	devUI := packr.New("ui/dist", "../../../../ui/dist")
	router.HandleFunc("/env", func(w http.ResponseWriter, r *http.Request) {
		env := Env{appName, isCommandAvailable("node")}
		js, err := json.Marshal(env)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	})
	router.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		res, err := http.Get("http://localhost:1317/node_info")
		if err != nil || res.StatusCode != 200 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 error"))
		} else if res.StatusCode == 200 {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("200 ok"))
		}
	})
	router.HandleFunc("/rpc", func(w http.ResponseWriter, r *http.Request) {
		res, err := http.Get("http://localhost:26657")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else if res.StatusCode == 200 {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
	router.HandleFunc("/frontend", func(w http.ResponseWriter, r *http.Request) {
		res, err := http.Get("http://localhost:8080")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else if res.StatusCode == 200 {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
	router.PathPrefix("/").Handler(http.FileServer(devUI))
	return http.ListenAndServe(":12345", router)
}

func isCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
