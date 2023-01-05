package network_test

import (
	"bytes"
	"context"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/blang/semver/v4"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/require"

	ignitecmd "github.com/ignite/cli/ignite/cmd"
	chainconfig "github.com/ignite/cli/ignite/config/chain"
	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/ignite/pkg/gomodule"
	"github.com/ignite/cli/ignite/pkg/xgit"
	envtest "github.com/ignite/cli/integration"
)

const (
	spnModule            = "github.com/tendermint/spn"
	spnRepoURL           = "https://" + spnModule
	spnConfigFile        = "config_2.yml"
	pluginNetworkRepoURL = "https://" + ignitecmd.PluginNetworkPath
)

// setupSPN executes the following tasks:
// - clone cli-plugin-network to get the SPN version from go.mod
// - add the cloned cli-plugin-network as a global plugin
// - clone SPN to the expected version
// - returns the path of the cloned SPN.
func setupSPN(env envtest.Env) string {
	var (
		t          = env.T()
		require    = require.New(t)
		path       = t.TempDir()
		pluginPath = filepath.Join(path, "cli-plugin-network")
		spnPath    = filepath.Join(path, "spn")
		spnVersion string
	)
	// Clone the cli-plugin-network with the expected version
	err := xgit.Clone(context.Background(), pluginNetworkRepoURL, pluginPath)
	require.NoError(err)
	t.Logf("Checkout cli-plugin-revision to ref %q", ignitecmd.PluginNetworkPath)
	// Add plugin to config
	env.Must(env.Exec("add plugin network",
		step.NewSteps(step.New(
			// NOTE(tb): to test cli-plugin-network locally (can happen during dev)
			// comment the first line below and uncomment the second, with the
			// correct path to the plugin.
			step.Exec(envtest.IgniteApp, "plugin", "add", "-g", pluginPath),
			// step.Exec(envtest.IgniteApp, "plugin", "add", "-g", "/home/tom/src/ignite/cli-plugin-network"),
		)),
	))

	// Read spn version from plugin's go.mod
	gm, err := gomodule.ParseAt(pluginPath)
	require.NoError(err)
	for _, r := range gm.Require {
		if r.Mod.Path == spnModule {
			spnVersion = r.Mod.Version
			break
		}
	}
	if spnVersion == "" {
		t.Fatal("unable to read spn version from go.mod")
	}
	t.Logf("spn version found %q", spnVersion)

	// Check if spnVersion is a tag or a pseudo-version
	v, err := semver.ParseTolerant(spnVersion)
	require.NoError(err)
	if n := len(v.Pre); n > 0 {
		// pseudo version, need to extract hash
		spnVersion = strings.Split(v.Pre[n-1].VersionStr, "-")[1]
	}

	// Clone spn
	spnRepo, err := git.PlainClone(spnPath, false, &git.CloneOptions{
		URL: spnRepoURL,
	})
	require.NoError(err)
	// Checkout expected version
	rev, err := spnRepo.ResolveRevision(plumbing.Revision(spnVersion))
	require.NoError(err, spnVersion)
	w, err := spnRepo.Worktree()
	require.NoError(err)
	err = w.Checkout(&git.CheckoutOptions{
		Hash: *rev,
	})
	require.NoError(err, rev)
	t.Logf("Checkout spn to ref %q", rev)

	return spnPath
}

func migrateSPNConfig(t *testing.T, spnPath string) {
	configPath := filepath.Join(spnPath, spnConfigFile)
	rawCfg, err := os.ReadFile(configPath)
	require.NoError(t, err)

	version, err := chainconfig.ReadConfigVersion(bytes.NewReader(rawCfg))
	require.NoError(t, err)
	if version != chainconfig.LatestVersion {
		t.Logf("migrating spn config from v%d to v%d", version, chainconfig.LatestVersion)

		file, err := os.OpenFile(configPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o755)
		require.NoError(t, err)

		defer file.Close()

		err = chainconfig.MigrateLatest(bytes.NewReader(rawCfg), file)
		require.NoError(t, err)
	}
}

func TestNetworkPublish(t *testing.T) {
	var (
		env     = envtest.New(t)
		spnPath = setupSPN(env)
		spn     = env.App(
			spnPath,
			envtest.AppHomePath(t.TempDir()),
			envtest.AppConfigPath(path.Join(spnPath, spnConfigFile)),
		)
	)

	var (
		ctx, cancel       = context.WithTimeout(env.Ctx(), envtest.ServeTimeout)
		isBackendAliveErr error
	)

	// Make sure that the SPN config file is at the latest version
	migrateSPNConfig(t, spnPath)

	validator := spn.Config().Validators[0]
	servers, err := validator.GetServers()
	require.NoError(t, err)

	go func() {
		defer cancel()

		if isBackendAliveErr = env.IsAppServed(ctx, servers.API.Address); isBackendAliveErr != nil {
			return
		}
		var b bytes.Buffer
		env.Exec("publish planet chain to spn",
			step.NewSteps(step.New(
				step.Exec(
					envtest.IgniteApp,
					"network", "chain", "publish",
					"https://github.com/ignite/example",
					"--local",
					// The hash is used to be sure the test uses the right config
					// version. Hash value must be updated to the latest when the
					// config version in the repository is updated to a new version.
					"--hash", "b8b2cc2876c982dd4a049ed16b9a6099eca000aa",
				),
				step.Stdout(&b),
			)),
		)
		require.False(t, env.HasFailed(), b.String())
		t.Log(b.String())
	}()

	env.Must(spn.Serve("serve spn chain", envtest.ExecCtx(ctx)))

	require.NoError(t, isBackendAliveErr, "spn cannot get online in time")
}

func TestNetworkPublishGenesisConfig(t *testing.T) {
	var (
		env     = envtest.New(t)
		spnPath = setupSPN(env)
		spn     = env.App(
			spnPath,
			envtest.AppHomePath(t.TempDir()),
			envtest.AppConfigPath(path.Join(spnPath, spnConfigFile)),
		)
	)

	var (
		ctx, cancel       = context.WithTimeout(env.Ctx(), envtest.ServeTimeout)
		isBackendAliveErr error
	)

	// Make sure that the SPN config file is at the latest version
	migrateSPNConfig(t, spnPath)

	validator := spn.Config().Validators[0]
	servers, err := validator.GetServers()
	require.NoError(t, err)

	go func() {
		defer cancel()

		if isBackendAliveErr = env.IsAppServed(ctx, servers.API.Address); isBackendAliveErr != nil {
			return
		}
		var b bytes.Buffer
		env.Exec("publish test chain to spn",
			step.NewSteps(step.New(
				step.Exec(
					envtest.IgniteApp,
					"network", "chain", "publish",
					"https://github.com/aljo242/test",
					"--local",
					// The hash is used to be sure the test uses the right config
					// version. Hash value must be updated to the latest when the
					// config version in the repository is updated to a new version.
					"--hash", "23761c53ea364c7b046043f8c551eef992b99d22",
					// custom config genesis
					"--genesis-config", "gen.yml",
				),
				step.Stdout(&b),
			)),
		)
		require.False(t, env.HasFailed(), b.String())
		t.Log(b.String())
	}()

	env.Must(spn.Serve("serve spn chain", envtest.ExecCtx(ctx)))

	require.NoError(t, isBackendAliveErr, "spn cannot get online in time")
}
