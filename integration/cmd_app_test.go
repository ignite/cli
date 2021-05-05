// +build !relayer

package integration_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/starport/starport/pkg/chaintest"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
)

func TestGenerateAnAppAndVerify(t *testing.T) {
	var (
		env  = chaintest.New(t)
		path = env.Scaffold("blog")
	)

	_, statErr := os.Stat(filepath.Join(path, "config.yml"))
	require.False(t, os.IsNotExist(statErr), "config.yml cannot be found")

	env.EnsureAppIsSteady(path)
}

func TestGenerateAnAppWithWasmAndVerify(t *testing.T) {
	var (
		env  = chaintest.New(t)
		path = env.Scaffold("blog")
	)

	env.Must(env.Exec("add Wasm module",
		step.NewSteps(step.New(
			step.Exec("starport", "module", "import", "wasm"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("should not add Wasm module second time",
		step.NewSteps(step.New(
			step.Exec("starport", "module", "import", "wasm"),
			step.Workdir(path),
		)),
		chaintest.ExecShouldError(),
	))

	env.EnsureAppIsSteady(path)
}

func TestGenerateAStargateAppWithEmptyModuleAndVerify(t *testing.T) {
	var (
		env  = chaintest.New(t)
		path = env.Scaffold("blog")
	)

	env.Must(env.Exec("create a module",
		step.NewSteps(step.New(
			step.Exec("starport", "module", "create", "example"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("should prevent creating an existing module",
		step.NewSteps(step.New(
			step.Exec("starport", "module", "create", "example"),
			step.Workdir(path),
		)),
		chaintest.ExecShouldError(),
	))

	env.EnsureAppIsSteady(path)
}
