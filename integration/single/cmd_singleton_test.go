//go:build !relayer
// +build !relayer

package single_test

import (
	"path/filepath"
	"testing"

	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	envtest "github.com/ignite/cli/integration"
)

func TestCreateSingletonWithStargate(t *testing.T) {
	var (
		env  = envtest.New(t)
		path = env.Scaffold("github.com/test/blog")
	)

	env.Must(env.Exec("create an singleton type",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "single", "--yes", "user", "email"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create an singleton type with custom path",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "single", "--yes", "appPath", "email", "--path", path),
			step.Workdir(filepath.Dir(path)),
		)),
	))

	env.Must(env.Exec("create an singleton type with no message",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "single", "--yes", "no-message", "email", "--no-message"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "module", "--yes", "example", "--require-registration"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create another type",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "list", "--yes", "user", "email", "--module", "example"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create another type with a custom field type",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "list", "--yes", "user-detail", "user:User", "--module", "example"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("should prevent creating an singleton type with a typename that already exist",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "single", "--yes", "user", "email", "--module", "example"),
			step.Workdir(path),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("create an singleton type in a custom module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "single", "--yes", "singleuser", "email", "--module", "example"),
			step.Workdir(path),
		)),
	))

	env.EnsureAppIsSteady(path)
}
