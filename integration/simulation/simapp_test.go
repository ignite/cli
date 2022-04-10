//go:build !relayer
// +build !relayer

package simulation_test

import (
	"testing"

	"github.com/ignite-hq/cli/ignite/pkg/cmdrunner/step"
	envtest "github.com/ignite-hq/cli/integration"
)

func TestGenerateAnAppAndSimulate(t *testing.T) {
	var (
		env  = envtest.New(t)
		path = env.Scaffold("blog")
	)

	env.Must(env.Exec("create a list",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "list", "foo", "foobar"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create an singleton type",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "single", "baz", "foobar"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create an singleton type",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "list", "noSimapp", "foobar", "--no-simulation"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a message",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "msgFoo", "foobar"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("scaffold a new module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "module", "new_module"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a map",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"s",
				"map",
				"bar",
				"foobar",
				"--module",
				"new_module",
			),
			step.Workdir(path),
		)),
	))

	env.Simulate(path, 100, 50)
}
