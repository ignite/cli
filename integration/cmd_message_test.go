// +build !relayer

package integration_test

import (
	"testing"

	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
)

func TestGenerateAnAppWithMessage(t *testing.T) {
	var (
		env  = newEnv(t)
		path = env.Scaffold("blog", Stargate)
	)

	env.Must(env.Exec("create a message",
		step.NewSteps(step.New(
			step.Exec("starport", "message", "foo", "foo bar foobar", "text", "vote:int", "like:bool", "--res", "foo,bar:int,foobar:bool"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("should prevent creating an existing message",
		step.NewSteps(step.New(
			step.Exec("starport", "message", "foo", "bar"),
			step.Workdir(path),
		)),
		ExecShouldError(),
	))

	env.Must(env.Exec("create a module",
		step.NewSteps(step.New(
			step.Exec("starport", "module", "create", "foo"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a message in a module",
		step.NewSteps(step.New(
			step.Exec("starport", "message", "foo", "foo bar foobar", "text", "--module", "foo"),
			step.Workdir(path),
		)),
	))

	env.EnsureAppIsSteady(path)
}
