package cmd

import (
	"encoding/json"
	"fmt"
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
)

// Env ...
type Env struct {
	ChainID string `json:"chain_id"`
	NodeJS  bool   `json:"node_js"`
}

func startServe(verbose bool) (*exec.Cmd, *exec.Cmd, *gocmd.Cmd) {
	appName, _ := getAppAndModule()
	cmdNpm := gocmd.NewCmd("npm", "run", "dev")
	cmdNpm.Dir = "frontend"
	cmdNpm.Start()
	fmt.Printf("\n📦 Installing dependencies...\n")
	cmdMod := exec.Command("/bin/sh", "-c", "go mod tidy")
	if verbose {
		cmdMod.Stdout = os.Stdout
	}
	if err := cmdMod.Run(); err != nil {
		log.Fatal("Error running go mod tidy. Please, check ./go.mod")
	}
	fmt.Printf("🚧 Building the application...\n")
	cmdMake := exec.Command("/bin/sh", "-c", "make")
	if verbose {
		cmdMake.Stdout = os.Stdout
	}
	if err := cmdMake.Run(); err != nil {
		log.Fatal("Error in building the application. Please, check ./Makefile")
	}
	fmt.Printf("💫 Initializing the chain...\n")
	cmdInit := exec.Command("/bin/sh", "-c", "sh init.sh")
	if verbose {
		cmdInit.Stdout = os.Stdout
	}
	if err := cmdInit.Run(); err != nil {
		log.Fatal("Error in initializing the chain. Please, check ./init.sh")
	}
	cmdTendermint := exec.Command(fmt.Sprintf("%[1]vd", appName), "start") //nolint:gosec // Subprocess launched with function call as argument or cmd arguments
	if verbose {
		fmt.Printf("🌍 Running a server at http://localhost:26657 (Tendermint)\n")
		cmdTendermint.Stdout = os.Stdout
	} else {
		fmt.Printf("🌍 Running a Cosmos '%[1]v' app with Tendermint.\n", appName)
	}
	if err := cmdTendermint.Start(); err != nil {
		log.Fatal(fmt.Sprintf("Error in running %[1]vd start", appName), err)
	}
	cmdREST := exec.Command(fmt.Sprintf("%[1]vcli", appName), "rest-server") //nolint:gosec // Subprocess launched with function call as argument or cmd arguments
	if verbose {
		fmt.Printf("🌍 Running a server at http://localhost:1317 (LCD)\n")
		cmdREST.Stdout = os.Stdout
	}
	if err := cmdREST.Start(); err != nil {
		log.Fatal(fmt.Sprintf("Error in running %[1]vcli rest-server", appName))
	}
	if verbose {
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
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else if res.StatusCode == 200 {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
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
	if !verbose {
		fmt.Printf("\n🚀 Get started: http://localhost:12345/\n\n")
	}
	return cmdTendermint, cmdREST, cmdNpm
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Launches a reloading server",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		cmdt, cmdr, cmdn := startServe(verbose)
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			<-c
			cmdn.Stop()
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
					cmdn.Stop()
					cmdr.Process.Kill()
					cmdt.Process.Kill()
					cmdt, cmdr, cmdn = startServe(verbose)
				case err := <-w.Error:
					log.Fatalln(err)
				case <-w.Closed:
					return
				}
			}
		}()
		if err := w.AddRecursive("."); err != nil {
			log.Fatalln(err)
		}
		if err := w.Ignore("./frontend"); err != nil {
			log.Fatalln(err)
		}
		if err := w.Ignore("./.git"); err != nil {
			log.Fatalln(err)
		}
		if err := w.Start(time.Millisecond * 100); err != nil {
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
