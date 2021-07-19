// +build !relayer

package integration_test

import (
	"testing"

	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
)

func TestCreateMapWithStargate(t *testing.T) {
	var (
		env  = newEnv(t)
		path = env.Scaffold("blog")
	)

	env.Must(env.Exec("create a map",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "map", "user", "email"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a map with no message",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "map", "nomessage", "email", "--no-message"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a module",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "module", "example", "--require-registration"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a list",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "list", "user", "email", "--module", "example"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("should prevent creating a map with a typename that already exist",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "map", "user", "email", "--module", "example"),
			step.Workdir(path),
		)),
		ExecShouldError(),
	))

	env.Must(env.Exec("create a map in a custom module",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "map", "mapuser", "email", "--module", "example"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a map with custom indexes",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "map", "map_with_index", "email", "--index", "foo:string,bar:int,foobar:uint,barFoo:bool"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("should prevent creating a map with duplicated indexes",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "map", "map_with_duplicated_index", "email", "--index", "foo,foo"),
			step.Workdir(path),
		)),
		ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating a map with an index present in fields",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "map", "map_with_invalid_index", "email", "--index", "email"),
			step.Workdir(path),
		)),
		ExecShouldError(),
	))

	env.EnsureAppIsSteady(path)
}
