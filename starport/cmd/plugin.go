package starportcmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/gomodulepath"
)

// NewPlugin returns a command that groups sub commands related to chain plugins.
func NewPlugin() *cobra.Command {
	c := &cobra.Command{
		Use:   "plugin [command]",
		Short: "Manage plugins specified in config file.",
		Long:  `Manage plugins specified in config file.`,
		Args:  cobra.ExactArgs(1),
	}

	c.AddCommand(NewPluginReload())
	c.AddCommand(NewPluginPull())
	c.AddCommand(NewPluginBuild())
	c.AddCommand(NewPluginList())

	return c
}

func promptConfig() string {
	var configFile string
	fmt.Println("We didn't find your config file. What is it's name? ")
	fmt.Scanln(configFile)
	return configFile
}

func appPathFromCmd(cmd *cobra.Command) (string, error) {
	flagPath := flagGetPath(cmd)
	absPath, err := filepath.Abs(flagPath)
	if err != nil {
		return "", err
	}

	_, appPath, err := gomodulepath.Find(absPath)
	if err != nil {
		return "", err
	}

	return appPath, nil
}
