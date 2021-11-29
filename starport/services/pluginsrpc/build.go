package pluginsrpc

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/tendermint/starport/starport/pkg/cmdrunner/exec"
	"github.com/tendermint/starport/starport/pkg/gocmd"
)

// Build will take plugins in the chain directory and build them into the output directory.
func (m *Manager) Build(ctx context.Context) error {
	fmt.Println("🛠️ Building plugins...")

	// Get plugin home
	pluginHome, err := formatPluginHome(m.ChainId, "")
	if err != nil {
		return err
	}

	outputDir := path.Join(pluginHome, "output")
	dir, err := ioutil.ReadDir(outputDir)
	for _, d := range dir {
		os.RemoveAll(path.Join(outputDir, d.Name()))
	}

	for _, cfgPlugin := range m.Config.Plugins {
		pluginId := getPluginId(cfgPlugin)
		pluginDir := path.Join(pluginHome, pluginId)

		// Enter plugins directory, go mod tidy
		// Somehow have to account for remote dependencies in individual plugins
		if err := gocmd.ModTidy(ctx, pluginDir); err != nil {
			return err
		}

		if err := gocmd.ModVerify(ctx, pluginDir); err != nil {
			return err
		}

		if len(cfgPlugin.Subdir) > 0 {
			for _, subdir := range cfgPlugin.Subdir {
				subdirPluginDir := path.Join(pluginDir, subdir)
				if err := traversePluginFiles(ctx, pluginId, subdirPluginDir, outputDir); err != nil {
					return err
				}
			}
		} else {
			pluginDir = path.Join(pluginHome, cfgPlugin.Name)
			if err := traversePluginFiles(ctx, pluginId, pluginDir, outputDir); err != nil {
				return err
			}
			log.Println("not exiting my guy")
		}
	}

	return nil
}

// traversePluginFiles will find all source files in the plugin directory, and build them.
func traversePluginFiles(ctx context.Context, pluginId string, pluginDir string, outputDir string) error {
	cmdPlugins, err := listFilesMatch(pluginDir, "*.cmd.go")
	if err != nil {
		return nil
	}

	if len(cmdPlugins) > 0 {
		for _, pluginFile := range cmdPlugins {
			fileName := pluginFile.Name()
			outputName := pluginId + "_" + strings.Trim(fileName, ".cmd.go") + "_cmd"
			inputFile := path.Join(pluginDir, fileName)
			outputFile := path.Join(outputDir, outputName)
			err := buildPlugin(ctx, outputFile, inputFile, pluginDir, []string{})
			if err != nil {
				return nil
			}
		}
	}

	hookPlugins, err := listFilesMatch(pluginDir, "*.hook.go")
	if err != nil {
		return nil
	}

	if len(hookPlugins) > 0 {
		for _, pluginFile := range hookPlugins {
			fileName := pluginFile.Name()
			outputName := pluginId + "_" + strings.Trim(fileName, ".hook.go") + "_hook"
			inputFile := path.Join(pluginDir, fileName)
			outputFile := path.Join(outputDir, outputName)
			err := buildPlugin(ctx, outputFile, inputFile, pluginDir, []string{})
			if err != nil {
				return nil
			}
		}
	}

	return nil
}

// buildPlugin will directly build the plugin.
func buildPlugin(ctx context.Context, output string, pluginFile string, pluginDir string, flags []string) error {
	// No -buildmode flag, because go-plugin does not take shared object files, instead uses RPC with a binary.
	command := []string{
		gocmd.Name(),
		gocmd.CommandBuild,
		gocmd.FlagOut,
		output,
		pluginFile,
	}

	// Go module not detected in workdirectory, so change into plugin directory.
	if err := os.Chdir(pluginDir); err != nil {
		return err
	}

	command = append(command, flags...)
	if err := exec.Exec(ctx, command); err != nil {
		return err
	}

	return nil
}
