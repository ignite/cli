package plugins

import (
	"context"
	"path"
	"plugin"

	"github.com/spf13/cobra"
	chaincfg "github.com/tendermint/starport/starport/chainconfig"
)

func (m *Manager) extractPlugins(ctx context.Context, parentCommand *cobra.Command, cfg chaincfg.Config) error {
	cmdPlugins, err := extractCommandPlugins(ctx, m.ChainId, parentCommand, cfg)
	if err != nil {
		return err
	}

	hookPlugins, err := extractHookPlugins(ctx, m.ChainId, parentCommand, cfg)
	if err != nil {
		return err
	}

	m.cmdPlugins = cmdPlugins
	m.hookPlugins = hookPlugins
	return nil
}

func extractCommandPlugins(
	ctx context.Context,
	pluginDir string,
	parentCommand *cobra.Command,
	cfg chaincfg.Config,
) ([]CmdPlugin, error) {
	pluginFiles, err := listFiles(pluginDir, `.*_cmd.so`)
	if err != nil {
		return nil, err
	}

	var cmdPlugins []CmdPlugin
	for _, pluginFile := range pluginFiles {
		pluginDir := path.Join(pluginDir, pluginFile.Name())
		plug, err := plugin.Open(pluginDir)
		if err != nil {
			return nil, err
		}

		symCmdPlugin, err := plug.Lookup(CMD_SYMBOL_NAME)
		if err != nil {
			return nil, err
		}

		var cmdPlugin CmdPlugin
		cmdPlugin, ok := symCmdPlugin.(CmdPlugin)
		if !ok {
			return nil, ErrCommandPluginNotRecognized
		}

		if err := cmdPlugin.Init(ctx); err != nil {
			return nil, err
		}

		for _, command := range cmdPlugin.Registry() {
			if err := validateParentCommand(parentCommand, command.ParentCommand()); err != nil {
				return nil, err
			}
		}

		cmdPlugins = append(cmdPlugins, cmdPlugin)
	}

	return cmdPlugins, nil
}

func extractHookPlugins(
	ctx context.Context,
	pluginDir string,
	parentCommand *cobra.Command,
	cfg chaincfg.Config,
) ([]HookPlugin, error) {
	pluginFiles, err := listFiles(pluginDir, `.*_hook.so`)
	if err != nil {
		return nil, err
	}

	var hookPlugins []HookPlugin
	for _, pluginFile := range pluginFiles {
		pluginDir := path.Join(pluginDir, pluginFile.Name())
		plug, err := plugin.Open(pluginDir)
		if err != nil {
			return nil, err
		}

		symCmdPlugin, err := plug.Lookup(HOOK_PLUGIN_NAME)
		if err != nil {
			return nil, err
		}

		var hookPlugin HookPlugin
		hookPlugin, ok := symCmdPlugin.(HookPlugin)
		if !ok {
			return nil, ErrCommandPluginNotRecognized
		}

		for _, command := range hookPlugin.Registry() {
			if err := validateParentCommand(parentCommand, command.ParentCommand()); err != nil {
				return nil, err
			}
		}

		hookPlugins = append(hookPlugins, hookPlugin)
	}

	return hookPlugins, nil
}
