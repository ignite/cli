// +build !relayer

package integration_test

import (
	"testing"

	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
)

func TestGenerateAnAppWithLaunchpadWithTypeAndVerify(t *testing.T) {
	var (
		env  = newEnv(t)
		path = env.Scaffold("blog")
	)

	env.Must(env.Exec("create a type",
		step.NewSteps(step.New(
			step.Exec("starport", "type", "user", "email"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a type with int",
		step.NewSteps(step.New(
			step.Exec("starport", "type", "employee", "name:string", "level:int"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a type with bool",
		step.NewSteps(step.New(
			step.Exec("starport", "type", "document", "signed:bool"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("should prevent creating a type with duplicated fields",
		step.NewSteps(step.New(
			step.Exec("starport", "type", "company", "name", "name"),
			step.Workdir(path),
		)),
		ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating a type with unrecognized field type",
		step.NewSteps(step.New(
			step.Exec("starport", "type", "employee", "level:itn"),
			step.Workdir(path),
		)),
		ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating an existing type",
		step.NewSteps(step.New(
			step.Exec("starport", "type", "user", "email"),
			step.Workdir(path),
		)),
		ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating a type whose name is a reserved word",
		step.NewSteps(step.New(
			step.Exec("starport", "type", "map", "size:int"),
			step.Workdir(path),
		)),
		ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating a type containing a field with a reserved word",
		step.NewSteps(step.New(
			step.Exec("starport", "type", "document", "type:int"),
			step.Workdir(path),
		)),
		ExecShouldError(),
	))

	env.EnsureAppIsSteady(path)
}

func TestGenerateAnAppWithStargateWithTypeAndVerify(t *testing.T) {
	var (
		env  = newEnv(t)
		path = env.Scaffold("blog")
	)

	env.Must(env.Exec("create a type",
		step.NewSteps(step.New(
			step.Exec("starport", "type", "user", "email"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a type with int",
		step.NewSteps(step.New(
			step.Exec("starport", "type", "employee", "name:string", "level:int"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a type with bool",
		step.NewSteps(step.New(
			step.Exec("starport", "type", "document", "signed:bool"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("should prevent creating a type with duplicated fields",
		step.NewSteps(step.New(
			step.Exec("starport", "type", "company", "name", "name"),
			step.Workdir(path),
		)),
		ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating a type with unrecognized field type",
		step.NewSteps(step.New(
			step.Exec("starport", "type", "employee", "level:itn"),
			step.Workdir(path),
		)),
		ExecShouldError(),
	))

	env.Must(env.Exec("create a type with legacy scaffold",
		step.NewSteps(step.New(
			step.Exec("starport", "type", "participant", "email", "--legacy"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("should prevent creating an existing type",
		step.NewSteps(step.New(
			step.Exec("starport", "type", "user", "email"),
			step.Workdir(path),
		)),
		ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating a type whose name is a reserved word",
		step.NewSteps(step.New(
			step.Exec("starport", "type", "map", "size:int"),
			step.Workdir(path),
		)),
		ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating a type containing a field with a reserved word",
		step.NewSteps(step.New(
			step.Exec("starport", "type", "document", "type:int"),
			step.Workdir(path),
		)),
		ExecShouldError(),
	))

	env.EnsureAppIsSteady(path)
}

func TestCreateTypeInCustomModuleWithLaunchpad(t *testing.T) {
	var (
		env  = newEnv(t)
		path = env.Scaffold("blog")
	)

	env.Must(env.Exec("create a module",
		step.NewSteps(step.New(
			step.Exec("starport", "module", "create", "example"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a type",
		step.NewSteps(step.New(
			step.Exec("starport", "type", "user", "email", "--module", "example"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a type in the app's module",
		step.NewSteps(step.New(
			step.Exec("starport", "type", "user", "email"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("should prevent creating a type in a non existant module",
		step.NewSteps(step.New(
			step.Exec("starport", "type", "user", "email", "--module", "idontexist"),
			step.Workdir(path),
		)),
		ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating an existing type",
		step.NewSteps(step.New(
			step.Exec("starport", "type", "user", "email", "--module", "example"),
			step.Workdir(path),
		)),
		ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating an indexed type",
		step.NewSteps(step.New(
			step.Exec("starport", "type", "indexeduser", "email", "--indexed"),
			step.Workdir(path),
		)),
		ExecShouldError(),
	))

	env.EnsureAppIsSteady(path)
}

func TestCreateTypeInCustomModuleWithStargate(t *testing.T) {
	var (
		env  = newEnv(t)
		path = env.Scaffold("blog")
	)

	env.Must(env.Exec("create a module",
		step.NewSteps(step.New(
			step.Exec("starport", "module", "create", "example"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a type",
		step.NewSteps(step.New(
			step.Exec("starport", "type", "user", "email", "--module", "example"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a type with legacy scaffold",
		step.NewSteps(step.New(
			step.Exec("starport", "type", "participant", "email", "--legacy", "--module", "example"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a type in the app's module",
		step.NewSteps(step.New(
			step.Exec("starport", "type", "user", "email"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("should prevent creating a type in a non existent module",
		step.NewSteps(step.New(
			step.Exec("starport", "type", "user", "email", "--module", "idontexist"),
			step.Workdir(path),
		)),
		ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating an existing type",
		step.NewSteps(step.New(
			step.Exec("starport", "type", "user", "email", "--module", "example"),
			step.Workdir(path),
		)),
		ExecShouldError(),
	))

	env.EnsureAppIsSteady(path)
}

func TestCreateIndexTypeWithStargate(t *testing.T) {
	var (
		env  = newEnv(t)
		path = env.Scaffold("blog")
	)

	env.Must(env.Exec("create an indexed type",
		step.NewSteps(step.New(
			step.Exec("starport", "type", "user", "email", "--indexed"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a module",
		step.NewSteps(step.New(
			step.Exec("starport", "module", "create", "example"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a type",
		step.NewSteps(step.New(
			step.Exec("starport", "type", "user", "email", "--module", "example"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("should prevent creating an indexed type with a typename that already exist",
		step.NewSteps(step.New(
			step.Exec("starport", "type", "user", "email", "--indexed", "--module", "example"),
			step.Workdir(path),
		)),
		ExecShouldError(),
	))

	env.Must(env.Exec("create an indexed type in a custom module",
		step.NewSteps(step.New(
			step.Exec("starport", "type", "indexeduser", "email", "--indexed", "--module", "example"),
			step.Workdir(path),
		)),
	))

	env.EnsureAppIsSteady(path)
}
