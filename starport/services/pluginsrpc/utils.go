package pluginsrpc

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
	chaincfg "github.com/tendermint/starport/starport/chainconfig"
)

func getConfigPlugins(configPath string) ([]chaincfg.Plugin, error) {
	cfg, err := chaincfg.ParseFile(configPath)
	if err != nil {
		return nil, err
	}

	return cfg.Plugins, nil
}

// Searches file, for now does not query the plugin repo
func pluginDownloaded(chainId, pluginId string) (bool, error) {
	dst, err := formatPluginHome(chainId, "")
	if err != nil {
		return false, err
	}

	pluginDirs, err := listDirs(dst)
	if err != nil {
		return false, err
	}

	for _, plugDir := range pluginDirs {
		if pluginId == plugDir.Name() {
			return true, nil
		}
	}

	return false, nil
}

// Check if plugin-specified configuration is different from downloaded plugins
// For now, ONLY CHECKS DIRECTORY NAMES
// This is not adequate, because one could delete files from directories
func PluginsChanged(cfg chaincfg.Config, chainId string) (bool, error) {
	var configPluginNames []string
	var fileConfigNames []string

	for _, plug := range cfg.Plugins {
		plugId := getPluginId(plug)
		configPluginNames = append(configPluginNames, plugId)
	}

	dst, err := formatPluginHome(chainId, "")
	if err != nil {
		return false, err
	}

	pluginDirs, err := listDirs(dst)
	if err != nil {
		if len(configPluginNames) > 0 && os.IsNotExist(err) {
			return true, nil
		}

		return false, err
	}

	for _, plugDir := range pluginDirs {
		fileConfigNames = append(fileConfigNames, plugDir.Name())
	}

	return !reflect.DeepEqual(configPluginNames, fileConfigNames), nil
}

func validateParentCommand(rootCommand *cobra.Command, subCommand []string) error {
	// this takes args, idk if that is the same as path
	innerCommand, _, err := rootCommand.Find(subCommand)
	if err != nil {
		return err
	}

	if innerCommand != nil {
		return nil
	}

	return ErrCommandNotFound
}

func getPluginId(plug chaincfg.Plugin) string {
	var plugId string
	if plug.Name != "" {
		plugId = plug.Name
	} else {
		repoSplit := strings.Split(plug.Repo, "/")
		plugId = repoSplit[len(repoSplit)-1]
	}

	return plugId
}

func formatPluginHome(chainId string, pluginId string) (string, error) {
	configDirPath, err := chaincfg.ConfigDirPath()
	if err != nil {
		return "", err
	}

	if pluginId != "" {
		return path.Join(configDirPath, "local-chains", chainId, PLUGINS_DIR, pluginId), nil
	}

	return path.Join(configDirPath, "local-chains", chainId, PLUGINS_DIR), nil
}

func listDirs(dir string) ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	filteredFiles := []os.FileInfo{}
	for _, file := range files {
		if file.IsDir() {
			filteredFiles = append(filteredFiles, file)
		}
	}

	return filteredFiles, nil
}

func listDirsMatch(dir, pattern string) ([]os.FileInfo, error) {
	var filteredFiles []os.FileInfo
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			// Test this, if it is proper for telling an exe
			if info.Mode()&0111 != 0 {
				filteredFiles = append(filteredFiles, info)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return filteredFiles, nil
}

func listFilesMatch(dir, pattern string) ([]os.FileInfo, error) {
	var filteredFiles []os.FileInfo
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			filteredFiles = append(filteredFiles, info)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return filteredFiles, nil
}
