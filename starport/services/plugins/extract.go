package plugins

func (m *Manager) extractPlugins(ctx context.Context, cfg chaincfg.Config) error {
	cmdPlugins, err := extractCommandPlugins(ctx, m.ChainId, cfg)
	if err != nil {
		return err
	}

	hookPlugins, err := extractHookPlugins(ctx, m.ChainId, cfg)
	if err != nil {
		return err
	}

	m.cmdPlugins = cmdPlugins
	m.hookPlugins = hookPlugins
	return nil
}

func extractCommandPlugins(ctx context.Context, chainId string, cfg chaincfg.Config) ([]CmdPlugin, error) {
	pluginsDir, err := formatPluginHome(chainId, "")
	if err != nil {
		return nil, err
	}

	pluginFiles, err := listFiles(pluginsDir, `.*_cmd.so`)
	if err != nil {
		return nil, err
	}

	var cmdPlugins []CmdPlugin
	for _, pluginFile := range pluginFiles {
		pluginDir := path.Join(pluginsDir, pluginFile.Name())
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

		if err := validateParentCommand(cmdPlugin.ParentCommand()); err != nil {
			return nil, err
		}

		cmdPlugins = append(cmdPlugins, cmdPlugin)
	}

	return cmdPlugins, nil
}

func extractHookPlugins(ctx context.Context, chainId string, cfg chaincfg.Config) ([]HookPlugin, error) {
	pluginsDir, err := formatPluginHome(chainId, "")
	if err != nil {
		return nil, err
	}

	pluginFiles, err := listFiles(pluginsDir, `.*_hook.so`)
	if err != nil {
		return nil, err
	}

	var hookPlugins []HookPlugin
	for _, pluginFile := range pluginFiles {
		pluginDir := path.Join(pluginsDir, pluginFile.Name())
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

		if err := validateParentCommand(cmdPlugin.ParentCommand()); err != nil {
			return nil, err
		}

		hookPlugins = append(hookPlugins, hookPlugin)
	}

	return hookPlugins, nil
}
