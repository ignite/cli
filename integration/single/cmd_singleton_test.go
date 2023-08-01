//go:build !relayer

package single_test

import (
	"testing"

	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	envtest "github.com/ignite/cli/integration"
)

func TestCreateSingleton(t *testing.T) {
	var (
		env = envtest.New(t)
		app = env.Scaffold("github.com/test/blog")
	)

	env.Must(env.Exec("create an singleton type",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "single", "--yes", "user", "email"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create an singleton type with custom path",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "single", "--yes", "appPath", "email", "--path", app.SourcePath()),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create an singleton type with no message",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "single", "--yes", "no-message", "email", "--no-message"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create a module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "module", "--yes", "example", "--require-registration"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create another type",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "list", "--yes", "user", "email", "--module", "example"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create another type with a custom field type",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "list", "--yes", "user-detail", "user:User", "--module", "example"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("should prevent creating an singleton type with a typename that already exist",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "single", "--yes", "user", "email", "--module", "example"),
			step.Workdir(app.SourcePath()),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("create an singleton type in a custom module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "single", "--yes", "singleuser", "email", "--module", "example"),
			step.Workdir(app.SourcePath()),
		)),
	))

	app.EnsureSteady()
}
