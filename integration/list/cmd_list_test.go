//go:build !relayer

package list_test

import (
	"testing"

	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	envtest "github.com/ignite/cli/integration"
)

func TestGenerateAnAppWithListAndVerify(t *testing.T) {
	var (
		env = envtest.New(t)
		app = env.Scaffold("github.com/test/blog")
	)

	env.Must(env.Exec("create a module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "module", "--yes", "example", "--require-registration"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create a list",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "list", "--yes", "user", "email"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create a list with custom path and module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"s",
				"list",
				"--yes",
				"AppPath",
				"email",
				"--path",
				"blog",
				"--module",
				"example",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create a custom type fields",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"s",
				"list",
				"--yes",
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
				"--no-simulation",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create a list with bool",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"s",
				"list",
				"--yes",
				"document",
				"signed:bool",
				"--module",
				"example",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create a list with custom field type",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"s",
				"list",
				"--yes",
				"custom",
				"document:Document",
				"--module",
				"example",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("should prevent creating a list with duplicated fields",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "list", "--yes", "company", "name", "name"),
			step.Workdir(app.SourcePath()),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating a list with unrecognized field type",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "list", "--yes", "employee", "level:itn"),
			step.Workdir(app.SourcePath()),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating an existing list",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "list", "--yes", "user", "email"),
			step.Workdir(app.SourcePath()),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating a list whose name is a reserved word",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "list", "--yes", "map", "size:int"),
			step.Workdir(app.SourcePath()),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating a list containing a field with a reserved word",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "list", "--yes", "document", "type:int"),
			step.Workdir(app.SourcePath()),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("create a list with no interaction message",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "list", "--yes", "nomessage", "email", "--no-message"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("should prevent creating a list in a non existent module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "list", "--yes", "user", "email", "--module", "idontexist"),
			step.Workdir(app.SourcePath()),
		)),
		envtest.ExecShouldError(),
	))

	app.EnsureSteady()
}
