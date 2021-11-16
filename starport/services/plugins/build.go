package plugins

func (m *Manager) build(ctx context.Context, cfg chaincfg.Config) error {
	outputDir := path.Join(dst, "output")

	// Get plugin home
	dst, err := formatPluginHome(m.ChainId, "")
	if err != nil {
		return err
	}

	// Enter plugins directory, go get .
	// Somehow have to account for remote dependencies in individual plugins
	pluginDirs, err := listDirs(dst)
	if err != nil {
		return err
	}

	for _, pluginSubDir := range pluginDirs {
		pluginDir := path.Join(dst, pluginSubDir.Name())

		cmdPlugins, err := listFiles(pluginDir, "*.cmd.go")
		if err != nil {
			return err
		}

		if len(cmdPlugins) > 0 {
			for _, pluginFile := range cmdPlugins {
				fileName := pluginFile.Name()
				outputName := strings.Trim(fileName, ".cmd.go")
				inputFileDir := path.Join(pluginDir, fileName)
				outputFileDir := path.Join(outputDir, outputName+"_cmd.so")
				buildPlugin(ctx, outputFileDir, inputFileDir, []string{})
			}
		}

		hookPlugins, err := listFiles(pluginDir, "*.hook.go")
		if err != nil {
			return err
		}

		if len(hookPlugins) > 0 {
			for _, pluginFile := range hookPlugins {
				fileName := pluginFile.Name()
				outputName := strings.Trim(fileName, ".hook.go")
				inputFileDir := path.Join(pluginDir, fileName)
				outputFileDir := path.Join(outputDir, outputName+"_hook.so")
				buildPlugin(ctx, outputFileDir, inputFileDir, []string{})
			}
		}
	}

	return nil
}

func buildPlugin(ctx context.Context, output string, path string, flags []string) error {
	command := []string{
		"go",
		gocmd.CommandBuild,
		FLAG_BUILD_MODE_PLUGIN,
		gocmd.FlagOut,
		output,
		path,
	}

	command = append(command, flags...)
	return exec.Exec(ctx, command)
}
