//go:build !relayer

package other_components_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner/step"
	envtest "github.com/ignite/cli/v29/integration"
)

func TestCreateChainRegistry(t *testing.T) {
	var (
		env = envtest.New(t)
		app = env.ScaffoldApp("github.com/test/mars")
	)

	env.Must(env.Exec("create chain-registry files",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"s",
				"chain-registry",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	_, statErr := os.Stat(filepath.Join(app.SourcePath(), "chain.json"))
	require.False(t, os.IsNotExist(statErr), "chain.json cannot be found")

	_, statErr = os.Stat(filepath.Join(app.SourcePath(), "assetlist.json"))
	require.False(t, os.IsNotExist(statErr), "assetlist.json cannot be found")

	app.EnsureSteady()
}
