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

	env.Must(env.Exec("add plugin locally",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "plugin", "add", "github.com/ignite/example-plugin"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("remove plugin locally",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "plugin", "remove", "github.com/ignite/example-plugin"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("add plugin globally",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "plugin", "add", "github.com/ignite/example-plugin", "-g"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("remove plugin globally",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "plugin", "remove", "github.com/ignite/example-plugin", "-g"),
			step.Workdir(app.SourcePath()),
		)),
	))

	app.EnsureSteady()
}

// TODO install network plugin test

func TestPluginScaffold(t *testing.T) {
	env := envtest.New(t)

	env.Must(env.Exec("add a plugin",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "plugin", "scaffold", "test"),
			step.Workdir(env.TmpDir()),
		)),
	))
}
