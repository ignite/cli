package starportcmd

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// New creates a new root command for `starport` with its sub commands.
func New() *cobra.Command {
	c := &cobra.Command{
		Use:   "starport",
		Short: "A tool for scaffolding out Cosmos applications",
	}
	c.AddCommand(NewApp())
	c.AddCommand(typedCmd)
	c.AddCommand(NewServe())
	c.AddCommand(addCmd)
	c.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	return c
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
