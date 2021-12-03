package starportcmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"

	"github.com/tendermint/starport/starport/chainconfig"
	"github.com/tendermint/starport/starport/services/plugin"
)

// TODO: Log issues.
// What is common method to log on Starport?

var (
	pluginHandler PluginCmdHandler = &pluginCmdHandler{}
)

// PluginCmdHandler provides interfaces to handle subcommands of plugin.
type PluginCmdHandler interface {
	HandleInstall(cmd *cobra.Command, args []string) error
	HandleList(cmd *cobra.Command, args []string) error
}

type pluginCmdHandler struct {
}

// GetConfig returns starport's config.
func GetConfig() (chainconfig.Config, error) {
	// ignore all path other than /usr/<user_name>/.starport
	projectPath, err := chainconfig.ConfigDirPath()
	if err != nil {
		return chainconfig.Config{}, nil
	}

	confPath, err := chainconfig.LocateDefault(projectPath)
	if err != nil {
		return chainconfig.Config{}, nil
	}

	conf, err := chainconfig.ParseFile(confPath)
	if err != nil {
		return chainconfig.Config{}, nil
	}

	if len(conf.Plugins) == 0 {
		fmt.Println("There's no plugins to be implemented.")
		return chainconfig.Config{}, nil
	}

	return conf, nil
}

func (p *pluginCmdHandler) HandleInstall(cmd *cobra.Command, args []string) error {
	var conf, _ = GetConfig()
	if len(conf.Plugins) == 0 {
		log.Println("There's no plugins to be implemented.")
		return nil
	}

	loader, err := plugin.NewLoader()
	if err != nil {
		log.Println("NewLoader", err)
		return err
	}

	var pluginIdx = -1
	for index, pluginSingle := range conf.Plugins {
		if args[0] == pluginSingle.Name {
			pluginIdx = index
			break
		}
	}

	if pluginIdx == -1 {
		log.Println("There's no plugin with given name")
		return err
	}

	selectedPlugin := conf.Plugins[pluginIdx]

	isInstalled := loader.IsInstalled(selectedPlugin)
	if isInstalled {
		log.Printf("Plugins %s already installed", selectedPlugin.Name)
		return nil
	}

	builder, err := plugin.NewBuilder()
	if err != nil {
		log.Println("NewBuilder", err)
		return err
	}

	err = builder.Build(selectedPlugin)
	if err != nil {
		log.Println("Build", err)
		return err
	}

	return nil
}

func (p *pluginCmdHandler) HandleList(cmd *cobra.Command, args []string) error {
	path, err := chainconfig.ConfigDirPath() // check if there's any config.yml
	if err != nil {
		return err
	}
	conf, err := chainconfig.ParseFile(filepath.Join(path, "config.yml"))

	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Name", "Repository URL", "Description"})

	if len(conf.Plugins) == 0 {
		rowConfigAutoMerge := table.RowConfig{AutoMerge: true}
		var noData = "No Plugin Data"
		t.AppendRow(table.Row{noData, noData, noData, noData}, rowConfigAutoMerge)

	} else {
		for index, pluginSingle := range conf.Plugins {
			t.AppendRows([]table.Row{
				{index, pluginSingle.Name, pluginSingle.RepositoryURL, pluginSingle.Description},
			})
		}
	}
	t.Render()
	return nil
}

// NewPlugin creates a new plugin command to manage plugin.
func NewPlugin() *cobra.Command {
	c := &cobra.Command{
		Use:   "plugin",
		Short: "Plugin list and install.",
		Args:  cobra.ExactArgs(1),
	}

	c.AddCommand(pluginListCmd())
	c.AddCommand(pluginInstallCmd())

	return c
}

func pluginListCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "list",
		Short: "List plugins cofigured",
		RunE:  pluginHandler.HandleList,
	}

	return c
}

func pluginInstallCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "install",
		Short: "Install new plugin",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return pluginHandler.HandleInstall(cmd, args)
		},
	}

	return c
}
