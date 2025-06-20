package ignitecmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	pluginsconfig "github.com/ignite/cli/v29/ignite/config/plugins"
	"github.com/ignite/cli/v29/ignite/pkg/clictx"
	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosanalysis"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/gomodule"
	"github.com/ignite/cli/v29/ignite/pkg/xfilepath"
	"github.com/ignite/cli/v29/ignite/pkg/xgit"
	"github.com/ignite/cli/v29/ignite/services/chain"
	"github.com/ignite/cli/v29/ignite/services/plugin"
)

const (
	flagPluginsGlobal = "global"
)

// plugins hold the list of plugin declared in the config.
// A global variable is used so the list is accessible to the plugin commands.
var plugins []*plugin.Plugin

// LoadPlugins tries to load all the plugins found in configurations.
// If no configurations found, it returns w/o error.
func LoadPlugins(ctx context.Context, cmd *cobra.Command, session *cliui.Session) error {
	var pluginsConfigs []pluginsconfig.Plugin
	localCfg, err := parseLocalPlugins()
	if err != nil && !errors.As(err, &cosmosanalysis.ErrPathNotChain{}) {
		return err
	} else if err == nil {
		pluginsConfigs = append(pluginsConfigs, localCfg.Apps...)
	}

	globalCfg, err := parseGlobalPlugins()
	if err == nil {
		pluginsConfigs = append(pluginsConfigs, globalCfg.Apps...)
	}
	ensureDefaultPlugins(cmd, globalCfg)

	if len(pluginsConfigs) == 0 {
		return nil
	}

	uniquePlugins := pluginsconfig.RemoveDuplicates(pluginsConfigs)
	plugins, err = plugin.Load(ctx, uniquePlugins, plugin.CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}
	if len(plugins) == 0 {
		return nil
	}

	return linkPlugins(ctx, cmd.Root(), plugins)
}

func parseLocalPlugins() (*pluginsconfig.Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, errors.Errorf("parse local apps: %w", err)
	}

	if err := cosmosanalysis.IsChainPath(wd); err != nil {
		return nil, err
	}

	return pluginsconfig.ParseDir(wd)
}

func parseGlobalPlugins() (cfg *pluginsconfig.Config, err error) {
	globalDir, err := plugin.PluginsPath()
	if err != nil {
		return cfg, err
	}

	cfg, err = pluginsconfig.ParseDir(globalDir)
	// if there is error parsing, return empty config and continue execution to load
	// local plugins if they exist.
	if err != nil {
		return &pluginsconfig.Config{}, nil
	}

	for i := range cfg.Apps {
		cfg.Apps[i].Global = true
	}
	return
}

func linkPlugins(ctx context.Context, rootCmd *cobra.Command, plugins []*plugin.Plugin) error {
	// Link plugins to related commands
	var linkErrors []*plugin.Plugin
	for _, p := range plugins {
		if p.Error != nil {
			linkErrors = append(linkErrors, p)
			continue
		}

		manifest, err := p.Interface.Manifest(ctx)
		if err != nil {
			p.Error = err
			linkErrors = append(linkErrors, p)
			continue
		}

		linkPluginHooks(rootCmd, p, manifest.Hooks)
		if p.Error != nil {
			linkErrors = append(linkErrors, p)
			continue
		}

		linkPluginCmds(rootCmd, p, manifest.Commands)
		if p.Error != nil {
			linkErrors = append(linkErrors, p)
			continue
		}
	}

	if len(linkErrors) > 0 {
		// unload any plugin that could have been loaded
		defer UnloadPlugins()

		if err := printPlugins(ctx, cliui.New(cliui.WithStdout(os.Stdout))); err != nil {
			// content of loadErrors is more important than a print error, so we don't
			// return here, just print the error.
			fmt.Printf("fail to print: %v\n", err)
		}

		var s strings.Builder
		for _, p := range linkErrors {
			fmt.Fprintf(&s, "%s: %v", p.Path, p.Error)
		}
		return errors.Errorf("fail to link: %v", s.String())
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

func linkPluginHooks(rootCmd *cobra.Command, p *plugin.Plugin, hooks []*plugin.Hook) {
	if p.Error != nil {
		return
	}
	for _, hook := range hooks {
		linkPluginHook(rootCmd, p, hook)
	}
}

func linkPluginHook(rootCmd *cobra.Command, p *plugin.Plugin, hook *plugin.Hook) {
	cmdPath := hook.CommandPath()
	cmd := findCommandByPath(rootCmd, cmdPath)
	if cmd == nil {
		p.Error = errors.Errorf("unable to find command path %q for app hook %q", cmdPath, hook.Name)
		return
	}
	if !cmd.Runnable() {
		p.Error = errors.Errorf("can't attach app hook %q to non executable command %q", hook.Name, hook.PlaceHookOn)
		return
	}

	newExecutedHook := func(hook *plugin.Hook, cmd *cobra.Command, args []string) *plugin.ExecutedHook {
		hook.ImportFlags(cmd)
		execHook := &plugin.ExecutedHook{
			Hook: hook,
			ExecutedCommand: &plugin.ExecutedCommand{
				Use:    cmd.Use,
				Path:   cmd.CommandPath(),
				Args:   args,
				OsArgs: os.Args,
				With:   p.With,
				Flags:  hook.Flags,
			},
		}
		execHook.ExecutedCommand.ImportFlags(cmd)
		return execHook
	}

	for _, f := range hook.Flags {
		var fs *flag.FlagSet
		if f.Persistent {
			fs = cmd.PersistentFlags()
		} else {
			fs = cmd.Flags()
		}

		if err := f.ExportToFlagSet(fs); err != nil {
			p.Error = errors.Errorf("can't attach hook flags %q to command %q", hook.Flags, hook.PlaceHookOn)
			return
		}
	}

	preRun := cmd.PreRunE
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if preRun != nil {
			err := preRun(cmd, args)
			if err != nil {
				return err
			}
		}

		api, err := newAppClientAPI(cmd)
		if err != nil {
			return err
		}

		ctx := cmd.Context()
		execHook := newExecutedHook(hook, cmd, args)
		err = p.Interface.ExecuteHookPre(ctx, execHook, api)
		if err != nil {
			return errors.Errorf("app %q ExecuteHookPre() error: %w", p.Path, err)
		}
		return nil
	}

	runCmd := cmd.RunE
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if runCmd != nil {
			err := runCmd(cmd, args)
			// if the command has failed the `PostRun` will not execute. here we execute the cleanup step before returning.
			if err != nil {
				api, err := newAppClientAPI(cmd)
				if err != nil {
					return err
				}

				ctx := cmd.Context()
				execHook := newExecutedHook(hook, cmd, args)
				err = p.Interface.ExecuteHookCleanUp(ctx, execHook, api)
				if err != nil {
					cmd.Printf("app %q ExecuteHookCleanUp() error: %v", p.Path, err)
				}
			}
			return err
		}

		time.Sleep(100 * time.Millisecond)
		return nil
	}

	postCmd := cmd.PostRunE
	cmd.PostRunE = func(cmd *cobra.Command, args []string) error {
		api, err := newAppClientAPI(cmd)
		if err != nil {
			return err
		}

		ctx := cmd.Context()
		execHook := newExecutedHook(hook, cmd, args)

		defer func() {
			err := p.Interface.ExecuteHookCleanUp(ctx, execHook, api)
			if err != nil {
				cmd.Printf("app %q ExecuteHookCleanUp() error: %v", p.Path, err)
			}
		}()

		if postCmd != nil {
			err := postCmd(cmd, args)
			if err != nil {
				// dont return the error, log it and let execution continue to `Run`
				return err
			}
		}

		err = p.Interface.ExecuteHookPost(ctx, execHook, api)
		if err != nil {
			return errors.Errorf("app %q ExecuteHookPost() error : %w", p.Path, err)
		}
		return nil
	}
}

// linkPluginCmds tries to add the plugin commands to the legacy ignite
// commands.
func linkPluginCmds(rootCmd *cobra.Command, p *plugin.Plugin, pluginCmds []*plugin.Command) {
	if p.Error != nil {
		return
	}
	for _, pluginCmd := range pluginCmds {
		linkPluginCmd(rootCmd, p, pluginCmd)
		if p.Error != nil {
			return
		}
	}
}

func linkPluginCmd(rootCmd *cobra.Command, p *plugin.Plugin, pluginCmd *plugin.Command) {
	cmdPath := pluginCmd.Path()
	cmd := findCommandByPath(rootCmd, cmdPath)
	if cmd == nil {
		p.Error = errors.Errorf("unable to find command path %q for app %q", cmdPath, p.Path)
		return
	}
	if cmd.Runnable() {
		p.Error = errors.Errorf("can't attach app command %q to runnable command %q", pluginCmd.Use, cmd.CommandPath())
		return
	}

	// Check for existing commands
	// pluginCmd.Use can be like `command [args]` so we need to remove those
	// extra args if any.
	pluginCmdName := strings.Split(pluginCmd.Use, " ")[0]
	for _, cmd := range cmd.Commands() {
		if cmd.Name() == pluginCmdName {
			p.Error = errors.Errorf("app command %q already exists in Ignite's commands", pluginCmdName)
			return
		}
	}

	newCmd, err := pluginCmd.ToCobraCommand()
	if err != nil {
		p.Error = err
		return
	}
	cmd.AddCommand(newCmd)

	if len(pluginCmd.Commands) == 0 {
		// pluginCmd has no sub commands, so it's runnable
		newCmd.RunE = func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			return clictx.Do(ctx, func() error {
				api, err := newAppClientAPI(cmd)
				if err != nil {
					return err
				}

				// Call the plugin Execute
				execCmd := &plugin.ExecutedCommand{
					Use:    cmd.Use,
					Path:   cmd.CommandPath(),
					Args:   args,
					OsArgs: os.Args,
					With:   p.With,
				}
				execCmd.ImportFlags(cmd)
				err = p.Interface.Execute(ctx, execCmd, api)

				return err
			})
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

// NewApp returns a command that groups Ignite App related sub commands.
func NewApp() *cobra.Command {
	c := &cobra.Command{
		Use:   "app [command]",
		Short: "Create and manage Ignite Apps",
	}

	c.AddCommand(
		NewAppList(),
		NewAppUpdate(),
		NewAppScaffold(),
		NewAppDescribe(),
		NewAppInstall(),
		NewAppUninstall(),
	)

	return c
}

func NewAppList() *cobra.Command {
	lstCmd := &cobra.Command{
		Use:   "list",
		Short: "List installed apps",
		Long:  "Prints status and information of all installed Ignite Apps.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			s := cliui.New(cliui.WithStdout(os.Stdout))
			return printPlugins(cmd.Context(), s)
		},
	}
	return lstCmd
}

func NewAppUpdate() *cobra.Command {
	return &cobra.Command{
		Use:   "update [path]",
		Short: "Update app",
		Long: `Updates an Ignite App specified by path.

If no path is specified all declared apps are updated.`,
		Example: "ignite app update github.com/org/my-app/",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 {
				// update all plugins
				return plugin.Update(plugins...)
			}
			pluginPath, err := getAppPath(args[0])
			if err != nil {
				return err
			}

			// find the plugin to update
			for _, p := range plugins {
				if p.HasPath(pluginPath) {
					return plugin.Update(p)
				}
			}
			return errors.Errorf("App %q not found", pluginPath)
		},
	}
}

func NewAppInstall() *cobra.Command {
	cmdPluginAdd := &cobra.Command{
		Use:   "install [path] [key=value]...",
		Short: "Install app",
		Long: `Installs an Ignite App.

Respects key value pairs declared after the app path to be added to the generated configuration definition.`,
		Example: "ignite app install github.com/org/my-app/ foo=bar baz=qux",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			session := cliui.New(
				cliui.WithStdout(os.Stdout),
				cliui.WithoutUserInteraction(getYes(cmd)),
			)
			defer session.End()

			var (
				conf *pluginsconfig.Config
				err  error
			)

			global := flagGetPluginsGlobal(cmd)
			if global {
				conf, err = parseGlobalPlugins()
			} else {
				conf, err = parseLocalPlugins()
			}
			if err != nil {
				return err
			}

			pluginPath, err := getAppPath(args[0])
			if err != nil {
				return err
			}

			for _, p := range conf.Apps {
				if p.HasPath(pluginPath) {
					return errors.Errorf("app %s is already installed", pluginPath)
				}
			}

			p := pluginsconfig.Plugin{
				Path:   pluginPath,
				With:   make(map[string]string),
				Global: global,
			}

			pluginsOptions := []plugin.Option{
				plugin.CollectEvents(session.EventBus()),
			}

			var pluginArgs []string
			if len(args) > 1 {
				pluginArgs = args[1:]
			}

			for _, pa := range pluginArgs {
				kv := strings.Split(pa, "=")
				if len(kv) != 2 {
					return errors.Errorf("malformed key=value arg: %s", pa)
				}
				p.With[kv[0]] = kv[1]
			}

			plugins, err := plugin.Load(cmd.Context(), []pluginsconfig.Plugin{p}, pluginsOptions...)
			if err != nil {
				return err
			}
			defer plugins[0].KillClient()

			if err := plugins[0].Error; err != nil {
				if strings.Contains(err.Error(), "go.mod file not found in current directory") {
					return errors.Errorf("unable to find an App at the root of this repository (%s). Please ensure your repository URL is correct. If you're trying to install an App under a subfolder, include the path at the end of your repository URL, e.g., github.com/ignite/apps/appregistry", pluginPath)
				}

				return errors.Errorf("error while loading app %q: %w", pluginPath, plugins[0].Error)
			}
			session.Println(icons.OK, "Done loading apps")
			conf.Apps = append(conf.Apps, p)

			if err := conf.Save(); err != nil {
				return err
			}

			session.Printf("%s Installed %s\n", icons.Tada, pluginPath)
			return nil
		},
	}

	cmdPluginAdd.Flags().AddFlagSet(flagSetPluginsGlobal())

	return cmdPluginAdd
}

func NewAppUninstall() *cobra.Command {
	cmdPluginRemove := &cobra.Command{
		Use:     "uninstall [path]",
		Aliases: []string{"rm"},
		Short:   "Uninstall app",
		Long:    "Uninstalls an Ignite App specified by path.",
		Example: "ignite app uninstall github.com/org/my-app/",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s := cliui.New(cliui.WithStdout(os.Stdout))

			var (
				conf *pluginsconfig.Config
				err  error
			)

			global := flagGetPluginsGlobal(cmd)
			if global {
				conf, err = parseGlobalPlugins()
			} else {
				conf, err = parseLocalPlugins()
			}
			if err != nil {
				return err
			}

			pluginPath, err := getAppPath(args[0])
			if err != nil {
				return err
			}

			removed := false
			for i, cp := range conf.Apps {
				if cp.HasPath(pluginPath) {
					conf.Apps = append(conf.Apps[:i], conf.Apps[i+1:]...)
					removed = true
					break
				}
			}

			if !removed {
				// return if no matching plugin path found
				return errors.Errorf("app %s not found", pluginPath)
			}

			if err := conf.Save(); err != nil {
				return err
			}

			s.Printf("%s %s uninstalled\n", icons.OK, pluginPath)
			s.Printf("\t%s updated\n", conf.Path())

			return nil
		},
	}

	cmdPluginRemove.Flags().AddFlagSet(flagSetPluginsGlobal())

	return cmdPluginRemove
}

func NewAppScaffold() *cobra.Command {
	return &cobra.Command{
		Use:   "scaffold [name]",
		Short: "Scaffold a new Ignite App",
		Long: `Scaffolds a new Ignite App in the current directory.

A git repository will be created with the given module name, unless the current directory is already a git repository.`,
		Example: "ignite app scaffold github.com/org/my-app/",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			session := cliui.New(
				cliui.StartSpinnerWithText(statusScaffolding),
				cliui.WithoutUserInteraction(getYes(cmd)),
			)
			defer session.End()

			wd, err := os.Getwd()
			if err != nil {
				return err
			}
			moduleName := args[0]
			path, err := plugin.Scaffold(cmd.Context(), session, wd, moduleName, false)
			if err != nil {
				return err
			}
			if err := xgit.InitAndCommit(path); err != nil {
				return err
			}

			message := `⭐️ Successfully created a new Ignite App '%[1]s'.

👉 Update app code at '%[2]s/main.go'

👉 Test Ignite App integration by installing the app within the chain directory:

  ignite app install %[2]s

Or globally:

  ignite app install -g %[2]s

👉 Once the app is pushed to a repository, replace the local path by the repository path.
`
			session.Printf(message, moduleName, path)
			return nil
		},
	}
}

func NewAppDescribe() *cobra.Command {
	return &cobra.Command{
		Use:     "describe [path]",
		Short:   "Print information about installed apps",
		Long:    "Print information about an installed Ignite App commands and hooks.",
		Example: "ignite app describe github.com/org/my-app/",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				s   = cliui.New(cliui.WithStdout(os.Stdout))
				ctx = cmd.Context()
			)

			pluginPath, err := getAppPath(args[0])
			if err != nil {
				return err
			}

			for _, p := range plugins {
				if p.HasPath(pluginPath) {
					manifest, err := p.Interface.Manifest(ctx)
					if err != nil {
						return errors.Errorf("error while loading app manifest: %w", err)
					}

					if len(manifest.Commands) > 0 {
						s.Println("Commands:")
						for i, c := range manifest.Commands {
							cmdPath := fmt.Sprintf("%s %s", c.Path(), c.Use)
							s.Printf("  %d) %s\n", i+1, cmdPath)
						}
					}

					if len(manifest.Hooks) > 0 {
						s.Println("Hooks:")
						for i, h := range manifest.Hooks {
							s.Printf("  %d) '%s' on command '%s'\n", i+1, h.Name, h.CommandPath())
						}
					}

					break
				}
			}

			return nil
		},
	}
}

func getPluginLocationName(p *plugin.Plugin) string {
	if p.IsGlobal() {
		return "global"
	}
	return "local"
}

func getPluginStatus(ctx context.Context, p *plugin.Plugin) string {
	if p.Error != nil {
		return fmt.Sprintf("%s Error: %v", icons.NotOK, p.Error)
	}

	_, err := p.Interface.Manifest(ctx)
	if err != nil {
		return fmt.Sprintf("%s Error: Manifest() returned %v", icons.NotOK, err)
	}

	return fmt.Sprintf("%s Loaded", icons.OK)
}

func printPlugins(ctx context.Context, session *cliui.Session) error {
	var entries [][]string
	for _, p := range plugins {
		entries = append(entries, []string{p.Path, getPluginLocationName(p), getPluginStatus(ctx, p)})
	}

	if err := session.PrintTable([]string{"Path", "Config", "Status"}, entries...); err != nil {
		return errors.Errorf("error while printing apps: %w", err)
	}
	return nil
}

func newAppClientAPI(cmd *cobra.Command) (plugin.ClientAPI, error) {
	// Get chain when the plugin runs inside an blockchain app
	c, err := chain.NewWithHomeFlags(cmd)
	if err != nil && !errors.Is(err, gomodule.ErrGoModNotFound) {
		return nil, err
	}

	var options []plugin.APIOption
	if c != nil {
		options = append(options, plugin.WithChain(c))
	}

	return plugin.NewClientAPI(options...), nil
}

func flagSetPluginsGlobal() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.BoolP(flagPluginsGlobal, "g", false, "use global plugins configuration ($HOME/.ignite/apps/igniteapps.yml)")
	return fs
}

func flagGetPluginsGlobal(cmd *cobra.Command) bool {
	global, _ := cmd.Flags().GetBool(flagPluginsGlobal)
	return global
}

func getAppPath(path string) (string, error) {
	if xfilepath.IsDir(path) {
		// if directory is relative, make it absolute
		pluginPathAbs, err := xfilepath.MustAbs(path)
		if err != nil {
			return "", errors.Wrapf(err, "failed to get absolute path of %s", path)
		}
		path = pluginPathAbs
	}
	return path, nil
}
