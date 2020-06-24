package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
	"github.com/radovskyb/watcher"
	"github.com/spf13/cobra"
)

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
	cmdInit := exec.Command("/bin/sh", "-c", "sh init.sh")
	if verbose {
		cmdInit.Stdout = os.Stdout
	}
	if err := cmdInit.Run(); err != nil {
		log.Fatal("Error in initializing the chain. Please, check ./init.sh")
	}
	fmt.Printf("🎨 Created a web front-end.\n")
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
	router.HandleFunc("/chain_id", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, appName)
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
		cmdt, cmdr := startServe(verbose)
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
					log.Fatalln(err)
				case <-w.Closed:
					return
				}
			}
		}()
		if err := w.AddRecursive("."); err != nil {
			log.Fatalln(err)
		}
		if err := w.Ignore("./ui"); err != nil {
			log.Fatalln(err)
		}
		if err := w.Start(time.Millisecond * 100); err != nil {
			log.Fatalln(err)
		}
	},
}
