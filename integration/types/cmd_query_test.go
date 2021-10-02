//go:build !relayer
// +build !relayer

package integration_test

import (
	"path/filepath"
	"testing"

	"github.com/tendermint/starport/integration"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
)

func TestGenerateAnAppWithQuery(t *testing.T) {
	var (
		env  = envtest.NewEnv(t)
		path = env.Scaffold("blog")
	)

	env.Must(env.Exec("create a query",
		step.NewSteps(step.New(
			step.Exec(
				"starport",
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
				"starport",
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
				"starport",
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
			step.Exec("starport", "s", "type", "custom-type", "customField:uint"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a query with the custom field type",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "query", "foobaz", "bar:CustomType"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create an empty query",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "query", "foobar"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("should prevent creating an existing query",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "query", "foo", "bar"),
			step.Workdir(path),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("create a module",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "module", "foo", "--require-registration"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a query in a module",
		step.NewSteps(step.New(
			step.Exec(
				"starport",
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
