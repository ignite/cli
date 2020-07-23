package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Pallinder/go-randomdata"
	"github.com/ilgooz/analytics-go"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/pkg/analyticsutil"
)

const (
	analyticsEndpoint = "https://analytics.starport.cloud"
	analyticsKey      = "ib6mwzNSLW6qIFRTyftezJL8cX4jWkQY"
)

var (
	analyticsc           *analyticsutil.Client
	starportDir          = ".starport"
	starportAnonIdentity = "anon"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "starport",
	Short: "A tool for scaffolding out Cosmos applications",
}

// Execute ...
func Execute() {
	defer func() {
		if r := recover(); r != nil {
			sendAnalytics(Metric{
				Err: fmt.Errorf("%s", r),
			})
			fmt.Println(r)
			os.Exit(1)
		}
	}()
	analyticsc = analyticsutil.New(analyticsEndpoint, analyticsKey)
	// TODO add version of new installation.
	name, hadLogin := prepLoginName()
	analyticsc.Login(name, "todo-version")
	if !hadLogin {
		sendAnalytics(Metric{
			IsInstallation: true,
		})
	}
	sendAnalytics(Metric{})
	err := rootCmd.Execute()
	sendAnalytics(Metric{
		IsExecutionDone: true,
		Err:             err,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(appCmd)
	rootCmd.AddCommand(typedCmd)
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
	appCmd.Flags().StringP("denom", "d", "token", "Token denomination")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getAppAndModule(path string) (string, string) {
	goModFile, err := ioutil.ReadFile(filepath.Join(path, "go.mod"))
	if err != nil {
		log.Fatal(err)
	}
	moduleString := strings.Split(string(goModFile), "\n")[0]
	modulePath := strings.ReplaceAll(moduleString, "module ", "")
	var appName string
	if t := strings.Split(modulePath, "/"); len(t) > 0 {
		appName = t[len(t)-1]
	}
	return appName, modulePath
}

func prepLoginName() (name string, hadLogin bool) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "any", false
	}
	anonPath := filepath.Join(home, starportDir, starportAnonIdentity)
	data, err := ioutil.ReadFile(anonPath)
	if len(data) != 0 {
		return string(data), true
	}
	name = randomdata.SillyName()
	ioutil.WriteFile(anonPath, []byte(name), 0644)
	return name, false
}

type Metric struct {
	IsInstallation  bool
	IsExecutionDone bool
	Err             error
}

func sendAnalytics(m Metric) {
	commandExecStatus := "pre"
	if m.IsExecutionDone {
		commandExecStatus = "post"
	}
	fullCommand := os.Args
	var rootCommand string
	if len(os.Args) > 1 { // first is starport (binary name).
		rootCommand = os.Args[1]
	}
	props := analytics.NewProperties()
	if !m.IsInstallation {
		props.Set("action", rootCommand)
		props.Set("label", strings.Join(fullCommand, " "))
		props.Set("commandExecStatus", commandExecStatus)
	}
	if m.Err != nil {
		props.Set("error", m.Err.Error())
	}
	var category string
	switch {
	case m.IsInstallation:
		category = "install"
	case m.Err == nil:
		category = "success"
	case m.Err != nil:
		category = "error"
	}
	props.Set("category", category)
	analyticsc.Track(analytics.Track{
		Event:      category,
		Properties: props,
	})
	if m.IsExecutionDone {
		// flush the message in the queue and close the client.
		analyticsc.Close()
	}
}
