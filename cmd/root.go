package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "starport",
	Short: "A tool for scaffolding out Cosmos applications",
}

// Execute ...
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(appCmd)
	rootCmd.AddCommand(typedCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(addCmd)
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
