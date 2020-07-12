package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"time"

	gocmd "github.com/go-cmd/cmd"
	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
	"github.com/radovskyb/watcher"
	"github.com/spf13/cobra"
	sr "github.com/tendermint/starport/pkg/cmdsteprunner"
	"github.com/tendermint/starport/pkg/lineprefixer"
)

// Env ...
type Env struct {
	ChainID string `json:"chain_id"`
	NodeJS  bool   `json:"node_js"`
}

// ctx for cancellation
// controlling stdout, stderr for testability.
func startServe(ctx context.Context, stdout, stderr io.Writer, isVerbose bool) (*exec.Cmd, *exec.Cmd) {
	// outOption logs step outputs (both regular and error) with a prefix for each log line
	// in the verbose mode.
	//
	// prefixing with line is helpful to debug long logs.
	//
	// we can also customize this func to color/style the lines or prefixes.
	outOption := func(prefix string) sr.StepOption {
		var stdout io.Writer = lineprefixer.NewWriter(prefix, stdout)
		var stderr io.Writer = lineprefixer.NewWriter(prefix, stderr)
		if !isVerbose {
			stdout = ioutil.Discard
			stderr = ioutil.Discard
		}
		return sr.StepStdouterr(stdout, stderr)
	}
	// initMessagePrinter is a shortcut to register a hook to the StepPreExec in order to
	// only print a message.
	initMessagePrinter := func(message string) func() error {
		fmt.Fprint(stdout, message)
		return nil

	}
	// proxyError is a shortcut to register a hook to the StepAfterExec in order to
	// replace command's execution error with a dummy(hidden) error.
	proxyError := func(message string) func(error) error {
		// err arg can be utilized to provide more details.
		// for ex exitErr := err.(*cmd.ExitError: https://golang.org/pkg/os/exec/#ExitError),
		// do something with the exit code to provide more detail to returned error.
		return func(err error) error {
			if err != nil {
				return errors.New(message)
			}
			return nil
		}
	}

	// NOTE: the above three funcs can be placed under a util pkg (pkg/a-util-pkg) in order to
	// reduce code from this package so it can be easier to read and understand main
	// functionality here.

	// userBuf keeps user data to be used later.
	userBuf := &bytes.Buffer{}

	// list of steps to run.
	steps := []sr.Step{
		sr.NewStep(
			sr.StepCommand("go", "mod", "tidy"),
			outOption("tidy: "),
			sr.StepPreExec(initMessagePrinter("\n📦 Installing dependencies...\n")),
			sr.StepAfterExec(proxyError("Error running go mod tidy. Please, check ./go.mod")),
		),
		sr.NewStep(
			sr.StepCommand("make"),
			outOption("all: "),
			sr.StepPreExec(initMessagePrinter("🚧 Building the application...\n")),
			sr.StepAfterExec(proxyError("Error in building the application. Please, check ./Makefile")),
		),
		sr.NewStep(
			sr.StepCommand("make", "init-pre"),
			outOption("init pre: "),
			sr.StepPreExec(initMessagePrinter("💫 Initializing the chain...\n")),
			sr.StepAfterExec(proxyError("Error in initializing the chain. Please, check ./init.sh")),
		),
		sr.NewStep(
			sr.StepCommand("make", "init-user", "-s"),
			outOption("init user: "),
			sr.StepStdout(userBuf), // this overwrites the given stdout in the above line but not stderr by intention.
			sr.StepAfterExec(func(err error) error {
				if err != nil { // if command exited with a non-zero code, nothing to do.
					return err
				}
				var user map[string]interface{}
				if err := json.NewDecoder(userBuf).Decode(&user); err != nil {
					return err
				}
				fmt.Fprintf(stdout, "🙂 Created an account. Password (mnemonic): %[1]v\n", user["mnemonic"])
				return nil
			}),
		),
		sr.NewStep(
			sr.StepCommand("make", "init-post"),
			outOption("init post: "),
		),
	}

	// New(options...) can have options to optionally run steps in parallel or for ex. continue
	// to other steps if previous one exited with a failure.
	r := sr.New()

	// ctx for cancallation. if context is cancalled by a timeout or programatically, next steps should not run
	// and current one should be interrupted. intrreption specially useful when we may want to use cmdsteprunner
	// for 'serve' functionality (starting long running servers).
	if err := r.Run(ctx, steps...); err != nil {
		// instead of this we should return with the error and Fatal in the caller function.
		// because a Fatal will end the program and we need to use it carefully. returning error
		// helps to testability as well.
		log.Fatal(err)
	}

	// NOTE: following code not relavent to this proposal.
	appName, _ := getAppAndModule()
	cmdTendermint := exec.Command(fmt.Sprintf("%[1]vd", appName), "start") //nolint:gosec // Subprocess launched with function call as argument or cmd arguments
	if isVerbose {
		fmt.Printf("🌍 Running a server at http://localhost:26657 (Tendermint)\n")
		cmdTendermint.Stdout = os.Stdout
	} else {
		fmt.Printf("🌍 Running a Cosmos '%[1]v' app with Tendermint.\n", appName)
	}
	if err := cmdTendermint.Start(); err != nil {
		log.Fatal(fmt.Sprintf("Error in running %[1]vd start", appName), err)
	}
	cmdREST := exec.Command(fmt.Sprintf("%[1]vcli", appName), "rest-server") //nolint:gosec // Subprocess launched with function call as argument or cmd arguments
	if isVerbose {
		fmt.Printf("🌍 Running a server at http://localhost:1317 (LCD)\n")
		cmdREST.Stdout = os.Stdout
	}
	if err := cmdREST.Start(); err != nil {
		log.Fatal(fmt.Sprintf("Error in running %[1]vcli rest-server", appName))
	}
	if isVerbose {
		fmt.Printf("🔧 Running dev interface at http://localhost:12345\n\n")
	}
	router := mux.NewRouter()
	devUI := packr.New("ui/dist", "../ui/dist")
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
	go func() {
		http.ListenAndServe(":12345", router)
	}()
	if !isVerbose {
		fmt.Printf("\n🚀 Get started: http://localhost:12345/\n\n")
	}
	return cmdTendermint, cmdREST
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Launches a reloading server",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		cmdNpm := gocmd.NewCmd("npm", "run", "dev")
		cmdNpm.Dir = "frontend"
		cmdNpm.Start()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cmdt, cmdr := startServe(ctx, os.Stdout, os.Stderr, verbose)
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			<-c
			// later in other PR, we can use cancel() here to cancel long running procceses.
			cmdNpm.Stop()
			cmdr.Process.Kill()
			cmdt.Process.Kill()
			os.Exit(0)
		}()
		w := watcher.New()
		w.SetMaxEvents(1)
		go func() {
			for {
				select {
				case <-w.Event:
					cmdr.Process.Kill()
					cmdt.Process.Kill()
					cmdt, cmdr = startServe(ctx, os.Stdout, os.Stderr, verbose)
				case err := <-w.Error:
					log.Println(err)
				case <-w.Closed:
					return
				}
			}
		}()
		if err := w.AddRecursive("./app"); err != nil {
			log.Fatalln(err)
		}
		if err := w.AddRecursive("./cmd"); err != nil {
			log.Fatalln(err)
		}
		if err := w.AddRecursive("./x"); err != nil {
			log.Fatalln(err)
		}
		if err := w.Start(time.Millisecond * 1000); err != nil {
			log.Fatalln(err)
		}
	},
}

func isCommandAvailable(name string) bool {
	cmd := exec.Command("/bin/sh", "-c", "command -v "+name)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
