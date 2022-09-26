package network_test

import (
	"bytes"
	"context"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/blang/semver"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/ignite/pkg/gomodule"
	envtest "github.com/ignite/cli/integration"
)

const (
	spnModule  = "github.com/tendermint/spn"
	spnRepoURL = "https://" + spnModule
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

func TestNetworkPublish(t *testing.T) {
	var (
		spnPath = cloneSPN(t)
		env     = envtest.New(t)
		spn     = env.App(
			spnPath,
			envtest.AppHomePath(t.TempDir()),
			envtest.AppConfigPath(path.Join(spnPath, "config_2.yml")),
		)
		servers = spn.Config().Host
	)

	var (
		ctx, cancel       = context.WithTimeout(env.Ctx(), envtest.ServeTimeout)
		isBackendAliveErr error
	)

	go func() {
		defer cancel()

		if isBackendAliveErr = env.IsAppServed(ctx, servers); isBackendAliveErr != nil {
			return
		}
		var b bytes.Buffer
		env.Exec("publish planet chain to spn",
			step.NewSteps(step.New(
				step.Exec(
					envtest.IgniteApp,
					"network", "chain", "publish",
					"https://github.com/lubtd/planet",
					"--local",
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
