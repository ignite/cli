// +build !relayer

package integration_test

import (
	"testing"

	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
)

func TestGenerateAnAppWithMessage(t *testing.T) {
	var (
		env  = newEnv(t)
		path = env.Scaffold("blog")
	)

	env.Must(env.Exec("create a message",
		step.NewSteps(step.New(
			step.Exec("starport", "message", "do-foo", "text", "vote:int", "like:bool", "-r", "foo,bar:int,foobar:bool"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("should prevent creating an existing message",
		step.NewSteps(step.New(
			step.Exec("starport", "message", "do-foo", "bar"),
			step.Workdir(path),
		)),
		ExecShouldError(),
	))

	env.Must(env.Exec("create a second message",
		step.NewSteps(step.New(
			step.Exec("starport", "message", "do-bar", "bar"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a module",
		step.NewSteps(step.New(
			step.Exec("starport", "module", "create", "foo", "--require-registration"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a message in a module",
		step.NewSteps(step.New(
			step.Exec("starport", "message", "do-foo", "text", "--module", "foo", "--desc", "foo bar foobar", "--response", "foo,bar:int,foobar:bool"),
			step.Workdir(path),
		)),
	))

	env.EnsureAppIsSteady(path)
}
