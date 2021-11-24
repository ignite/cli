package plugins

import (
	"context"
	"fmt"
	"path"
	"strings"

	chaincfg "github.com/tendermint/starport/starport/chainconfig"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/exec"
	"github.com/tendermint/starport/starport/pkg/gocmd"
)

func (m *Manager) build(ctx context.Context, cfg chaincfg.Config) error {
	// Get plugin home
	dst, err := formatPluginHome(m.ChainId, "")
	if err != nil {
		return err
	}

	outputDir := path.Join(dst, "output")

	for _, cfgPlugin := range cfg.Plugins {
		pluginDir := path.Join(dst, getPluginId(cfgPlugin))
		// Enter plugins directory, go get .
		// Somehow have to account for remote dependencies in individual plugins
		if err := gocmd.GetAll(ctx, pluginDir, nil); err != nil {
			return err
		}

		if cfgPlugin.Subdir != "" {
			pluginDir = path.Join(pluginDir, cfgPlugin.Subdir)
		} else {
			pluginDir = path.Join(dst, cfgPlugin.Name)
		}

		if err := traversePluginFiles(ctx, pluginDir, outputDir); err != nil {
			return err
		}
	}

	return nil
}

// Context?
func traversePluginFiles(ctx context.Context, pluginDir string, outputDir string) error {
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
			err := buildPlugin(ctx, outputFileDir, inputFileDir, []string{})
			if err != nil {
				fmt.Println(err.Error())
				return err
			}
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
			err := buildPlugin(ctx, outputFileDir, inputFileDir, []string{})
			if err != nil {
				fmt.Println(err.Error())
				return err
			}
		}
	}

	return nil
}

// ERROR IN HERE
func buildPlugin(ctx context.Context, output string, path string, flags []string) error {
	command := []string{
		"go",
		gocmd.CommandBuild,
		FLAG_BUILD_MODE_PLUGIN,
		gocmd.FlagOut,
		output,
		path,
	}

	// Command is not using relative paths, main reason for error i think
	fmt.Println(command)
	command = append(command, flags...)
	return exec.Exec(ctx, command)
}
