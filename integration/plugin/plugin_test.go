package plugin_test

import (
	"testing"

	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	envtest "github.com/ignite/cli/integration"
)

func TestAddPlugin(t *testing.T) {
	var (
		env = envtest.New(t)
		app = env.Scaffold("github.com/test/blog")
	)

	// TODO use network plugin once finalized

	env.Must(env.Exec("add network plugin locally",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "plugin", "add", "github.com/aljo242/test-plugin"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("remove network plugin locally",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "plugin", "remove", "github.com/aljo242/test-plugin"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("add network plugin globally",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "plugin", "add", "github.com/aljo242/test-plugin", "-g"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("remove network plugin globally",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "plugin", "remove", "github.com/aljo242/test-plugin", "-g"),
			step.Workdir(app.SourcePath()),
		)),
	))

	app.EnsureSteady()
}

func TestPluginScaffold(t *testing.T) {
	var (
		env = envtest.New(t)
	)

	env.Must(env.Exec("add a plugin",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "plugin", "scaffold", "test"),
			step.Workdir(env.TmpDir()),
		)),
	))
}
