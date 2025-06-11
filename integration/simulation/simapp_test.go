//go:build !relayer

package simulation_test

import (
	"testing"

	envtest "github.com/ignite/cli/v29/integration"
)

func TestGenerateAnAppAndSimulate(t *testing.T) {
	var (
		env = envtest.New(t)
		app = env.ScaffoldApp("github.com/test/blog")
	)

	app.Scaffold(
		"create a list",
		false,
		"list", "foo", "foobar",
	)

	app.Scaffold(
		"create an singleton type",
		false,
		"single", "baz", "foobar",
	)

	app.Scaffold(
		"create an singleton type",
		false,
		"list", "noSimapp", "foobar", "--no-simulation",
	)

	app.Scaffold(
		"create a message",
		false,
		"message", "msgFoo", "foobar",
	)

	app.Scaffold(
		"scaffold a new module",
		false,
		"module", "new_module",
	)

	app.Scaffold(
		"create a map",
		false,
		"map",
		"bar",
		"foobar",
		"--module",
		"new_module",
	)

	app.Simulate(100, 50)
}
