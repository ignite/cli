package starportcmd

import (
	"fmt"
	"path/filepath"
	"plugin"

	"github.com/spf13/cobra"
)

func NewPluginsUse() *cobra.Command {
	c := &cobra.Command{
		Use:   "use",
		Short: "use a plugin listed in config",
		RunE:  pluginsUseHandler,
	}

	flagSetPath(c)

	return c
}

func pluginsUseHandler(cmd *cobra.Command, args []string) error {
	appPath := flagGetPath(cmd)
	pluginPath := filepath.Join(appPath, fmt.Sprintf("plugins/%s/main.so", args[0]))

	plug, err := plugin.Open(pluginPath)
	if err != nil {
		return err
	}

	newCmdsSymbol, err := plug.Lookup("NewCmds")
	if err != nil {
		return err
	}

	c := newCmdsSymbol.(func() *cobra.Command)()
	c.SetArgs(args[1:])
	_, err = c.ExecuteC()
	if err != nil {
		return err
	}

	return nil
}
