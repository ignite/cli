package cosmosgen_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner/step"
	envtest "github.com/ignite/cli/v29/integration"
)

func TestCosmosGenScaffold(t *testing.T) {
	if envtest.IsCI {
		t.Skip("Skipping CosmosGenScaffold test in CI environment")
	}

	var (
		env = envtest.New(t)
		app = env.ScaffoldApp("github.com/test/blog")
	)

	const (
		withMsgModuleName    = "withmsg"
		withoutMsgModuleName = "withoutmsg"
	)

	env.Must(env.Exec("add custom module with message",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"s",
				"module",
				"--yes",
				withMsgModuleName,
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create a message",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"s",
				"message",
				"--yes",
				"mymessage",
				"myfield1",
				"myfield2:bool",
				"--module",
				withMsgModuleName,
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("add custom module without message",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"s",
				"module",
				"--yes",
				withoutMsgModuleName,
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create a type",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"s",
				"type",
				"--yes",
				"mytype",
				"mytypefield",
				"--module",
				withoutMsgModuleName,
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create a query",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"s",
				"query",
				"--yes",
				"myQuery",
				"mytypefield",
				"--module",
				withoutMsgModuleName,
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	tsDirGenerated := filepath.Join(app.SourcePath(), "ts-client")
	require.NoError(t, os.RemoveAll(tsDirGenerated))

	env.Must(env.Exec("generate typescript",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"g",
				"ts-client",
				"--yes",
				"--clear-cache",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	expectedModules := []string{
		"cosmos.auth.v1beta1",
		"cosmos.authz.v1beta1",
		"cosmos.bank.v1beta1",
		"cosmos.base.tendermint.v1beta1",
		"cosmos.distribution.v1beta1",
		"cosmos.evidence.v1beta1",
		"cosmos.feegrant.v1beta1",
		"cosmos.gov.v1beta1",
		"cosmos.gov.v1",
		"cosmos.group.v1",
		"cosmos.mint.v1beta1",
		"cosmos.nft.v1beta1",
		"cosmos.params.v1beta1",
		"cosmos.slashing.v1beta1",
		"cosmos.staking.v1beta1",
		"cosmos.tx.v1beta1",
		"cosmos.upgrade.v1beta1",
		"cosmos.vesting.v1beta1",
		// custom modules
		"blog.blog.v1",
		"blog.withmsg.v1",
		"blog.withoutmsg.v1",
	}

	for _, mod := range expectedModules {
		_, err := os.Stat(filepath.Join(tsDirGenerated, mod))
		if assert.False(t, os.IsNotExist(err), "missing module %q in %s", mod, tsDirGenerated) {
			assert.NoError(t, err)
		}
	}

	if t.Failed() {
		// list ts-client files
		tsFiles, err := os.ReadDir(tsDirGenerated)
		require.NoError(t, err)
		t.Log("TS files:", len(tsFiles))
		for _, file := range tsFiles {
			t.Logf(" - %s", file.Name())
		}
	}
}
