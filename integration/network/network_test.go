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

	"github.com/ignite/cli/ignite/chainconfig"
	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/ignite/pkg/gomodule"
	envtest "github.com/ignite/cli/integration"
)

const (
	spnModule     = "github.com/tendermint/spn"
	spnRepoURL    = "https://" + spnModule
	spnConfigFile = "config_2.yml"
)

func cloneSPN(t *testing.T) string {
	path := t.TempDir()
	repo, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:      spnRepoURL,
		Progress: os.Stdout,
	})
	require.NoError(t, err)

	// Read spn version from go.mod
	gm, err := gomodule.ParseAt("../..")
	require.NoError(t, err)
	var spnVersion string
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
	require.NoError(t, err)
	if n := len(v.Pre); n > 0 {
		// pseudo version, need to extract hash
		spnVersion = strings.Split(v.Pre[n-1].VersionStr, "-")[1]
	}
	rev, err := repo.ResolveRevision(plumbing.Revision(spnVersion))
	require.NoError(t, err)
	w, err := repo.Worktree()
	require.NoError(t, err)
	err = w.Checkout(&git.CheckoutOptions{
		Hash: *rev,
	})
	require.NoError(t, err)
	t.Logf("Checkout spn to ref %q", rev)

	return path
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
		spnPath = cloneSPN(t)
		env     = envtest.New(t)
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

func TestNetworkPublishConfigGenesis(t *testing.T) {
	var (
		spnPath = cloneSPN(t)
		env     = envtest.New(t)
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
