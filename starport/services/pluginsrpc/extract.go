package pluginsrpc

import (
	"context"
	"os/exec"
	"path"

	"github.com/hashicorp/go-plugin"
	"github.com/spf13/cobra"
	plugintypes "github.com/tendermint/starport/starport/services/pluginsrpc/types"
)

func (m *Manager) extractPlugins(ctx context.Context, rootCmd *cobra.Command, args []string) error {
	outputDir, err := formatPluginHome(m.ChainId, "output")
	if err != nil {
		return err
	}

	for i := 0; i < len(m.Config.Plugins); i++ {
		err := m.extractCommandPlugins(ctx, outputDir, rootCmd)
		if err != nil {
			return err
		}

		err = m.extractHookPlugins(ctx, outputDir, rootCmd)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) extractCommandPlugins(
	ctx context.Context,
	outputDir string,
	parentCommand *cobra.Command,
) error {
	pluginFiles, err := listDirsMatch(outputDir, `*_cmd`)
	if err != nil {
		return err
	}

	if len(pluginFiles) == 0 {
		return nil
	}

	// Remove pluginFiles that are not specified in the config

	var extractedCommandModules []ExtractedCommandModule
	for _, pluginFile := range pluginFiles {
		pluginDir := path.Join(outputDir, pluginFile.Name())
		PluginMap := BasePluginMap

		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: HandshakeConfig,
			Plugins:         PluginMap,
			Cmd:             exec.Command(pluginDir),
			Logger:          pluginLogger,
		})

		rpcClient, err := client.Client()
		if err != nil {
			return err
		}

		raw, err := rpcClient.Dispense("command_map")
		if err != nil {
			return err
		}

		cmdMapper := raw.(plugintypes.CommandMapper)
		storedCommands := cmdMapper.Commands()
		client.Kill()

		for _, loadedModule := range storedCommands {
			NewPluginMap := map[string]plugin.Plugin{
				loadedModule: &plugintypes.CommandModulePlugin{},
			}

			client2 := plugin.NewClient(&plugin.ClientConfig{
				HandshakeConfig: HandshakeConfig,
				Plugins:         NewPluginMap,
				Cmd:             exec.Command(pluginDir),
				Logger:          pluginLogger,
			})

			rpcClient2, err := client2.Client()
			if err != nil {
				return err
			}

			raw2, err := rpcClient2.Dispense(loadedModule)
			if err != nil {
				return err
			}

			cmdModule, ok := raw2.(plugintypes.CommandModule)
			if !ok {
				return ErrCommandPluginNotRecognized
			}

			extractedCommandModules = append(extractedCommandModules, ExtractedCommandModule{
				ModuleName:    loadedModule,
				PluginDir:     pluginDir,
				ParentCommand: cmdModule.GetParentCommand(),
				Name:          cmdModule.GetName(),
				Usage:         cmdModule.GetUsage(),
				ShortDesc:     cmdModule.GetShortDesc(),
				LongDesc:      cmdModule.GetLongDesc(),
				NumArgs:       cmdModule.GetNumArgs(),
				Exec: func(cmd *cobra.Command, args []string) error {
					client := plugin.NewClient(&plugin.ClientConfig{
						HandshakeConfig: HandshakeConfig,
						Plugins:         NewPluginMap,
						Cmd:             exec.Command(pluginDir),
						Logger:          pluginLogger,
					})

					rpcClient, err := client.Client()
					if err != nil {
						return err
					}

					raw, err := rpcClient.Dispense(loadedModule)
					if err != nil {
						return err
					}

					cmdModuleExec := raw.(plugintypes.CommandModule)
					err = cmdModuleExec.Exec(cmd, args)
					if err != nil {
						return err
					}

					client.Kill()
					return nil
				},
			})

			client2.Kill()
		}
	}

	m.cmdPlugins = append(m.cmdPlugins, extractedCommandModules...)
	return nil
}

func (m *Manager) extractHookPlugins(
	ctx context.Context,
	outputDir string,
	parentCommand *cobra.Command,
) error {
	pluginFiles, err := listDirsMatch(outputDir, `*_hook`)
	if err != nil {
		return err
	}

	if len(pluginFiles) == 0 {
		return nil
	}

	// Remove pluginFiles that are not specified in the config

	var extractedHookModules []ExtractedHookModule
	for _, pluginFile := range pluginFiles {
		pluginDir := path.Join(outputDir, pluginFile.Name())
		PluginMap := BasePluginMap

		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: HandshakeConfig,
			Plugins:         PluginMap,
			Cmd:             exec.Command(pluginDir),
			Logger:          pluginLogger,
		})

		rpcClient, err := client.Client()
		if err != nil {
			return err
		}

		raw, err := rpcClient.Dispense("hook_map")
		if err != nil {
			return err
		}

		hookMapper := raw.(plugintypes.HookMapper)

		storedHooks := hookMapper.Hooks()

		client.Kill()

		// Edit pluginMap off of that, and then make execute functions for each thing
		for _, loadedModule := range storedHooks {
			NewPluginMap := map[string]plugin.Plugin{
				loadedModule: &plugintypes.HookModulePlugin{},
			}

			client2 := plugin.NewClient(&plugin.ClientConfig{
				HandshakeConfig: HandshakeConfig,
				Plugins:         NewPluginMap,
				Cmd:             exec.Command(pluginDir),
				Logger:          pluginLogger,
			})

			rpcClient2, err := client2.Client()
			if err != nil {
				return err
			}

			raw2, err := rpcClient2.Dispense(loadedModule)
			if err != nil {
				return err
			}

			// stops here on cast, making me think something is wrong with plugintypes
			hookModule, ok := raw2.(plugintypes.HookModule)
			if !ok {
				return ErrHookPluginNotRecognized
			}

			extractedHookModules = append(extractedHookModules, ExtractedHookModule{
				ModuleName:    loadedModule,
				PluginDir:     pluginDir,
				ParentCommand: hookModule.GetParentCommand(),
				Name:          hookModule.GetName(),
				HookType:      hookModule.GetType(),
				PreRun: func(cmd *cobra.Command, args []string) error {
					client := plugin.NewClient(&plugin.ClientConfig{
						HandshakeConfig: HandshakeConfig,
						Plugins:         NewPluginMap,
						Cmd:             exec.Command(pluginDir),
						Logger:          pluginLogger,
					})

					rpcClient, err := client.Client()
					if err != nil {
						return err
					}

					raw, err := rpcClient.Dispense(loadedModule)
					if err != nil {
						return err
					}

					hookModuleExec := raw.(plugintypes.HookModule)
					err = hookModuleExec.PreRun(cmd, args)
					if err != nil {
						return err
					}

					client.Kill()
					return nil
				},
				PostRun: func(cmd *cobra.Command, args []string) error {
					client := plugin.NewClient(&plugin.ClientConfig{
						HandshakeConfig: HandshakeConfig,
						Plugins:         NewPluginMap,
						Cmd:             exec.Command(pluginDir),
						Logger:          pluginLogger,
					})

					rpcClient, err := client.Client()
					if err != nil {
						return err
					}

					raw, err := rpcClient.Dispense(loadedModule)
					if err != nil {
						return err
					}

					hookModuleExec := raw.(plugintypes.HookModule)
					err = hookModuleExec.PostRun(cmd, args)
					if err != nil {
						return err
					}

					client.Kill()
					return nil
				},
			})

			client2.Kill()
		}
	}

	m.hookPlugins = append(m.hookPlugins, extractedHookModules...)
	return nil
}
