package cosmosgen_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	envtest "github.com/ignite/cli/integration"
)

func TestCosmosGen(t *testing.T) {
	var (
		env          = envtest.New(t)
		path         = env.Scaffold("github.com/test/blog")
		dirGenerated = filepath.Join(path, "vue/src/store/generated")
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
			step.Workdir(path),
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
			step.Workdir(path),
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
			step.Workdir(path),
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
			step.Workdir(path),
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
			step.Workdir(path),
		)),
	))

	require.NoError(t, os.RemoveAll(dirGenerated))

	env.Must(env.Exec("generate vuex",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"g",
				"vuex",
				"--yes",
				"--proto-all-modules",
			),
			step.Workdir(path),
		)),
	))

	var expectedCosmosModules = []string{
		"cosmos.auth.v1beta1",
		"cosmos.authz.v1beta1",
		"cosmos.bank.v1beta1",
		"cosmos.base.tendermint.v1beta1",
		"cosmos.crisis.v1beta1",
		"cosmos.distribution.v1beta1",
		"cosmos.evidence.v1beta1",
		"cosmos.feegrant.v1beta1",
		"cosmos.gov.v1beta1",
		"cosmos.mint.v1beta1",
		"cosmos.params.v1beta1",
		"cosmos.slashing.v1beta1",
		"cosmos.staking.v1beta1",
		"cosmos.tx.v1beta1",
		"cosmos.upgrade.v1beta1",
		"cosmos.vesting.v1beta1",
	}

	var expectedCustomModules = []string{
		"test.blog.blog",
		"test.blog.withmsg",
		"test.blog.withoutmsg",
	}

	for _, chainModule := range expectedCustomModules {
		_, statErr := os.Stat(filepath.Join(dirGenerated, "test/blog", chainModule))
		require.False(t, os.IsNotExist(statErr), fmt.Sprintf("the %s vuex store should have be generated", chainModule))
		require.NoError(t, statErr)
	}

	chainDir, err := os.ReadDir(filepath.Join(dirGenerated, "test/blog"))
	require.Equal(t, len(expectedCustomModules), len(chainDir), "no extra modules should have been generated for test/blog")
	require.NoError(t, err)

	for _, cosmosModule := range expectedCosmosModules {
		_, statErr := os.Stat(filepath.Join(dirGenerated, "cosmos/cosmos-sdk", cosmosModule))
		require.False(t, os.IsNotExist(statErr), fmt.Sprintf("the %s code generation for module should have be made", cosmosModule))
		require.NoError(t, statErr)
	}

	cosmosDirs, err := os.ReadDir(filepath.Join(dirGenerated, "cosmos/cosmos-sdk"))
	require.Equal(t, len(expectedCosmosModules), len(cosmosDirs), "no extra modules should have been generated for cosmos/cosmos-sdk")
	require.NoError(t, err)
}
