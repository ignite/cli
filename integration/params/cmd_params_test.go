//go:build !relayer

package params_test

import (
	"testing"

	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	envtest "github.com/ignite/cli/integration"
)

func TestCreateModuleParameters(t *testing.T) {
	var (
		env = envtest.New(t)
		app = env.Scaffold("github.com/test/mars")
	)

	env.Must(env.Exec("create an module with parameter",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"s",
				"module",
				"--yes",
				"foo",
				"--params",
				"bla,baz:uint,bar:bool",
				"--require-registration"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("should prevent creating parameter field that already exist",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"s",
				"params",
				"--yes",
				"bla",
				"buu:uint",
				"--module",
				"foo",
			),
			step.Workdir(app.SourcePath()),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("create an new module parameters in the foo module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"s",
				"params",
				"--yes",
				"bol",
				"buu:uint",
				"plk:bool",
				"--module",
				"foo",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create an new module parameters in the mars module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"s",
				"params",
				"--yes",
				"foo",
				"bar:uint",
				"baz:bool",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	app.EnsureSteady()
}
