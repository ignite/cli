//go:build !relayer

package params_test

import (
	"context"
	"testing"

	envtest "github.com/ignite/cli/v29/integration"
)

func TestCreateModuleConfigs(t *testing.T) {
	var (
		name      = "mars"
		namespace = "github.com/test/" + name

		env     = envtest.New(t)
		app     = env.ScaffoldApp(namespace)
		servers = app.RandomizeServerPorts()
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

	ctx, cancel := context.WithCancel(env.Ctx())
	defer cancel()

	go func() {
		app.MustServe(ctx)
	}()

	app.WaitChainUp(ctx, servers.API)
}
