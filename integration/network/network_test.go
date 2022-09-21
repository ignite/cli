package network_test

import (
	"bytes"
	"context"
	"os"
	"path"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	envtest "github.com/ignite/cli/integration"
)

const (
	spnURL    = "git@github.com:tendermint/spn"
	spnBranch = "develop"
	spnHash   = "5da0c7ae019d376f782fa3baeb2c0ac5654f2d1f"
)

func cloneSPN(t *testing.T) string {
	path, err := os.MkdirTemp("", "spn")
	require.NoError(t, err)
	repo, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:           spnURL,
		ReferenceName: plumbing.NewBranchReferenceName("develop"),
		Progress:      os.Stdout,
	})
	w, err := repo.Worktree()
	require.NoError(t, err)
	err = w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(spnHash),
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		os.RemoveAll(path)
	})
	return path
}

func TestNetworkPublish(t *testing.T) {
	var (
		spnPath = cloneSPN(t)
		env     = envtest.New(t)
		spn     = env.App(
			spnPath,
			envtest.AppHomePath("/tmp/spnhome"),
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
	}()

	env.Must(spn.Serve("serve spn chain", envtest.ExecCtx(ctx)))

	require.NoError(t, isBackendAliveErr, "spn cannot get online in time")
}
