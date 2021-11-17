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
	log.Println("Handle plugin Install", args)
	var conf, _ = GetConfig()
	if len(conf.Plugins) == 0 {
		fmt.Println("There's no plugins to be implemented.")
		return nil
	}

	loader, err := plugin.NewLoader()
	if err != nil {
		return err
	}
	var pluginIdx = -1
	// TODO: jkkim: Search plugin with name from []PluginConfig.
	// How to set a selected Plugins? always first plugin? or need to find by its name?

	for index, pluginSingle := range conf.Plugins {
		if args[0] == pluginSingle.Name {
			pluginIdx = index
			break
		}
	}

	if pluginIdx == -1 {
		fmt.Println("There's no plugin with given name")
		return err
	}

	selectedPlugin := conf.Plugins[pluginIdx]

	// TODO: jkkim: Check installed by call PluginConfig.IsInstalled()
	isInstalled := loader.IsInstalled(selectedPlugin)
	if isInstalled {
		fmt.Println("Selected Plugins", selectedPlugin.Name)
		return nil
	}

	builder, err := plugin.NewBuilder()
	if err != nil {
		return err
	}

	// TODO: jkkim: Install plugin Builder.Build()
	err = builder.Build(selectedPlugin)
	if err != nil {
		return err
	}
	return nil
}

func (p *pluginCmdHandler) HandleList(cmd *cobra.Command, args []string) error {
	log.Println("HandleList: printout all list in configs")
	path, err := chainconfig.ConfigDirPath() // check if there's any config.yml
	if err != nil {
		return err
	}
	conf, err := chainconfig.ParseFile(filepath.Join(path, "config.yml"))
	// var tempPath string = "/Users/dongyookang/ffplay/GolandProjects/starport-development/mars"
	// conf, err := chainconfig.ParseFile(filepath.Join(tempPath, "config.yml"))

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
