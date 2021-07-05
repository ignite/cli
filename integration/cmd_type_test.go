// +build !relayer

package integration_test

import (
	"testing"

	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
)

func TestGenerateAnAppWithStargateWithListAndVerify(t *testing.T) {
	var (
		env  = newEnv(t)
		path = env.Scaffold("blog")
	)

	env.Must(env.Exec("create a list",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "list", "user", "email"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a list with int",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "list", "employee", "name:string", "level:int"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a list with bool",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "list", "document", "signed:bool"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("should prevent creating a list with duplicated fields",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "list", "company", "name", "name"),
			step.Workdir(path),
		)),
		ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating a list with unrecognized field type",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "list", "employee", "level:itn"),
			step.Workdir(path),
		)),
		ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating an existing list",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "list", "user", "email"),
			step.Workdir(path),
		)),
		ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating a list whose name is a reserved word",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "list", "map", "size:int"),
			step.Workdir(path),
		)),
		ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating a list containing a field with a reserved word",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "list", "document", "type:int"),
			step.Workdir(path),
		)),
		ExecShouldError(),
	))

	env.Must(env.Exec("create a list with no interaction message",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "list", "nomessage", "email", "--no-message"),
			step.Workdir(path),
		)),
	))

	env.EnsureAppIsSteady(path)
}

func TestCreateListInCustomModuleWithStargate(t *testing.T) {
	var (
		env  = newEnv(t)
		path = env.Scaffold("blog")
	)

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

	env.Must(env.Exec("create a list in the app's module",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "list", "user", "email"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("should prevent creating a list in a non existent module",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "list", "user", "email", "--module", "idontexist"),
			step.Workdir(path),
		)),
		ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating an existing list",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "list", "user", "email", "--module", "example"),
			step.Workdir(path),
		)),
		ExecShouldError(),
	))

	env.EnsureAppIsSteady(path)
}

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

	env.EnsureAppIsSteady(path)
}

func TestCreateSingletonTypeWithStargate(t *testing.T) {
	var (
		env  = newEnv(t)
		path = env.Scaffold("blog")
	)

	env.Must(env.Exec("create an singleton type",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "single", "user", "email"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create an singleton type with no message",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "single", "no-message", "email", "--no-message"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a module",
		step.NewSteps(step.New(
			step.Exec("starport", "module", "create", "example", "--require-registration"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create another type",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "list", "user", "email", "--module", "example"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("should prevent creating an singleton type with a typename that already exist",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "single", "user", "email", "--module", "example"),
			step.Workdir(path),
		)),
		ExecShouldError(),
	))

	env.Must(env.Exec("create an singleton type in a custom module",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "single", "singleuser", "email", "--module", "example"),
			step.Workdir(path),
		)),
	))

	env.EnsureAppIsSteady(path)
}
