package plugins

// import (
// 	"context"
// 	"path"
// 	"plugin"

// 	"github.com/spf13/cobra"
// 	chaincfg "github.com/tendermint/starport/starport/chainconfig"
// )

// func (m *Manager) extractPlugins(ctx context.Context, rootCmd *cobra.Command) error {
// 	pluginHome, err := formatPluginHome(m.ChainId, "")
// 	if err != nil {
// 		return err
// 	}

// 	for _, cfgPlugin := range m.Config.Plugins {
// 		pluginId := getPluginId(cfgPlugin)
// 		pluginDir := path.Join(pluginHome, pluginId)
// 		outputDir := path.Join(pluginHome, "output")

// 		cmdPlugins, err := extractCommandPlugins(ctx, pluginDir, outputDir, rootCmd, m.Config)
// 		if err != nil {
// 			return err
// 		}

// 		hookPlugins, err := extractHookPlugins(ctx, pluginDir, outputDir, rootCmd, m.Config)
// 		if err != nil {
// 			return err
// 		}

// 		m.cmdPlugins = append(m.cmdPlugins, cmdPlugins...)
// 		m.hookPlugins = append(m.hookPlugins, hookPlugins...)
// 	}

// 	return nil
// }

// func extractCommandPlugins(
// 	ctx context.Context,
// 	pluginDir string,
// 	outputDir string,
// 	parentCommand *cobra.Command,
// 	cfg chaincfg.Config,
// ) ([]CmdPlugin, error) {
// 	pluginFiles, err := listFilesMatch(outputDir, `*_cmd.so`)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if len(pluginFiles) == 0 {
// 		return []CmdPlugin{}, nil
// 	}

// 	var cmdPlugins []CmdPlugin
// 	for _, pluginFile := range pluginFiles {
// 		pluginDir := path.Join(outputDir, pluginFile.Name())
// 		plug, err := plugin.Open(pluginDir)
// 		if err != nil {
// 			return nil, err
// 		}

// 		symCmdPlugin, err := plug.Lookup(CMD_SYMBOL_NAME)
// 		if err != nil {
// 			return nil, err
// 		}

// 		cmdPlugin, ok := symCmdPlugin.(CmdPlugin)
// 		if !ok {
// 			return nil, ErrCommandPluginNotRecognized
// 		}

// 		for _, command := range cmdPlugin.Registry() {
// 			if err := validateParentCommand(parentCommand, command.ParentCommand()); err != nil {
// 				return nil, err
// 			}
// 		}

// 		cmdPlugins = append(cmdPlugins, cmdPlugin)
// 	}

// 	return cmdPlugins, nil
// }

// func extractHookPlugins(
// 	ctx context.Context,
// 	pluginDir string,
// 	outputDir string,
// 	parentCommand *cobra.Command,
// 	cfg chaincfg.Config,
// ) ([]HookPlugin, error) {
// 	pluginFiles, err := listFilesMatch(outputDir, `.*_hook.so`)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if len(pluginFiles) == 0 {
// 		return []HookPlugin{}, nil
// 	}

// 	var hookPlugins []HookPlugin
// 	for _, pluginFile := range pluginFiles {
// 		pluginDir := path.Join(outputDir, pluginFile.Name())
// 		plug, err := plugin.Open(pluginDir)
// 		if err != nil {
// 			return nil, err
// 		}

// 		symCmdPlugin, err := plug.Lookup(HOOK_PLUGIN_NAME)
// 		if err != nil {
// 			return nil, err
// 		}

// 		var hookPlugin HookPlugin
// 		hookPlugin, ok := symCmdPlugin.(HookPlugin)
// 		if !ok {
// 			return nil, ErrCommandPluginNotRecognized
// 		}

// 		for _, command := range hookPlugin.Registry() {
// 			if err := validateParentCommand(parentCommand, command.ParentCommand()); err != nil {
// 				return nil, err
// 			}
// 		}

// 		hookPlugins = append(hookPlugins, hookPlugin)
// 	}

// 	return hookPlugins, nil
// }
