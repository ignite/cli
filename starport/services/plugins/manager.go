package plugins

import (
	"context"
	"io/ioutil"
	"os"
	"path"
	"plugin"
	"regexp"
	"strings"

	gogetter "github.com/hashicorp/go-getter"
	chaincfg "github.com/tendermint/starport/starport/chainconfig"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/exec"
	gocmd "github.com/tendermint/starport/starport/pkg/gocmd"
)

// NOTE: Go plugins were added in 1.8, so using them is questionable
// unless we can build the plugins with a greater go version and somehow
// grab the exposed functions with that build.

const (
	PLUGINS_DIR            = "plugins"
	CMD_SYMBOL_NAME        = "CmdPlugin"
	FLAG_BUILD_MODE_PLUGIN = "-buildmode=plugin"
)

// Plugin manager
type Manager struct {
	ChainId string `yaml:"chain_id"`

	cmdPlugins  []CmdPlugin
	hookPlugins []HookPlugin
}

func NewManager(chainId string) Manager {
	return Manager{
		ChainId: chainId,
	}
}

// Run will go through all of the steps:
// - Downloading
// - Building
// - Cache .so files
// - Execution and Injection
func (m *Manager) Run(ctx context.Context, cfg chaincfg.Config) error {
	// Download files, will overwrite (maybe check for remote changes?)
	if err := m.pull(ctx, cfg); err != nil {
		return err
	}

	// Build

	// Cache .so files

	// Extraction

	return nil
}

// MUST BE RAN BEFORE BUILD
func (m *Manager) pull(ctx context.Context, cfg chaincfg.Config) error {
	for _, plug := range cfg.Plugins {
		// Seperate individual plugins by ID
		var plugId string
		if plug.Name != "" {
			plugId = plug.Name
		} else {
			repoSplit := strings.Split(plug.Repo, "/")
			repoName := repoSplit[len(repoSplit)-1]
			if plug.Subdir != "" {
				plugId = repoName + "_" + plug.Subdir
			} else {
				plugId = repoName
			}
		}

		// Get plugin home
		dst, err := formatPluginHome(m.ChainId, plugId)
		if err != nil {
			return err
		}

		if err := download(plug.Repo, plug.Subdir, dst); err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) build(ctx context.Context, cfg chaincfg.Config) error {
	// Check for a change in file contents

	// Get plugin home
	dst, err := formatPluginHome(m.ChainId, "")
	if err != nil {
		return err
	}

	outputDir := path.Join(dst, "output")

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

		for _, pluginFile := range cmdPlugins {
			fileName := pluginFile.Name()
			inputFileDir := path.Join(pluginDir, fileName)
			outputFileDir := path.Join(outputDir, fileName)
			buildPlugin(ctx, outputFileDir, inputFileDir, []string{})
		}

		hookPlugins, err := listFiles(pluginDir, "*.hook.go")
		if err != nil {
			return err
		}

		for _, pluginFile := range hookPlugins {
			fileName := pluginFile.Name()
			inputFileDir := path.Join(pluginDir, fileName)
			outputFileDir := path.Join(outputDir, fileName)
			buildPlugin(ctx, outputFileDir, inputFileDir, []string{})
		}
	}

	return nil
}

// Caches .so files for rebuilding
func (m *Manager) cache(ctx context.Context, cfg chaincfg.Config) error {
	return nil
}

func (m *Manager) extractCommandPlugin(ctx context.Context, cfg chaincfg.Config) ([]CmdPlugin, error) {
	pluginsDir, err := formatPluginHome(m.ChainId, "")
	if err != nil {
		return nil, err
	}

	pluginFiles, err := listFiles(pluginsDir, `.*_command.so`)
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

		cmdPlugins = append(cmdPlugins, cmdPlugin)
	}

	m.cmdPlugins = cmdPlugins
	return cmdPlugins, nil
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

func download(repo string, subdir string, dst string) error {
	url := "https://" + repo + ".git"
	if subdir != "" {
		url += "//" + subdir
	}

	if err := gogetter.Get(dst, url); err != nil {
		return err
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

func listFiles(dir, pattern string) ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	filteredFiles := []os.FileInfo{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		matched, err := regexp.MatchString(pattern, file.Name())
		if err != nil {
			return nil, err
		}

		if matched {
			filteredFiles = append(filteredFiles, file)
		}
	}

	return filteredFiles, nil
}
