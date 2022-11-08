//go:build !relayer

package map_test

import (
	"path/filepath"
	"testing"

	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	envtest "github.com/ignite/cli/integration"
)

func TestCreateMap(t *testing.T) {
	var (
		env = envtest.New(t)
		app = env.Scaffold("github.com/test/blog")
	)

	env.Must(env.Exec("create a map",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "map", "--yes", "user", "user-id", "email"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create a map with custom path",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "map", "--yes", "appPath", "email", "--path", filepath.Join(app.SourcePath(), "app")),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create a map with no message",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "map", "--yes", "nomessage", "email", "--no-message"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create a module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "module", "--yes", "example", "--require-registration"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create a list",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"s",
				"list",
				"--yes",
				"user",
				"email",
				"--module",
				"example",
				"--no-simulation",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("should prevent creating a map with a typename that already exist",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "map", "--yes", "user", "email", "--module", "example"),
			step.Workdir(app.SourcePath()),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("create a map in a custom module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "map", "--yes", "mapUser", "email", "--module", "example"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create a map with a custom field type",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "map", "--yes", "mapDetail", "user:MapUser", "--module", "example"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create a map with Coin and []Coin",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"s",
				"map",
				"--yes",
				"salary",
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
				"--module",
				"example",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create a map with index",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"s",
				"map",
				"--yes",
				"map_with_index",
				"email",
				"emailIds:ints",
				"--index",
				"foo:string,bar:int,foobar:uint,barFoo:bool",
				"--module",
				"example",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create a map with invalid index",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"s",
				"map",
				"--yes",
				"map_with_invalid_index",
				"email",
				"--index",
				"foo:strings,bar:ints",
				"--module",
				"example",
			),
			step.Workdir(app.SourcePath()),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("create a message and a map with no-message flag to check conflicts",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "message", "--yes", "create-scavenge", "description"),
			step.Exec(envtest.IgniteApp, "s", "map", "--yes", "scavenge", "description", "--no-message"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("should prevent creating a map with duplicated indexes",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "map", "--yes", "map_with_duplicated_index", "email", "--index", "foo,foo"),
			step.Workdir(app.SourcePath()),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating a map with an index present in fields",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "map", "--yes", "map_with_invalid_index", "email", "--index", "email"),
			step.Workdir(app.SourcePath()),
		)),
		envtest.ExecShouldError(),
	))

	app.EnsureSteady()
}
