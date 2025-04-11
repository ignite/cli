//go:build !relayer

package params_test

import (
	"context"
	"testing"

	envtest "github.com/ignite/cli/v29/integration"
)

func TestCreateModuleParameters(t *testing.T) {
	var (
		name      = "mars"
		namespace = "github.com/test/" + name

		env     = envtest.New(t)
		app     = env.ScaffoldApp(namespace)
		servers = app.RandomizeServerPorts()
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

	ctx, cancel := context.WithCancel(env.Ctx())
	defer cancel()

	go func() {
		app.MustServe(ctx)
	}()

	app.WaitChainUp(ctx, servers.API)
}
