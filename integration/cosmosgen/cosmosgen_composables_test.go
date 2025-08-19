package cosmosgen_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	envtest "github.com/ignite/cli/v29/integration"
)

func TestCosmosGenScaffoldComposables(t *testing.T) {
	if envtest.IsCI {
		t.Skip("Skipping TestCosmosGenScaffoldComposables test in CI environment")
	}

	var (
		env = envtest.New(t)
		app = env.ScaffoldApp("github.com/test/blog")
	)

	const (
		withMsgModuleName    = "withmsg"
		withoutMsgModuleName = "withoutmsg"
	)

	app.Scaffold("add custom module with message", false, "module", withMsgModuleName)

	app.Scaffold(
		"create a message",
		false,
		"message",
		"mymessage",
		"myfield1",
		"myfield2:bool",
		"--module",
		withMsgModuleName,
	)

	app.Scaffold(
		"add custom module without message",
		false,
		"module",
		withoutMsgModuleName,
	)

	app.Scaffold(
		"create a type",
		false,
		"type",
		"mytype",
		"mytypefield",
		"--module",
		withoutMsgModuleName,
	)

	app.Scaffold(
		"create a query",
		false,
		"query",
		"myQuery",
		"mytypefield",
		"--module",
		withoutMsgModuleName,
	)

	composablesDirGenerated := filepath.Join(app.SourcePath(), "vue/src/composables")
	require.NoError(t, os.RemoveAll(composablesDirGenerated))

	app.Scaffold(
		"scaffold vue",
		false,
		"vue",
	)

	app.Generate(
		"generate composables",
		false,
		"composables",
		"--clear-cache",
	)

	expectedQueryModules := []string{
		"useCosmosAuthV1Beta1",
		"useCosmosAuthzV1Beta1",
		"useCosmosBankV1Beta1",
		"useCosmosBaseTendermintV1Beta1",
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
		"useBlogBlogV1",
		"useBlogWithmsgV1",
		"useBlogWithoutmsgV1",
	}

	for _, mod := range expectedQueryModules {
		_, err := os.Stat(filepath.Join(composablesDirGenerated, mod))
		if assert.False(t, os.IsNotExist(err), "missing composable %q in %s", mod, composablesDirGenerated) {
			assert.NoError(t, err)
		}
	}

	if t.Failed() {
		// list composables files
		composablesFiles, err := os.ReadDir(composablesDirGenerated)
		require.NoError(t, err)
		t.Log("Composables files:", len(composablesFiles))
		for _, file := range composablesFiles {
			t.Logf(" - %s", file.Name())
		}
	}
}
