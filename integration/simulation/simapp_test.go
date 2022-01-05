//go:build !relayer
// +build !relayer

package simulation_test

import (
	"testing"

	"github.com/tendermint/starport/integration"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
)

func TestGenerateAnAppAndSimulate(t *testing.T) {
	var (
		env  = envtest.New(t)
		path = env.Scaffold("blog")
	)

	env.Must(env.Exec("create a list",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "list", "foo", "foobar"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create an singleton type",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "single", "baz", "foobar"),
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
			step.Exec("starport", "s", "msgFoo", "foobar"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("scaffold a new module",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "module", "new_module"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a map",
		step.NewSteps(step.New(
			step.Exec(
				"starport",
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
