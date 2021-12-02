package pluginsrpc

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/hashicorp/go-plugin"
	"github.com/spf13/cobra"
	plugintypes "github.com/tendermint/starport/starport/services/pluginsrpc/types"
)

// CachedCommand represents the information the cache stores on a command
type CachedCommand struct {
	ModuleName string
	PluginDir  string
	plugintypes.Command
}

// CachedHook represents the information the cache stores on a hook
type CachedHook struct {
	ModuleName string
	PluginDir  string
	plugintypes.Hook
}

// Cache caches the stored plugins
func (m *Manager) Cache(ctx context.Context) error {
	cacheHome, err := formatPluginHome(m.ChainId, "cached")
	if err != nil {
		return err
	}

	for _, cmdPlugin := range m.cmdPlugins {
		fileParts := strings.Split(cmdPlugin.PluginDir, "/")
		fileName := strings.Trim(fileParts[len(fileParts)-1], "_cmd")
		targetFile := path.Join(cacheHome, (fileName + ".cmd.json"))

		_, err = os.Stat(targetFile)
		if !errors.Is(err, os.ErrNotExist) {
			continue
		} else if err == nil {
			continue
		}

		cachedCommand := CachedCommand{
			ModuleName: cmdPlugin.ModuleName,
			PluginDir:  cmdPlugin.PluginDir,
			Command: plugintypes.Command{
				ParentCommand: cmdPlugin.ParentCommand,
				Name:          cmdPlugin.Name,
				Usage:         cmdPlugin.Usage,
				ShortDesc:     cmdPlugin.ShortDesc,
				LongDesc:      cmdPlugin.LongDesc,
				NumArgs:       cmdPlugin.NumArgs,
			},
		}

		file, err := json.MarshalIndent(cachedCommand, "", "")
		if err != nil {
			return err
		}

		err = os.WriteFile(targetFile, file, 0644)
		if err != nil {
			return err
		}
	}

	for _, hookPlugin := range m.hookPlugins {
		fileParts := strings.Split(hookPlugin.PluginDir, "/")
		fileName := strings.Trim(fileParts[len(fileParts)-1], "_hook")
		targetFile := path.Join(cacheHome, (fileName + ".hook.json"))

		_, err = os.Stat(targetFile)
		if !errors.Is(err, os.ErrNotExist) {
			continue
		} else if err == nil {
			continue
		}

		cachedCommand := CachedHook{
			ModuleName: hookPlugin.ModuleName,
			PluginDir:  hookPlugin.PluginDir,
			Hook: plugintypes.Hook{
				ParentCommand: hookPlugin.ParentCommand,
				Name:          hookPlugin.Name,
				HookType:      hookPlugin.HookType,
			},
		}

		file, err := json.MarshalIndent(cachedCommand, "", "")
		if err != nil {
			return err
		}

		err = os.WriteFile(targetFile, file, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

// Retrieve gets the cached plugins into memory
func (m *Manager) RetrieveCached(ctx context.Context) error {
	cacheHome, err := formatPluginHome(m.ChainId, "cached")
	if err != nil {
		return err
	}

	cachedCommandFiles, err := listFilesMatch(cacheHome, "*.cmd.json")
	if err != nil {
		return err
	}

	cachedHookFiles, err := listFilesMatch(cacheHome, "*.hook.json")
	if err != nil {
		return err
	}

	// Add checking techniques that clean up unused plugins from cache, etc.
	var cachedCommands []CachedCommand
	for _, cachedCommand := range cachedCommandFiles {
		cachedFileName := cachedCommand.Name()
		cachedFileNamePath := path.Join(cacheHome, cachedFileName)
		file, err := os.ReadFile(cachedFileNamePath)
		if err != nil {
			return err
		}

		cachedCommand := CachedCommand{}
		err = json.Unmarshal([]byte(file), &cachedCommand)
		if err != nil {
			return err
		}

		cachedCommands = append(cachedCommands, cachedCommand)
	}

	var cachedHooks []CachedHook
	for _, cachedHook := range cachedHookFiles {
		cachedFileName := cachedHook.Name()
		cachedFileNamePath := path.Join(cacheHome, cachedFileName)
		file, err := os.ReadFile(cachedFileNamePath)
		if err != nil {
			return err
		}

		cachedHook := CachedHook{}
		err = json.Unmarshal([]byte(file), &cachedHook)
		if err != nil {
			return err
		}

		cachedHooks = append(cachedHooks, cachedHook)
	}

	for _, cachedCommand := range cachedCommands {
		NewPluginMap := map[string]plugin.Plugin{
			cachedCommand.ModuleName: &plugintypes.CommandModulePlugin{},
		}

		m.cmdPlugins = append(m.cmdPlugins, ExtractedCommandModule{
			ModuleName:    cachedCommand.ModuleName,
			PluginDir:     cachedCommand.PluginDir,
			ParentCommand: cachedCommand.Command.ParentCommand,
			Name:          cachedCommand.Command.Name,
			Usage:         cachedCommand.Command.Usage,
			ShortDesc:     cachedCommand.Command.ShortDesc,
			LongDesc:      cachedCommand.Command.LongDesc,
			NumArgs:       cachedCommand.Command.NumArgs,
			Exec: func(cmd *cobra.Command, args []string) error {
				client := plugin.NewClient(&plugin.ClientConfig{
					HandshakeConfig: HandshakeConfig,
					Plugins:         NewPluginMap,
					Cmd:             exec.Command(cachedCommand.PluginDir),
					Logger:          pluginLogger,
				})

				rpcClient, err := client.Client()
				if err != nil {
					return err
				}

				raw, err := rpcClient.Dispense(cachedCommand.ModuleName)
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
	}

	for _, cachedHook := range cachedHooks {
		NewPluginMap := map[string]plugin.Plugin{
			cachedHook.ModuleName: &plugintypes.HookModulePlugin{},
		}

		m.hookPlugins = append(m.hookPlugins, ExtractedHookModule{
			ModuleName:    cachedHook.ModuleName,
			PluginDir:     cachedHook.PluginDir,
			ParentCommand: cachedHook.Hook.ParentCommand,
			Name:          cachedHook.Hook.Name,
			HookType:      cachedHook.Hook.HookType,
			PreRun: func(cmd *cobra.Command, args []string) error {
				client := plugin.NewClient(&plugin.ClientConfig{
					HandshakeConfig: HandshakeConfig,
					Plugins:         NewPluginMap,
					Cmd:             exec.Command(cachedHook.PluginDir),
					Logger:          pluginLogger,
				})

				rpcClient, err := client.Client()
				if err != nil {
					return err
				}

				raw, err := rpcClient.Dispense(cachedHook.ModuleName)
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
					Cmd:             exec.Command(cachedHook.PluginDir),
					Logger:          pluginLogger,
				})

				rpcClient, err := client.Client()
				if err != nil {
					return err
				}

				raw, err := rpcClient.Dispense(cachedHook.ModuleName)
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
	}

	return nil
}
