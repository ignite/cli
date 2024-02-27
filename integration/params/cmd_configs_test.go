//go:build !relayer

package params_test

import (
	"testing"

	"github.com/ignite/cli/v28/ignite/pkg/cmdrunner/step"
	envtest "github.com/ignite/cli/v28/integration"
)

func TestCreateModuleConfigs(t *testing.T) {
	var (
		env = envtest.New(t)
		app = env.Scaffold("github.com/test/mars")
	)

	env.Must(env.Exec("create a new module with configs",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"s",
				"module",
				"--yes",
				"foo",
				"--module-configs",
				"bla,baz:uint,bar:bool",
				"--require-registration"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("should prevent creating configs field that already exist",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"s",
				"configs",
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

	env.Must(env.Exec("create a new module configs in the foo module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"s",
				"configs",
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

	env.Must(env.Exec("create a new module configs in the mars module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"s",
				"configs",
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
