// +build !relayer

package integration_test

import (
	"path/filepath"
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
			step.Exec(
				"starport",
				"s",
				"message",
				"do-foo",
				"text",
				"vote:int",
				"like:bool",
				"-r",
				"foo,bar:int,foobar:bool",
			),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a message with custom path",
		step.NewSteps(step.New(
			step.Exec(
				"starport",
				"s",
				"message",
				"app-path",
				"text",
				"vote:int",
				"like:bool",
				"-r",
				"foo,bar:int,foobar:bool",
				"--path",
				"blog",
			),
			step.Workdir(filepath.Dir(path)),
		)),
	))

	env.Must(env.Exec("should prevent creating an existing message",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "message", "do-foo", "bar"),
			step.Workdir(path),
		)),
		ExecShouldError(),
	))

	env.Must(env.Exec("create a message with a custom signer name",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "message", "do-bar", "bar", "--signer", "bar-doer"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a custom field type",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "type", "custom-type", "customField:uint"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a message with the custom field type",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "message", "foo-baz", "customField:CustomType"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a module",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "module", "foo", "--require-registration"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a message in a module",
		step.NewSteps(step.New(
			step.Exec(
				"starport",
				"s",
				"message",
				"do-foo",
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
