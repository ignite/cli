//go:build !relayer
// +build !relayer

package integration_test

import (
	"path/filepath"
	"testing"

	"github.com/tendermint/starport/integration"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
)

func TestGenerateAnAppWithStargateWithListAndVerify(t *testing.T) {
	var (
		env  = envtest.NewEnv(t)
		path = env.Scaffold("blog")
	)

	env.Must(env.Exec("create a list",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "list", "user", "email"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a list with custom path",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "list", "AppPath", "email", "--path", "blog"),
			step.Workdir(filepath.Dir(path)),
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

	env.Must(env.Exec("create a list with custom field type",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "list", "custom", "document:Document"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("should prevent creating a list with duplicated fields",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "list", "company", "name", "name"),
			step.Workdir(path),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating a list with unrecognized field type",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "list", "employee", "level:itn"),
			step.Workdir(path),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating an existing list",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "list", "user", "email"),
			step.Workdir(path),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating a list whose name is a reserved word",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "list", "map", "size:int"),
			step.Workdir(path),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating a list containing a field with a reserved word",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "list", "document", "type:int"),
			step.Workdir(path),
		)),
		envtest.ExecShouldError(),
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
		env  = envtest.NewEnv(t)
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
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating an existing list",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "list", "user", "email", "--module", "example"),
			step.Workdir(path),
		)),
		envtest.ExecShouldError(),
	))

	env.EnsureAppIsSteady(path)
}
