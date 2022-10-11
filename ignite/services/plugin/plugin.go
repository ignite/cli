// Package plugin implements ignite plugin management.
// A ignite plugin is a binary which communicates with the ignite binary
// via RPC thanks to the github.com/hashicorp/go-plugin library.
package plugin

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/hashicorp/go-hclog"
	hplugin "github.com/hashicorp/go-plugin"
	"github.com/pkg/errors"

	"github.com/ignite/cli/ignite/chainconfig"
	"github.com/ignite/cli/ignite/pkg/cliui/clispinner"
	"github.com/ignite/cli/ignite/pkg/gocmd"
	"github.com/ignite/cli/ignite/pkg/xfilepath"
	"github.com/ignite/cli/ignite/services/chain"
)

// pluginsPath holds the plugin cache directory.
var pluginsPath = xfilepath.Join(
	chainconfig.ConfigDirPath,
	xfilepath.Path("plugins"),
)

// Plugin represents a ignite plugin.
type Plugin struct {
	// Embed the plugin configuration
	chainconfig.Plugin
	// Interface allows to communicate with the plugin via net/rpc.
	Interface Interface
	// If any error occurred during the plugin load, it's stored here
	Error error

	repoPath   string
	cloneURL   string
	cloneDir   string
	reference  string
	srcPath    string
	binaryName string

	client *hplugin.Client
}

// Load loads the plugins found in the chain config.
//
// There's 2 kinds of plugins, local or remote.
// Local plugins have their path starting with a `/`, while remote plugins
// don't.
// Local plugins are useful for development purpose.
// Remote plugins require to be fetched first, in $HOME/.ignite/plugins
// folder, then they are loaded from there.
//
// If an error occurs during a plugin load, it's not returned but rather stored
// in the Plugin.Error field. This prevents the loading of other plugins to be
// interrupted.
func Load(ctx context.Context, c *chain.Chain) ([]*Plugin, error) {
	conf, err := c.Config()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	pluginsDir, err := pluginsPath()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var plugins []*Plugin
	for _, cp := range conf.Plugins {
		p := newPlugin(pluginsDir, cp)
		p.load(ctx)
		plugins = append(plugins, p)
	}
	return plugins, nil
}

// Update removes the cache directory of plugins and fetch them again.
func Update(plugins ...*Plugin) error {
	for _, p := range plugins {
		err := p.clean()
		if err != nil {
			return err
		}
		p.fetch()
	}
	return nil
}

// newPlugin creates a Plugin from configuration.
func newPlugin(pluginsDir string, cp chainconfig.Plugin) *Plugin {
	var (
		p          = &Plugin{Plugin: cp}
		pluginPath = cp.Path
	)
	if pluginPath == "" {
		p.Error = errors.Errorf(`missing plugin property "path"`)
		return p
	}
	if strings.HasPrefix(pluginPath, "/") {
		// This is a local plugin, check if the file exists
		st, err := os.Stat(pluginPath)
		if err != nil {
			p.Error = errors.Wrapf(err, "local plugin path %q not found", pluginPath)
			return p
		}
		if !st.IsDir() {
			p.Error = errors.Errorf("local plugin path %q is not a dir", pluginPath)
			return p
		}
		p.srcPath = pluginPath
		p.binaryName = path.Base(pluginPath)
		return p
	}
	// This is a remote plugin, parse the URL
	if i := strings.LastIndex(pluginPath, "@"); i != -1 {
		// path contains a reference
		p.reference = pluginPath[i+1:]
		pluginPath = pluginPath[:i]
	}
	parts := strings.Split(pluginPath, "/")
	if len(parts) < 3 {
		p.Error = errors.Errorf("plugin path %q is not a valid repository URL", pluginPath)
		return p
	}
	p.repoPath = path.Join(parts[:3]...)
	p.cloneURL = "https://" + p.repoPath
	if len(p.reference) > 0 {
		p.repoPath += "@" + p.reference
	}
	p.cloneDir = path.Join(pluginsDir, p.repoPath)
	p.srcPath = path.Join(pluginsDir, p.repoPath, path.Join(parts[3:]...))
	p.binaryName = path.Base(pluginPath)
	return p
}

func (p *Plugin) KillClient() {
	if p.client != nil {
		p.client.Kill()
	}
}

func (p *Plugin) isLocal() bool {
	return p.cloneURL == ""
}

func (p *Plugin) binaryPath() string {
	return path.Join(p.srcPath, p.binaryName)
}

// load tries to fill p.Interface, ensuring the plugin is usable.
func (p *Plugin) load(ctx context.Context) {
	if p.Error != nil {
		return
	}
	_, err := os.Stat(p.srcPath)
	if err != nil {
		// srcPath found, need to fetch the plugin
		p.fetch()
		if p.Error != nil {
			return
		}
	}
	if p.isLocal() {
		// trigger rebuild for local plugin if binary is outdated
		if p.outdatedBinary() {
			p.build(ctx)
		}
	} else {
		// Check if binary is already build
		_, err = os.Stat(p.binaryPath())
		if err != nil {
			// binary not found, need to build it
			p.build(ctx)
		}
	}
	if p.Error != nil {
		return
	}
	// pluginMap is the map of plugins we can dispense.
	pluginMap := map[string]hplugin.Plugin{
		p.binaryName: &InterfacePlugin{},
	}
	// Create an hclog.Logger
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   fmt.Sprintf("plugin %s", p.Path),
		Output: os.Stderr,
		Level:  hclog.Error,
	})
	// We're a host! Start by launching the plugin process.
	p.client = hplugin.NewClient(&hplugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Logger:          logger,
		Cmd:             exec.Command(p.binaryPath()),
		SyncStderr:      os.Stderr,
		SyncStdout:      os.Stdout,
	})

	// Connect via RPC
	rpcClient, err := p.client.Client()
	if err != nil {
		p.Error = errors.Wrapf(err, "connecting")
		return
	}

	// Request the plugin
	raw, err := rpcClient.Dispense(p.binaryName)
	if err != nil {
		p.Error = errors.Wrapf(err, "dispensing")
		return
	}

	// We should have an Interface now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	p.Interface = raw.(Interface)
}

// fetch clones the plugin repository at the expected reference.
func (p *Plugin) fetch() {
	if p.isLocal() {
		return
	}
	if p.Error != nil {
		return
	}
	defer clispinner.New().SetText(fmt.Sprintf("Fetching plugin %q...", p.cloneURL)).Stop()

	var err error
	if p.reference == "" {
		// No reference provided, just clone
		_, err = git.PlainClone(p.cloneDir, false, &git.CloneOptions{
			URL: p.cloneURL,
		})
	} else {
		// Reference provided, clone using tag or branch reference, one of the two
		// should work. SHA-1 aren't supported.
		for _, ref := range []plumbing.ReferenceName{
			plumbing.NewTagReferenceName(p.reference),
			plumbing.NewBranchReferenceName(p.reference),
		} {
			_, err = git.PlainClone(p.cloneDir, false, &git.CloneOptions{
				URL:           p.cloneURL,
				ReferenceName: ref,
				// Try to limit number of commits but this option doesn't seem to work well
				Depth: 1,
			})
			if err == nil {
				break
			}
		}
	}
	if err != nil {
		p.Error = errors.Wrapf(err, "cloning %q", p.cloneURL)
	}
}

// build compiles the plugin binary.
func (p *Plugin) build(ctx context.Context) {
	if p.Error != nil {
		return
	}
	defer clispinner.New().SetText(fmt.Sprintf("Building plugin %q...", p.Path)).Stop()

	if err := gocmd.ModTidy(ctx, p.srcPath); err != nil {
		p.Error = errors.Wrapf(err, "go mod tidy")
		return
	}
	if err := gocmd.BuildAll(ctx, p.binaryName, p.srcPath, nil); err != nil {
		p.Error = errors.Wrapf(err, "go build")
		return
	}
}

// clean removes the plugin cache (only for remote plugins).
func (p *Plugin) clean() error {
	if p.Error != nil {
		// Dont try to clean plugins with error
		return nil
	}
	if p.isLocal() {
		// Not a remote plugin, nothing to clean
		return nil
	}
	// Clean the cloneDir, next time the ignite command will be invoked, the
	// plugin will be fetched again.
	err := os.RemoveAll(p.cloneDir)
	return errors.WithStack(err)
}

// outdatedBinary returns true if the plugin binary is older than the other
// files in p.srcPath.
// Also returns true if the plugin binary is absent.
func (p *Plugin) outdatedBinary() bool {
	var (
		binaryTime time.Time
		mostRecent time.Time
	)
	err := filepath.Walk(p.srcPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if path == p.binaryPath() {
			binaryTime = info.ModTime()
			return nil
		}
		t := info.ModTime()
		if mostRecent.IsZero() || t.After(mostRecent) {
			mostRecent = t
		}
		return nil
	})
	if err != nil {
		fmt.Printf("error while walking plugin source path %q\n", p.srcPath)
		return false
	}
	return mostRecent.After(binaryTime)
}
