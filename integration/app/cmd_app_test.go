//go:build !relayer

package app_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	envtest "github.com/ignite/cli/integration"
)

func TestGenerateAnApp(t *testing.T) {
	var (
		env = envtest.New(t)
		app = env.Scaffold("github.com/test/blog")
	)

	_, statErr := os.Stat(filepath.Join(app.SourcePath(), "x", "blog"))
	require.False(t, os.IsNotExist(statErr), "the default module should be scaffolded")

	app.EnsureSteady()
}

// TestGenerateAnAppWithName tests scaffolding a new chain using a local name instead of a GitHub URI.
func TestGenerateAnAppWithName(t *testing.T) {
	var (
		env = envtest.New(t)
		app = env.Scaffold("blog")
	)

	_, statErr := os.Stat(filepath.Join(app.SourcePath(), "x", "blog"))
	require.False(t, os.IsNotExist(statErr), "the default module should be scaffolded")

	app.EnsureSteady()
}

func TestGenerateAnAppWithNoDefaultModule(t *testing.T) {
	var (
		env = envtest.New(t)
		app = env.Scaffold("github.com/test/blog", "--no-module")
	)

	_, statErr := os.Stat(filepath.Join(app.SourcePath(), "x", "blog"))
	require.True(t, os.IsNotExist(statErr), "the default module should not be scaffolded")

	app.EnsureSteady()
}

func TestGenerateAnAppWithNoDefaultModuleAndCreateAModule(t *testing.T) {
	var (
		env = envtest.New(t)
		app = env.Scaffold("github.com/test/blog", "--no-module")
	)

	defer app.EnsureSteady()

	env.Must(env.Exec("should scaffold a new module into a chain that never had modules before",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "module", "--yes", "first_module"),
			step.Workdir(app.SourcePath()),
		)),
	))
}

func TestGenerateAnAppWithWasm(t *testing.T) {
	t.Skip()

	var (
		env = envtest.New(t)
		app = env.Scaffold("github.com/test/blog")
	)

	env.Must(env.Exec("add Wasm module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "wasm", "--yes"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("should not add Wasm module second time",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "wasm", "--yes"),
			step.Workdir(app.SourcePath()),
		)),
		envtest.ExecShouldError(),
	))

	app.EnsureSteady()
}

func TestGenerateAppWithEmptyModule(t *testing.T) {
	var (
		env = envtest.New(t)
		app = env.Scaffold("github.com/test/blog")
	)

	env.Must(env.Exec("create a module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "module", "--yes", "example", "--require-registration"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("should prevent creating an existing module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "module", "--yes", "example", "--require-registration"),
			step.Workdir(app.SourcePath()),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating a module with an invalid name",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "module", "--yes", "example1", "--require-registration"),
			step.Workdir(app.SourcePath()),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating a module with a reserved name",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "module", "--yes", "tx", "--require-registration"),
			step.Workdir(app.SourcePath()),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating a module with a forbidden prefix",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "module", "--yes", "ibcfoo", "--require-registration"),
			step.Workdir(app.SourcePath()),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating a module prefixed with an existing module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "module", "--yes", "examplefoo", "--require-registration"),
			step.Workdir(app.SourcePath()),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("create a module with dependencies",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"s",
				"module",
				"--yes",
				"with_dep",
				"--dep",
				"account,bank,staking,slashing,example",
				"--require-registration",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("should prevent creating a module with invalid dependencies",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"s",
				"module",
				"--yes",
				"with_wrong_dep",
				"--dep",
				"dup,dup",
				"--require-registration",
			),
			step.Workdir(app.SourcePath()),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating a module with a non registered dependency",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"s",
				"module",
				"--yes",
				"with_no_dep",
				"--dep",
				"inexistent",
				"--require-registration",
			),
			step.Workdir(app.SourcePath()),
		)),
		envtest.ExecShouldError(),
	))

	app.EnsureSteady()
}
