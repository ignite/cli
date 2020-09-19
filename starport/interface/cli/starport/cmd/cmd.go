package starportcmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
)

// New creates a new root command for `starport` with its sub commands.
func New() *cobra.Command {
	c := &cobra.Command{
		Use:   "starport",
		Short: "A tool for scaffolding out Cosmos applications",
	}
	c.AddCommand(NewApp())
	c.AddCommand(NewType())
	c.AddCommand(NewServe())
	c.AddCommand(NewModule())
	c.AddCommand(NewVersion())
	c.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	return c
}

const sdkVersionFlag = "sdk-version"

func addSdkVersionFlag(c *cobra.Command) {
	c.Flags().String(sdkVersionFlag, string(cosmosver.Launchpad), fmt.Sprintf("Target Cosmos-SDK Version %s", cosmosver.MajorVersions))
}

func sdkVersion(c *cobra.Command) (cosmosver.MajorVersion, error) {
	v, _ := c.Flags().GetString(sdkVersionFlag)
	parsed, err := cosmosver.MajorVersions.Parse(v)
	if err != nil {
		return "", fmt.Errorf("%q is an unkown sdk version", v)
	}
	return parsed, nil
}

func getModule(path string) string {
	goModFile, err := ioutil.ReadFile(filepath.Join(path, "go.mod"))
	if err != nil {
		log.Fatal(err)
	}
	moduleString := strings.Split(string(goModFile), "\n")[0]
	modulePath := strings.ReplaceAll(moduleString, "module ", "")
	return modulePath
}
