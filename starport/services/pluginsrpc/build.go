package pluginsrpc

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/tendermint/starport/starport/pkg/cmdrunner/exec"
	"github.com/tendermint/starport/starport/pkg/gocmd"
)

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
		pluginDir := path.Join(pluginHome, getPluginId(cfgPlugin))

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
				if err := traversePluginFiles(ctx, subdirPluginDir, outputDir); err != nil {
					return err
				}
			}
		} else {
			pluginDir = path.Join(pluginHome, cfgPlugin.Name)
			if err := traversePluginFiles(ctx, pluginDir, outputDir); err != nil {
				return err
			}
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
			outputName := strings.Trim(fileName, ".cmd.go") + "_cmd"
			inputFile := path.Join(pluginDir, fileName)
			outputFile := path.Join(outputDir, outputName)
			err := buildPlugin(ctx, outputFile, inputFile, pluginDir, []string{})
			if err != nil {
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
			outputName := strings.Trim(fileName, ".hook.go") + "_hook"
			inputFile := path.Join(pluginDir, fileName)
			outputFile := path.Join(outputDir, outputName)
			err := buildPlugin(ctx, outputFile, inputFile, pluginDir, []string{})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ERROR IN HERE
// maybe build is being run in some other directory without go.mod?
// its because the work directory is the underlying chain, and the chain does not have the dependencies.
// probably change work directory?
func buildPlugin(ctx context.Context, output string, pluginFile string, pluginDir string, flags []string) error {
	command := []string{
		gocmd.Name(),
		gocmd.CommandBuild,
		gocmd.FlagOut,
		output,
		pluginFile,
	}

	if err := os.Chdir(pluginDir); err != nil {
		return err
	}

	// Command is not using relative paths, main reason for error i think
	command = append(command, flags...)
	if err := exec.Exec(ctx, command); err != nil {
		return err
	}

	return nil
}
