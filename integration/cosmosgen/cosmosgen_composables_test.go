package cosmosgen_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	envtest "github.com/ignite/cli/integration"
)

func TestCosmosGenScaffoldComposables(t *testing.T) {
	t.Skip()

	var (
		env = envtest.New(t)
		app = env.Scaffold("github.com/test/blog")
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

	composablesDirGenerated := filepath.Join(app.SourcePath(), "vue/src/composables")
	require.NoError(t, os.RemoveAll(composablesDirGenerated))

	env.Must(env.Exec("generate composables",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"g",
				"composables",
				"--yes",
				"--clear-cache",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	expectedQueryModules := []string{
		"useCosmosAuthV1Beta1",
		"useCosmosAuthzV1Beta1",
		"useCosmosBankV1Beta1",
		"useCosmosBaseTendermintV1Beta1",
		"useCosmosCrisisV1Beta1",
		"useCosmosDistributionV1Beta1",
		"useCosmosEvidenceV1Beta1",
		"useCosmosFeegrantV1Beta1",
		"useCosmosGovV1Beta1",
		"useCosmosGovV1",
		"useCosmosGroupV1",
		"useCosmosMintV1Beta1",
		"useCosmosNftV1Beta1",
		"useCosmosParamsV1Beta1",
		"useCosmosSlashingV1Beta1",
		"useCosmosStakingV1Beta1",
		"useCosmosTxV1Beta1",
		"useCosmosUpgradeV1Beta1",
		"useCosmosVestingV1Beta1",
		// custom modules
		"useTestBlogBlog",
		"useTestBlogWithmsg",
		"useTestBlogWithoutmsg",
	}

	for _, mod := range expectedQueryModules {

		_, err := os.Stat(filepath.Join(composablesDirGenerated, mod))
		if assert.False(t, os.IsNotExist(err), "missing composable %q in %s", mod, composablesDirGenerated) {
			assert.NoError(t, err)
		}
	}
}
