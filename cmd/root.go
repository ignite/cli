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
			addMetric(Metric{
				Err: fmt.Errorf("%s", r),
			})
			analyticsc.Close()
			fmt.Println(r)
			os.Exit(1)
		}
	}()
	analyticsc = analyticsutil.New(analyticsEndpoint, analyticsKey)
	// TODO add version of new installation.
	name, hadLogin := prepLoginName()
	analyticsc.Login(name, "todo-version")
	if !hadLogin {
		addMetric(Metric{
			Login:          name,
			IsInstallation: true,
		})
	}
	if len(os.Args) > 1 && os.Args[1] == "serve" {
		addMetric(Metric{})
	}
	err := rootCmd.Execute()
	addMetric(Metric{
		Err: err,
	})
	analyticsc.Close()
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
	// IsInstallation sets metrics type as an installation metric.
	IsInstallation bool

	// Err sets metrics type as an error metric.
	Err error

	// Login is the name of anon user.
	Login string
}

func addMetric(m Metric) {
	fullCommand := os.Args
	var rootCommand string
	if len(os.Args) > 1 { // first is starport (binary name).
		rootCommand = os.Args[1]
	}
	var event string
	props := analytics.NewProperties()
	switch {
	case m.IsInstallation:
		event = "install"
	case m.Err == nil:
		event = "success"
	case m.Err != nil:
		event = "error"
		props.Set("value", m.Err.Error())
	}
	if m.IsInstallation {
		props.Set("action", m.Login)
	} else {
		props.
			Set("action", rootCommand).
			Set("label", strings.Join(fullCommand, " "))
	}
	props.Set("category", event)
	analyticsc.Track(analytics.Track{
		Event:      event,
		Properties: props,
	})
}
