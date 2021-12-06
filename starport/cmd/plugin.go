package starportcmd

import (
	"fmt"
	"os"
	"path"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/services/plugin"
)

// plugins hold the list of plugin declared in the config.
// A global variable is used so the list is accessible to the plugin commands.
var plugins []*plugin.Plugin

// loadPlugins tries to load all the plugins found in configuration.
// If no configuration found, it returns w/o error.
func loadPlugins(rootCmd *cobra.Command) error {
	// NOTE(tb) Not sure if it's the right place to load this.
	c, err := NewChainWithHomeFlags(rootCmd)
	if err != nil {
		// Binary is run outside of an chain app, plugins can't be loaded
		return nil
	}
	plugins, err = plugin.Load(c)
	if err != nil {
		return err
	}
	// Link plugins to related commands
	var loadErrors bool
	for _, p := range plugins {
		linkPluginCmds(rootCmd, p)
		if p.Error != nil {
			loadErrors = true
		}
	}
	if loadErrors {
		fmt.Println("Error(s) detected during plugin load:")
		printPlugins()
		os.Exit(1)
	}
	return nil
}

// linkPluginCmds tries to add the plugin commands to the legacy starport
// commands.
func linkPluginCmds(rootCmd *cobra.Command, p *plugin.Plugin) {
	if p.Error != nil {
		return
	}
	for _, pluginCmd := range p.Interface.Commands() {
		linkPluginCmd(rootCmd, p, pluginCmd)
		if p.Error != nil {
			return
		}
	}
}

func linkPluginCmd(rootCmd *cobra.Command, p *plugin.Plugin, pluginCmd plugin.Command) {
	cmdPath := pluginCmd.PlaceCommandUnder
	if !strings.HasPrefix(cmdPath, "starport") {
		// cmdPath must start with `starport ` before comparison with
		// cmd.CommandPath()
		cmdPath = "starport " + cmdPath
	}
	cmdPath = strings.TrimSpace(cmdPath)

	cmd := findCommandByPath(rootCmd, cmdPath)
	if cmd == nil {
		p.Error = errors.Errorf("unable to find commandPath %q for plugin %q", cmdPath, p.Name)
		return
	}
	if cmd.Runnable() {
		p.Error = errors.Errorf("can't attach plugin command %q to runnable command %q", pluginCmd.Use, cmd.CommandPath())
		return
	}
	for _, cmd := range cmd.Commands() {
		if cmd.Name() == pluginCmd.Use {
			p.Error = errors.Errorf("plugin command %q already exists in starport's commands", pluginCmd.Use)
			return
		}
	}
	newCmd := &cobra.Command{
		Use:   pluginCmd.Use,
		Short: pluginCmd.Short,
		Long:  pluginCmd.Long,
	}
	cmd.AddCommand(newCmd)
	if len(pluginCmd.Commands) == 0 {
		// pluginCmd has no sub commands, so it's runnable
		newCmd.RunE = func(cmd *cobra.Command, args []string) error {
			// Pass config parameters
			pluginCmd.With = p.With
			// Pass cobra cmd
			pluginCmd.CobraCmd = cmd
			// Call the plugin Execute
			err := p.Interface.Execute(pluginCmd, args)
			// NOTE(tb): This pause gives enough time for go-plugin to sync the
			// output from stdout/stderr of the plugin. Without that pause, this
			// output can be discarded and not printed in the user console.
			time.Sleep(100 * time.Millisecond)
			return err
		}
	} else {
		for _, pluginCmd := range pluginCmd.Commands {
			pluginCmd.PlaceCommandUnder = newCmd.CommandPath()
			linkPluginCmd(newCmd, p, pluginCmd)
			if p.Error != nil {
				return
			}
		}
	}
}

func findCommandByPath(cmd *cobra.Command, cmdPath string) *cobra.Command {
	if cmd.CommandPath() == cmdPath {
		return cmd
	}
	for _, cmd := range cmd.Commands() {
		if cmd := findCommandByPath(cmd, cmdPath); cmd != nil {
			return cmd
		}
	}
	return nil
}

// NewPlugin returns a command that groups plugin related sub commands.
func NewPlugin() *cobra.Command {
	c := &cobra.Command{
		Use:   "plugin [command]",
		Short: "Handle plugins",
	}

	c.AddCommand(NewPluginList())
	c.AddCommand(NewPluginUpdate())
	c.AddCommand(NewPluginScaffold())
	return c
}

func NewPluginList() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List declared plugins and status",
		RunE: func(cmd *cobra.Command, args []string) error {
			printPlugins()
			return nil
		},
	}
}

func NewPluginUpdate() *cobra.Command {
	return &cobra.Command{
		Use:   "update [name]",
		Short: "Update plugins",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				// update all plugins
				err := plugin.Update(plugins...)
				if err != nil {
					return err
				}
				fmt.Printf("All plugins updated.\n")
				return nil
			}
			// find the plugin to update
			for _, p := range plugins {
				if p.Name == args[0] {
					err := plugin.Update(p)
					if err != nil {
						return err
					}
					fmt.Printf("Plugin %q updated.\n", p.Name)
					return nil
				}
			}
			return errors.Errorf("Plugin %q not found", args[0])
		},
	}
}

func NewPluginScaffold() *cobra.Command {
	return &cobra.Command{
		Use:   "scaffold [github.com/org/repo]",
		Short: "Scaffold a new plugin",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			wd, err := os.Getwd()
			if err != nil {
				return err
			}
			moduleName := args[0]
			err = plugin.Scaffold(wd, moduleName)
			if err != nil {
				return err
			}
			name := path.Base(moduleName)
			fullpath := path.Join(wd, name)

			message := `
⭐️ Successfully created a new plugin '%[1]s'.
👉 update plugin code at '%[2]s/main.go'

👉 test plugin integration by adding the following lines in a chain config.yaml:
plugins:
  - name: %[1]s
    path: %[2]s

👉 once the plugin is pushed to a repository, the config becomes:
plugins:
  - name: %[1]s
    path: %[3]s
`
			fmt.Printf(message, name, fullpath, moduleName)
			return nil
		},
	}
}

func printPlugins() {
	if len(plugins) == 0 {
		return
	}
	w := &tabwriter.Writer{}
	w.Init(os.Stdout, 0, 8, 1, '\t', 0)

	fmt.Fprintln(w, "name\tpath\tstatus")

	for _, p := range plugins {
		status := "✅ Loaded"
		if p.Error != nil {
			status = fmt.Sprintf("❌ Error: %v", p.Error)
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n",
			p.Name,
			p.Path,
			status,
		)
	}

	fmt.Fprintln(w)
	w.Flush()
}
