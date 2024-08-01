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
	"github.com/ignite/cli/v29/ignite/pkg/xgit"
	"github.com/ignite/cli/v29/ignite/services/chain"
	"github.com/ignite/cli/v29/ignite/services/plugin"
)

const (
	flagExtensionsGlobal = "global"
)

// extensions hold the list of extensions declared in the config.
// A global variable is used so the list is accessible to the extension commands.
var extensions []*plugin.Plugin

// LoadExtensions tries to load all the extensions found in configurations.
// If no configurations found, it returns w/o error.
func LoadExtensions(ctx context.Context, cmd *cobra.Command, session *cliui.Session) error {
	var (
		rootCmd           = cmd.Root()
		extensionsConfigs []pluginsconfig.Plugin
	)
	localCfg, err := parseLocalExtensions(rootCmd)
	if err != nil && !errors.As(err, &cosmosanalysis.ErrPathNotChain{}) {
		return err
	} else if err == nil {
		extensionsConfigs = append(extensionsConfigs, localCfg.Extensions...)
	}

	globalCfg, err := parseGlobalExtensions()
	if err == nil {
		extensionsConfigs = append(extensionsConfigs, globalCfg.Extensions...)
	}
	ensureDefaultExtensions(cmd, globalCfg)

	if len(extensionsConfigs) == 0 {
		return nil
	}

	uniqueExtensions := pluginsconfig.RemoveDuplicates(extensionsConfigs)
	extensions, err = plugin.Load(ctx, uniqueExtensions, plugin.CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}
	if len(extensions) == 0 {
		return nil
	}

	return linkExtensions(ctx, rootCmd, extensions)
}

func parseLocalExtensions(cmd *cobra.Command) (*pluginsconfig.Config, error) {
	// FIXME(tb): like other commands that works on a chain directory,
	// parseLocalExtensions should rely on `-p` flag to guess that chain directory.
	// Unfortunately parseLocalExtensions is invoked before flags are parsed, so
	// we cannot rely on `-p` flag. As a workaround, we use the working dir.
	// The drawback is we cannot load chain's plugin when using `-p`.
	_ = cmd
	wd, err := os.Getwd()
	if err != nil {
		return nil, errors.Errorf("parse local extensions: %w", err)
	}
	if err := cosmosanalysis.IsChainPath(wd); err != nil {
		return nil, err
	}
	return pluginsconfig.ParseDir(wd)
}

func parseGlobalExtensions() (cfg *pluginsconfig.Config, err error) {
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

	for i := range cfg.Extensions {
		cfg.Extensions[i].Global = true
	}
	return
}

func linkExtensions(ctx context.Context, rootCmd *cobra.Command, plugins []*plugin.Plugin) error {
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

		linkExtensionHooks(rootCmd, p, manifest.Hooks)
		if p.Error != nil {
			linkErrors = append(linkErrors, p)
			continue
		}

		linkExtensionCmds(rootCmd, p, manifest.Commands)
		if p.Error != nil {
			linkErrors = append(linkErrors, p)
			continue
		}
	}

	if len(linkErrors) > 0 {
		// unload any plugin that could have been loaded
		defer UnloadExtensions()

		if err := printExtensions(ctx, cliui.New(cliui.WithStdout(os.Stdout))); err != nil {
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

// UnloadExtensions releases any loaded extensions, which is basically killing the
// plugin server instance.
func UnloadExtensions() {
	for _, p := range extensions {
		p.KillClient()
	}
}

func linkExtensionHooks(rootCmd *cobra.Command, p *plugin.Plugin, hooks []*plugin.Hook) {
	if p.Error != nil {
		return
	}
	for _, hook := range hooks {
		linkExtensionHook(rootCmd, p, hook)
	}
}

func linkExtensionHook(rootCmd *cobra.Command, p *plugin.Plugin, hook *plugin.Hook) {
	cmdPath := hook.CommandPath()
	cmd := findCommandByPath(rootCmd, cmdPath)
	if cmd == nil {
		p.Error = errors.Errorf("unable to find command path %q for extension hook %q", cmdPath, hook.Name)
		return
	}
	if !cmd.Runnable() {
		p.Error = errors.Errorf("can't attach extension hook %q to non executable command %q", hook.Name, hook.PlaceHookOn)
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

		api, err := newExtensionClientAPI(cmd)
		if err != nil {
			return err
		}

		ctx := cmd.Context()
		execHook := newExecutedHook(hook, cmd, args)
		err = p.Interface.ExecuteHookPre(ctx, execHook, api)
		if err != nil {
			return errors.Errorf("extension %q ExecuteHookPre() error: %w", p.Path, err)
		}
		return nil
	}

	runCmd := cmd.RunE
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if runCmd != nil {
			err := runCmd(cmd, args)
			// if the command has failed the `PostRun` will not execute. here we execute the cleanup step before returning.
			if err != nil {
				api, err := newExtensionClientAPI(cmd)
				if err != nil {
					return err
				}

				ctx := cmd.Context()
				execHook := newExecutedHook(hook, cmd, args)
				err = p.Interface.ExecuteHookCleanUp(ctx, execHook, api)
				if err != nil {
					cmd.Printf("extension %q ExecuteHookCleanUp() error: %v", p.Path, err)
				}
			}
			return err
		}

		time.Sleep(100 * time.Millisecond)
		return nil
	}

	postCmd := cmd.PostRunE
	cmd.PostRunE = func(cmd *cobra.Command, args []string) error {
		api, err := newExtensionClientAPI(cmd)
		if err != nil {
			return err
		}

		ctx := cmd.Context()
		execHook := newExecutedHook(hook, cmd, args)

		defer func() {
			err := p.Interface.ExecuteHookCleanUp(ctx, execHook, api)
			if err != nil {
				cmd.Printf("extension %q ExecuteHookCleanUp() error: %v", p.Path, err)
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
			return errors.Errorf("extension %q ExecuteHookPost() error : %w", p.Path, err)
		}
		return nil
	}
}

// linkExtensionCmds tries to add the plugin commands to the legacy ignite
// commands.
func linkExtensionCmds(rootCmd *cobra.Command, p *plugin.Plugin, pluginCmds []*plugin.Command) {
	if p.Error != nil {
		return
	}
	for _, pluginCmd := range pluginCmds {
		linkExtensionCmd(rootCmd, p, pluginCmd)
		if p.Error != nil {
			return
		}
	}
}

func linkExtensionCmd(rootCmd *cobra.Command, p *plugin.Plugin, pluginCmd *plugin.Command) {
	cmdPath := pluginCmd.Path()
	cmd := findCommandByPath(rootCmd, cmdPath)
	if cmd == nil {
		p.Error = errors.Errorf("unable to find command path %q for extension %q", cmdPath, p.Path)
		return
	}
	if cmd.Runnable() {
		p.Error = errors.Errorf("can't attach extension command %q to runnable command %q", pluginCmd.Use, cmd.CommandPath())
		return
	}

	// Check for existing commands
	// pluginCmd.Use can be like `command [args]` so we need to remove those
	// extra args if any.
	pluginCmdName := strings.Split(pluginCmd.Use, " ")[0]
	for _, cmd := range cmd.Commands() {
		if cmd.Name() == pluginCmdName {
			p.Error = errors.Errorf("extension command %q already exists in Ignite's commands", pluginCmdName)
			return
		}
	}

	newCmd, err := pluginCmd.ToCobraCommand()
	if err != nil {
		p.Error = err
		return
	}
	cmd.AddCommand(newCmd)

	// NOTE(tb) we could probably simplify by removing this condition and call the
	// plugin even if the invoked command isn't runnable. If we do so, the plugin
	// will be responsible for outputing the standard cobra output, which implies
	// it must use cobra too. This is how cli-plugin-network works, but to make
	// it for all, we need to change the `plugin scaffold` output (so it outputs
	// something similar than the cli-plugin-network) and update the docs.
	if len(pluginCmd.Commands) == 0 {
		// pluginCmd has no sub commands, so it's runnable
		newCmd.RunE = func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			return clictx.Do(ctx, func() error {
				api, err := newExtensionClientAPI(cmd)
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

				// NOTE(tb): This pause gives enough time for go-plugin to sync the
				// output from stdout/stderr of the plugin. Without that pause, this
				// output can be discarded and not printed in the user console.
				time.Sleep(100 * time.Millisecond)
				return err
			})
		}
	} else {
		for _, pluginCmd := range pluginCmd.Commands {
			pluginCmd.PlaceCommandUnder = newCmd.CommandPath()
			linkExtensionCmd(newCmd, p, pluginCmd)
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

// NewExtension returns a command that groups Ignite Extensions related sub commands.
func NewExtension() *cobra.Command {
	c := &cobra.Command{
		Use:     "extension [command]",
		Aliases: []string{"ext", "extn", "app", "plugin", "extensions"},
		Short:   "Create and manage Ignite Extensions",
	}

	c.AddCommand(
		NewExtensionList(),
		NewExtensionUpdate(),
		NewExtensionScaffold(),
		NewExtensionDescribe(),
		NewExtensionInstall(),
		NewExtensionUninstall(),
	)

	return c
}

func NewExtensionList() *cobra.Command {
	lstCmd := &cobra.Command{
		Use:   "list",
		Short: "List installed extensions",
		Long:  "Prints status and information of all installed Ignite Extensions.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			s := cliui.New(cliui.WithStdout(os.Stdout))
			return printExtensions(cmd.Context(), s)
		},
	}
	return lstCmd
}

func NewExtensionUpdate() *cobra.Command {
	return &cobra.Command{
		Use:   "update [path]",
		Short: "Updates an Ignite Extension",
		Long: `Updates an Ignite Extension specified by path.

If no path is specified all declared apps are updated.`,
		Example: "ignite extension update github.com/org/my-extension/",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 {
				// update all plugins
				return plugin.Update(extensions...)
			}
			// find the plugin to update
			for _, p := range extensions {
				if p.HasPath(args[0]) {
					return plugin.Update(p)
				}
			}
			return errors.Errorf("Extension %q not found", args[0])
		},
	}
}

func NewExtensionInstall() *cobra.Command {
	cmdExtensionInstall := &cobra.Command{
		Use:   "install [path] [key=value]...",
		Short: "Installs an Ingite Extension",
		Long: `Installs an Ignite Extension.

Respects key value pairs declared after the extension path to be added to the generated configuration definition.`,
		Example: "ignite extension install github.com/org/my-extension/ foo=bar baz=qux",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			session := cliui.New(cliui.WithStdout(os.Stdout))
			defer session.End()

			var (
				conf *pluginsconfig.Config
				err  error
			)

			global := flagGetExtensionsGlobal(cmd)
			if global {
				conf, err = parseGlobalExtensions()
			} else {
				conf, err = parseLocalExtensions(cmd)
			}
			if err != nil {
				return err
			}

			for _, p := range conf.Extensions {
				if p.HasPath(args[0]) {
					return errors.Errorf("extension %s is already installed", args[0])
				}
			}

			p := pluginsconfig.Plugin{
				Path:   args[0],
				With:   make(map[string]string),
				Global: global,
			}

			extensionsOptions := []plugin.Option{
				plugin.CollectEvents(session.EventBus()),
			}

			var extensionArgs []string
			if len(args) > 1 {
				extensionArgs = args[1:]
			}

			for _, pa := range extensionArgs {
				kv := strings.Split(pa, "=")
				if len(kv) != 2 {
					return errors.Errorf("malformed key=value arg: %s", pa)
				}
				p.With[kv[0]] = kv[1]
			}

			extensions, err := plugin.Load(cmd.Context(), []pluginsconfig.Plugin{p}, extensionsOptions...)
			if err != nil {
				return err
			}
			defer extensions[0].KillClient()

			if extensions[0].Error != nil {
				return errors.Errorf("error while loading extension %q: %w", args[0], extensions[0].Error)
			}
			session.Println(icons.OK, "Done loading extensions")
			conf.Extensions = append(conf.Extensions, p)

			if err := conf.Save(); err != nil {
				return err
			}

			session.Printf("%s Installed %s\n", icons.Tada, args[0])
			return nil
		},
	}

	cmdExtensionInstall.Flags().AddFlagSet(flagSetExtensionsGlobal())

	return cmdExtensionInstall
}

func NewExtensionUninstall() *cobra.Command {
	cmdExtensionUninstall := &cobra.Command{
		Use:     "uninstall [path]",
		Aliases: []string{"rm"},
		Short:   "Uninstall an Ignite Extension xtension",
		Long:    "Uninstalls an Ignite Extension specified by path.",
		Example: "ignite extension uninstall github.com/org/my-extension/",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s := cliui.New(cliui.WithStdout(os.Stdout))

			var (
				conf *pluginsconfig.Config
				err  error
			)

			global := flagGetExtensionsGlobal(cmd)
			if global {
				conf, err = parseGlobalExtensions()
			} else {
				conf, err = parseLocalExtensions(cmd)
			}
			if err != nil {
				return err
			}

			removed := false
			for i, cp := range conf.Extensions {
				if cp.HasPath(args[0]) {
					conf.Extensions = append(conf.Extensions[:i], conf.Extensions[i+1:]...)
					removed = true
					break
				}
			}

			if !removed {
				// return if no matching plugin path found
				return errors.Errorf("extension %s not found", args[0])
			}

			if err := conf.Save(); err != nil {
				return err
			}

			s.Printf("%s %s uninstalled\n", icons.OK, args[0])
			s.Printf("\t%s updated\n", conf.Path())

			return nil
		},
	}

	cmdExtensionUninstall.Flags().AddFlagSet(flagSetExtensionsGlobal())

	return cmdExtensionUninstall
}

func NewExtensionScaffold() *cobra.Command {
	return &cobra.Command{
		Use:   "scaffold [name]",
		Short: "Scaffold a new Ignite Extension",
		Long: `Scaffolds a new Ignite Extension in the current directory.

A git repository will be created with the given module name, unless the current directory is already a git repository.`,
		Example: "ignite extension scaffold github.com/org/my-extension/",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			session := cliui.New(cliui.StartSpinnerWithText(statusScaffolding))
			defer session.End()

			wd, err := os.Getwd()
			if err != nil {
				return err
			}
			moduleName := args[0]
			path, err := plugin.Scaffold(cmd.Context(), wd, moduleName, false)
			if err != nil {
				return err
			}
			if err := xgit.InitAndCommit(path); err != nil {
				return err
			}

			message := `⭐️ Successfully created a new Ignite Extension '%[1]s'.

👉 Update extension code at '%[2]s/main.go'

👉 Test Ignite Extension integration by installing the extension within the chain directory:

  ignite extension install %[2]s

Or globally:

  ignite extension install -g %[2]s

👉 Once the extension is pushed to a repository, replace the local path by the repository path.
`
			session.Printf(message, moduleName, path)
			return nil
		},
	}
}

func NewExtensionDescribe() *cobra.Command {
	return &cobra.Command{
		Use:     "describe [path]",
		Short:   "Print information about installed extensions",
		Long:    "Print information about an installed Ignite Extension commands and hooks.",
		Example: "ignite extension describe github.com/org/my-extension/",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s := cliui.New(cliui.WithStdout(os.Stdout))
			ctx := cmd.Context()

			for _, p := range extensions {
				if p.HasPath(args[0]) {
					manifest, err := p.Interface.Manifest(ctx)
					if err != nil {
						return errors.Errorf("error while loading extension manifest: %w", err)
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

func getExtensionLocationName(p *plugin.Plugin) string {
	if p.IsGlobal() {
		return "global"
	}
	return "local"
}

func getExtensionStatus(ctx context.Context, p *plugin.Plugin) string {
	if p.Error != nil {
		return fmt.Sprintf("%s Error: %v", icons.NotOK, p.Error)
	}

	_, err := p.Interface.Manifest(ctx)
	if err != nil {
		return fmt.Sprintf("%s Error: Manifest() returned %v", icons.NotOK, err)
	}

	return fmt.Sprintf("%s Loaded", icons.OK)
}

func printExtensions(ctx context.Context, session *cliui.Session) error {
	var entries [][]string
	for _, p := range extensions {
		entries = append(entries, []string{p.Path, getExtensionLocationName(p), getExtensionStatus(ctx, p)})
	}

	if err := session.PrintTable([]string{"Path", "Config", "Status"}, entries...); err != nil {
		return errors.Errorf("error while printing extensions: %w", err)
	}
	return nil
}

func newExtensionClientAPI(cmd *cobra.Command) (plugin.ClientAPI, error) {
	// Get chain when the extension runs inside a blockchain app
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

func flagSetExtensionsGlobal() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.BoolP(flagExtensionsGlobal, "g", false, "use global extensions configuration ($HOME/.ignite/extensions/extensions.yml)")
	return fs
}

func flagGetExtensionsGlobal(cmd *cobra.Command) bool {
	global, _ := cmd.Flags().GetBool(flagExtensionsGlobal)
	return global
}

// Backward compat.
var (
	LoadPlugins   = LoadExtensions
	UnloadPlugins = UnloadExtensions
	NewAppInstall = NewExtensionInstall
)
