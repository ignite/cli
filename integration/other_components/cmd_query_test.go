//go:build !relayer
// +build !relayer

package other_components_test

import (
	"path/filepath"
	"testing"

	"github.com/ignite-hq/cli/ignite/pkg/cmdrunner/step"
	envtest "github.com/ignite-hq/cli/integration"
)

func TestGenerateAnAppWithQuery(t *testing.T) {
	var (
		env  = envtest.New(t)
		path = env.Scaffold("blog")
	)

	env.Must(env.Exec("create a query",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"s",
				"query",
				"foo",
				"text",
				"vote:int",
				"like:bool",
				"-r",
				"foo,bar:int,foobar:bool",
			),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a query with custom path",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"s",
				"query",
				"AppPath",
				"text",
				"vote:int",
				"like:bool",
				"-r",
				"foo,bar:int,foobar:bool",
				"--path",
				"./blog",
			),
			step.Workdir(filepath.Dir(path)),
		)),
	))

	env.Must(env.Exec("create a paginated query",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"s",
				"query",
				"bar",
				"text",
				"vote:int",
				"like:bool",
				"-r",
				"foo,bar:int,foobar:bool",
				"--paginated",
			),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a custom field type",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"s",
				"type",
				"custom-type",
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

	env.Must(env.Exec("create a query with the custom field type as a response",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "query", "foobaz", "-r", "bar:CustomType"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("should prevent using custom type in request params",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "query", "bur", "bar:CustomType"),
			step.Workdir(path),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("create an empty query",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "query", "foobar"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("should prevent creating an existing query",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "query", "foo", "bar"),
			step.Workdir(path),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("create a module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "module", "foo", "--require-registration"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a query in a module",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"s",
				"query",
				"foo",
				"text",
				"--module",
				"foo",
				"--desc",
				"foo bar foobar",
				"--response",
				"foo,bar:int,foobar:bool",
			),
			step.Workdir(path),
		)),
	))

	env.EnsureAppIsSteady(path)
}
