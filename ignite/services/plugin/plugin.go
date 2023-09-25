// Package plugin implements ignite plugin management.
// An ignite plugin is a binary which communicates with the ignite binary
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

	"github.com/hashicorp/go-hclog"
	hplugin "github.com/hashicorp/go-plugin"
	"github.com/pkg/errors"

	"github.com/ignite/cli/ignite/config"
	pluginsconfig "github.com/ignite/cli/ignite/config/plugins"
	"github.com/ignite/cli/ignite/pkg/env"
	"github.com/ignite/cli/ignite/pkg/events"
	"github.com/ignite/cli/ignite/pkg/gocmd"
	"github.com/ignite/cli/ignite/pkg/xfilepath"
	"github.com/ignite/cli/ignite/pkg/xgit"
	"github.com/ignite/cli/ignite/pkg/xurl"
)

// PluginsPath holds the plugin cache directory.
var PluginsPath = xfilepath.Mkdir(xfilepath.Join(
	config.DirPath,
	xfilepath.Path("plugins"),
))

// Plugin represents a ignite plugin.
type Plugin struct {
	// Embed the plugin configuration
	pluginsconfig.Plugin
	// Interface allows to communicate with the plugin via net/rpc.
	Interface Interface
	// If any error occurred during the plugin load, it's stored here
	Error error

	name      string
	repoPath  string
	cloneURL  string
	cloneDir  string
	reference string
	srcPath   string

	client *hplugin.Client

	// holds a cache of the plugin manifest to prevent mant calls over the rpc boundary
	manifest Manifest
	// If a plugin's ShareHost flag is set to true, isHost is used to discern if a
	// plugin instance is controlling the rpc server.
	isHost bool

	ev events.Bus
}

// Option configures Plugin.
type Option func(*Plugin)

// CollectEvents collects events from the chain.
func CollectEvents(ev events.Bus) Option {
	return func(p *Plugin) {
		p.ev = ev
	}
}

// Load loads the plugins found in the chain config.
//
// There's 2 kinds of plugins, local or remote.
// Local plugins have their path starting with a `/`, while remote plugins don't.
// Local plugins are useful for development purpose.
// Remote plugins require to be fetched first, in $HOME/.ignite/plugins folder, then they are loaded
// from there.
//
// If an error occurs during a plugin load, it's not returned but rather stored in the `Plugin.Error`
// field. This prevents the loading of other plugins to be interrupted.
func Load(ctx context.Context, plugins []pluginsconfig.Plugin, options ...Option) ([]*Plugin, error) {
	pluginsDir, err := PluginsPath()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var loaded []*Plugin
	for _, cp := range plugins {
		p := newPlugin(pluginsDir, cp, options...)
		p.load(ctx)

		loaded = append(loaded, p)
	}
	return loaded, nil
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
func newPlugin(pluginsDir string, cp pluginsconfig.Plugin, options ...Option) *Plugin {
	var (
		p = &Plugin{
			Plugin: cp,
		}
		pluginPath = cp.Path
	)
	if pluginPath == "" {
		p.Error = errors.Errorf(`missing plugin property "path"`)
		return p
	}

	// Apply the options
	for _, apply := range options {
		apply(p)
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
		p.name = path.Base(pluginPath)
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
	p.cloneURL, _ = xurl.HTTPS(p.repoPath)

	if len(p.reference) > 0 {
		ref := strings.ReplaceAll(p.reference, "/", "-")
		p.cloneDir = path.Join(pluginsDir, fmt.Sprintf("%s-%s", p.repoPath, ref))
		p.repoPath += "@" + p.reference
	} else {
		p.cloneDir = path.Join(pluginsDir, p.repoPath)
	}

	// Plugin can have a subpath within its repository. For example,
	// "github.com/ignite/plugins/plugin1" where "plugin1" is the subpath.
	repoSubPath := path.Join(parts[3:]...)

	p.srcPath = path.Join(p.cloneDir, repoSubPath)
	p.name = path.Base(pluginPath)

	return p
}

// KillClient kills the running plugin client.
func (p *Plugin) KillClient() {
	if p.manifest.SharedHost && !p.isHost {
		// Don't send kill signal to a shared-host plugin when this process isn't
		// the one who initiated it.
		return
	}

	if p.client != nil {
		p.client.Kill()
	}

	if p.isHost {
		deleteConfCache(p.Path)
		p.isHost = false
	}
}

func (p Plugin) binaryName() string {
	return fmt.Sprintf("%s.app", p.name)
}

func (p Plugin) binaryPath() string {
	return path.Join(p.srcPath, p.binaryName())
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

	if p.IsLocalPath() {
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
		p.name: &InterfacePlugin{},
	}
	// Create an hclog.Logger
	logLevel := hclog.Error
	if env.DebugEnabled() {
		logLevel = hclog.Trace
	}
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   fmt.Sprintf("plugin %s", p.Path),
		Output: os.Stderr,
		Level:  logLevel,
	})

	if checkConfCache(p.Path) {
		rconf, err := readConfigCache(p.Path)
		if err != nil {
			p.Error = err
			return
		}

		// We're attaching to an existing server, supply attachment configuration
		p.client = hplugin.NewClient(&hplugin.ClientConfig{
			HandshakeConfig: handshakeConfig,
			Plugins:         pluginMap,
			Logger:          logger,
			Reattach:        &rconf,
			SyncStderr:      os.Stderr,
			SyncStdout:      os.Stdout,
		})

	} else {
		// We're a host! Start by launching the plugin process.
		p.client = hplugin.NewClient(&hplugin.ClientConfig{
			HandshakeConfig: handshakeConfig,
			Plugins:         pluginMap,
			Logger:          logger,
			Cmd:             exec.Command(p.binaryPath()),
			SyncStderr:      os.Stderr,
			SyncStdout:      os.Stdout,
		})
	}

	// :Connect via RPC
	rpcClient, err := p.client.Client()
	if err != nil {
		p.Error = errors.Wrapf(err, "connecting")
		return
	}

	// Request the plugin
	raw, err := rpcClient.Dispense(p.name)
	if err != nil {
		p.Error = errors.Wrapf(err, "dispensing")
		return
	}

	// We should have an Interface now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	p.Interface = raw.(Interface)

	m, err := p.Interface.Manifest()
	if err != nil {
		p.Error = errors.Wrapf(err, "manifest load")
	}

	p.manifest = m

	// write the rpc context to cache if the plugin is declared as host.
	// writing it to cache as lost operation within load to assure rpc client's reattach config
	// is hydrated.
	if m.SharedHost && !checkConfCache(p.Path) {
		err := writeConfigCache(p.Path, *p.client.ReattachConfig())
		if err != nil {
			p.Error = err
			return
		}

		// set the plugin's rpc server as host so other plugin clients may share
		p.isHost = true
	}
}

// fetch clones the plugin repository at the expected reference.
func (p *Plugin) fetch() {
	if p.IsLocalPath() {
		return
	}
	if p.Error != nil {
		return
	}
	p.ev.Send(fmt.Sprintf("Fetching plugin %q", p.cloneURL), events.ProgressStart())
	defer p.ev.Send(fmt.Sprintf("Plugin fetched %q", p.cloneURL), events.ProgressFinish())

	urlref := strings.Join([]string{p.cloneURL, p.reference}, "@")
	err := xgit.Clone(context.Background(), urlref, p.cloneDir)
	if err != nil {
		p.Error = errors.Wrapf(err, "cloning %q", p.repoPath)
	}
}

// build compiles the plugin binary.
func (p *Plugin) build(ctx context.Context) {
	if p.Error != nil {
		return
	}
	p.ev.Send(fmt.Sprintf("Building plugin %q", p.Path), events.ProgressStart())
	defer p.ev.Send(fmt.Sprintf("Plugin built %q", p.Path), events.ProgressFinish())

	if err := gocmd.ModTidy(ctx, p.srcPath); err != nil {
		p.Error = errors.Wrapf(err, "go mod tidy")
		return
	}
	if err := gocmd.Build(ctx, p.binaryName(), p.srcPath, nil); err != nil {
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
	if p.IsLocalPath() {
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
