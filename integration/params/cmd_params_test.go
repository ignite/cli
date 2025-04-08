//go:build !relayer

package params_test

import (
	"testing"

	envtest "github.com/ignite/cli/v29/integration"
)

func TestCreateModuleParameters(t *testing.T) {
	var (
		env = envtest.New(t)
		app = env.Scaffold("github.com/test/mars")
	)

	app.Scaffold(
		"create a new module with parameter",
		false,
		"module",
		"foo",
		"--params",
		"bla,baz:uint,bar:bool",
		"--require-registration",
	)

	app.Scaffold(
		"should prevent creating parameter field that already exist",
		true,
		"params",
		"bla",
		"buu:uint",
		"--module",
		"foo",
	)

	app.Scaffold(
		"create a new module parameters in the foo module",
		false,
		"params",
		"bol",
		"buu:uint",
		"plk:bool",
		"--module",
		"foo",
	)

	app.Scaffold(
		"create a new module parameters in the mars module",
		false,
		"params",
		"foo",
		"bar:uint",
		"baz:bool",
	)

	app.EnsureSteady()
}
