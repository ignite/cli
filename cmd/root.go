package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Pallinder/go-randomdata"
	"github.com/segmentio/analytics-go"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/pkg/analyticsutil"
)

const (
	analyticsEndpoint = "https://analytics.starport.cloud"
	analyticsKey      = "pWSXBMIF3tQsHTtA63Lb63zAfIA80Bhy"
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
			sendAnalytics(true, fmt.Errorf("%s", r))
			fmt.Println(r)
			os.Exit(1)
		}
	}()
	sendAnalytics(false, nil)
	err := rootCmd.Execute()
	sendAnalytics(true, err)
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

	analyticsc = analyticsutil.New(analyticsEndpoint, analyticsKey)
	// TODO add starport version.
	analyticsc.Login(loginName(), "todo")
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

func loginName() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "any"
	}
	anonPath := filepath.Join(home, starportDir, starportAnonIdentity)
	data, err := ioutil.ReadFile(anonPath)
	if len(data) == 0 {
		name := randomdata.SillyName()
		ioutil.WriteFile(anonPath, []byte(name), 0644)
		return name
	}
	return string(data)
}

func sendAnalytics(isDone bool, err error) {
	hook := "pre"
	if isDone {
		hook = "post"
	}
	props := analytics.NewProperties().
		Set("name", strings.Join(os.Args, " ")).
		Set("hook", hook)
	if err != nil {
		props.Set("err", err.Error())
	}
	analyticsc.Track(analytics.Track{
		Event:      "command",
		Properties: props,
	})
	if isDone {
		// flush the message in the queue and close the client.
		analyticsc.Close()
	}
}
