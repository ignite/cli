//go:build !relayer
// +build !relayer

package list_test

import (
	"path/filepath"
	"testing"

	envtest "github.com/ignite-hq/cli/integration"
	"github.com/ignite-hq/cli/ignite/pkg/cmdrunner/step"
)

func TestGenerateAnAppWithStargateWithListAndVerify(t *testing.T) {
	var (
		env  = envtest.New(t)
		path = env.Scaffold("blog")
	)

	env.Must(env.Exec("create a list",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "list", "user", "email"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a list with custom path",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "list", "AppPath", "email", "--path", "blog"),
			step.Workdir(filepath.Dir(path)),
		)),
	))

	env.Must(env.Exec("create a custom type fields",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"s",
				"list",
				"employee",
				"numInt:int",
				"numsInt:array.int",
				"numsIntAlias:ints",
				"numUint:uint",
				"numsUint:array.uint",
				"numsUintAlias:uints",
				"textString:string",
				"textStrings:array.string",
				"textStringsAlias:strings",
				"textCoin:coin",
				"textCoins:array.coin",
				"textCoinsAlias:coins",
			),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a list with bool",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "list", "document", "signed:bool"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a list with custom field type",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "list", "custom", "document:Document"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("should prevent creating a list with duplicated fields",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "list", "company", "name", "name"),
			step.Workdir(path),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating a list with unrecognized field type",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "list", "employee", "level:itn"),
			step.Workdir(path),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating an existing list",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "list", "user", "email"),
			step.Workdir(path),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating a list whose name is a reserved word",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "list", "map", "size:int"),
			step.Workdir(path),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating a list containing a field with a reserved word",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "list", "document", "type:int"),
			step.Workdir(path),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("create a list with no interaction message",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "list", "nomessage", "email", "--no-message"),
			step.Workdir(path),
		)),
	))

	env.EnsureAppIsSteady(path)
}

func TestCreateListInCustomModuleWithStargate(t *testing.T) {
	var (
		env  = envtest.New(t)
		path = env.Scaffold("blog")
	)

	env.Must(env.Exec("create a module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "module", "example", "--require-registration"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a list",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "list", "user", "email", "--module", "example"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a list in the app's module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "list", "user", "email"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("should prevent creating a list in a non existent module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "list", "user", "email", "--module", "idontexist"),
			step.Workdir(path),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating an existing list",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "list", "user", "email", "--module", "example"),
			step.Workdir(path),
		)),
		envtest.ExecShouldError(),
	))

	env.EnsureAppIsSteady(path)
}
