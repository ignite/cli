package ignitecmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui/entrywriter"
	"github.com/ignite/cli/ignite/pkg/xgit"
	"github.com/ignite/cli/ignite/services/plugin"
)

// plugins hold the list of plugin declared in the config.
// A global variable is used so the list is accessible to the plugin commands.
var plugins []*plugin.Plugin

// LoadPlugins tries to load all the plugins found in configuration.
// If no configuration found, it returns w/o error.
func LoadPlugins(ctx context.Context, rootCmd *cobra.Command) error {
	// NOTE(tb) Not sure if it's the right place to load this.
	chain, err := NewChainWithHomeFlags(rootCmd)
	if err != nil {
		// Binary is run outside of an chain app, plugins can't be loaded
		return nil
	}
	plugins, err = plugin.Load(ctx, chain)
	if err != nil {
		return err
	}
	// Link plugins to related commands
	var loadErrors []string
	for _, p := range plugins {
		linkPluginCmds(rootCmd, p)
		if p.Error != nil {
			loadErrors = append(loadErrors, p.Path)
		}
	}
	if len(loadErrors) > 0 {
		// unload any plugin that could have been loaded
		UnloadPlugins()
		printPlugins()
		return errors.Errorf("fail to load: %v", strings.Join(loadErrors, ","))
	}
	return nil
}

// UnloadPlugins releases any loaded plugins, which is basically killing the
// plugin server instance.
func UnloadPlugins() {
	for _, p := range plugins {
		p.KillClient()
	}
}

// linkPluginCmds tries to add the plugin commands to the legacy ignite
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
	if !strings.HasPrefix(cmdPath, "ignite") {
		// cmdPath must start with `ignite ` before comparison with
		// cmd.CommandPath()
		cmdPath = "ignite " + cmdPath
	}
	cmdPath = strings.TrimSpace(cmdPath)

	cmd := findCommandByPath(rootCmd, cmdPath)
	if cmd == nil {
		p.Error = errors.Errorf("unable to find commandPath %q for plugin %q", cmdPath, p.Path)
		return
	}
	if cmd.Runnable() {
		p.Error = errors.Errorf("can't attach plugin command %q to runnable command %q", pluginCmd.Use, cmd.CommandPath())
		return
	}
	for _, cmd := range cmd.Commands() {
		if cmd.Name() == pluginCmd.Use {
			p.Error = errors.Errorf("plugin command %q already exists in ignite's commands", pluginCmd.Use)
			return
		}
	}
	newCmd := &cobra.Command{
		Use:   pluginCmd.Use,
		Short: pluginCmd.Short,
		Long:  pluginCmd.Long,
	}
	newCmd.Flags().AddFlagSet(pluginCmd.Flags())
	cmd.AddCommand(newCmd)
	if len(pluginCmd.Commands) == 0 {
		// pluginCmd has no sub commands, so it's runnable
		newCmd.RunE = func(cmd *cobra.Command, args []string) error {
			// Pass config parameters
			pluginCmd.With = p.With
			// Pass flags
			pluginCmd.SetFlags(cmd.Flags())
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
		Use:   "update [path]",
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
				if p.Path == args[0] {
					err := plugin.Update(p)
					if err != nil {
						return err
					}
					fmt.Printf("Plugin %q updated.\n", p.Path)
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
			path, err := plugin.Scaffold(wd, moduleName)
			if err != nil {
				return err
			}
			if err := xgit.InitAndCommit(path); err != nil {
				return err
			}

			message := `
⭐️ Successfully created a new plugin '%[1]s'.
👉 update plugin code at '%[2]s/main.go'

👉 test plugin integration by adding the following lines in a chain config.yaml:
plugins:
- path: %[2]s

👉 once the plugin is pushed to a repository, the config becomes:
plugins:
- path: %[1]s
`
			fmt.Printf(message, moduleName, path)
			return nil
		},
	}
}

func printPlugins() {
	if len(plugins) == 0 {
		fmt.Println("No plugin found")
		return
	}
	var entries [][]string
	for _, p := range plugins {
		status := "✅ Loaded"
		if p.Error != nil {
			status = fmt.Sprintf("❌ Error: %v", p.Error)
		}
		entries = append(entries, []string{p.Path, status})
	}
	entrywriter.MustWrite(os.Stdout, []string{"path", "status"}, entries...)
}
