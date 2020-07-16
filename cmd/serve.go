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
	packr "github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
	"github.com/radovskyb/watcher"
	"github.com/spf13/cobra"
)

// Env ...
type Env struct {
	ChainID string `json:"chain_id"`
	NodeJS  bool   `json:"node_js"`
}

func startServe(verbose bool) (*exec.Cmd, *exec.Cmd) {
	appName, _ := getAppAndModule()
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
	cmdInitPre := exec.Command("make", "init-pre")
	if err := cmdInitPre.Run(); err != nil {
		log.Fatal("Error in initializing the chain. Please, check ./init.sh")
	}
	if verbose {
		cmdInitPre.Stdout = os.Stdout
	}
	userString1, err := exec.Command("make", "init-user1", "-s").Output()
	if err != nil {
		log.Fatal(err)
	}
	var userJSON map[string]interface{}
	json.Unmarshal(userString1, &userJSON)
	fmt.Printf("🙂 Created an account. Password (mnemonic): %[1]v\n", userJSON["mnemonic"])
	userString2, err := exec.Command("make", "init-user2", "-s").Output()
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(userString2, &userJSON)
	fmt.Printf("🙂 Created an account. Password (mnemonic): %[1]v\n", userJSON["mnemonic"])
	cmdInitPost := exec.Command("make", "init-post")
	if err := cmdInitPost.Run(); err != nil {
		log.Fatal(err)
	}
	if verbose {
		cmdInitPost.Stdout = os.Stdout
	}
	if verbose {
		cmdInitPost.Stdout = os.Stdout
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
	if !verbose {
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
		cmdt, cmdr := startServe(verbose)
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			<-c
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
					cmdt, cmdr = startServe(verbose)
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
	_, err := exec.LookPath(name)
	return err == nil
}
