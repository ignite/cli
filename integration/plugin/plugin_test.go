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

	env.Must(env.Exec("add network plugin locally",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "plugin", "add", "github.com/ignite/cli-plugin-network@feb7a963661612ae1a6c1f9ccdbe3ea9f7cb8dd3"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("remove network plugin locally",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "plugin", "remove", "github.com/ignite/cli-plugin-network"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("add network plugin globally",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "plugin", "add", "github.com/ignite/cli-plugin-network@feb7a963661612ae1a6c1f9ccdbe3ea9f7cb8dd3", "-g"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("remove network plugin globally",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "plugin", "remove", "github.com/ignite/cli-plugin-network", "-g"),
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
