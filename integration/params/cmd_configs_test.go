//go:build !relayer

package params_test

import (
	"testing"

	envtest "github.com/ignite/cli/v29/integration"
)

func TestCreateModuleConfigs(t *testing.T) {
	var (
		env = envtest.New(t)
		app = env.Scaffold("github.com/test/mars")
	)

	app.Scaffold(
		"create a new module with configs",
		false,
		"module",
		"foo",
		"--module-configs",
		"bla,baz:uint,bar:bool",
		"--require-registration",
	)

	app.Scaffold(
		"should prevent creating configs field that already exist",
		true,
		"configs",
		"--yes",
		"bla",
		"buu:uint",
		"--module",
		"foo",
	)

	app.Scaffold(
		"create a new module configs in the foo module",
		false,
		"configs",
		"bol",
		"buu:uint",
		"plk:bool",
		"--module",
		"foo",
	)

	app.Scaffold(
		"create a new module configs in the mars module",
		false,
		"configs",
		"foo",
		"bar:uint",
		"baz:bool",
	)

	app.EnsureSteady()
}
