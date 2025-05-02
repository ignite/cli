//go:build !relayer

package single_test

import (
	"testing"

	envtest "github.com/ignite/cli/v29/integration"
)

func TestCreateSingleton(t *testing.T) {
	var (
		env     = envtest.New(t)
		app     = env.ScaffoldApp("github.com/test/blog")
		servers = app.RandomizeServerPorts()
	)

	app.Scaffold(
		"create an singleton type",
		false,
		"single",
		"user", "email",
	)

	app.Scaffold(
		"create an singleton type with custom path",
		false,
		"single",
		"appPath", "email", "--path", app.SourcePath(),
	)

	app.Scaffold(
		"create an singleton type with no message",
		false,
		"single",
		"no-message", "email", "--no-message",
	)

	app.Scaffold(
		"create a module",
		false,
		"module",
		"example", "--require-registration",
	)

	app.Scaffold(
		"create another type",
		false,
		"list",
		"user", "email", "--module", "example",
	)

	app.Scaffold(
		"create another type with a custom field type",
		false,
		"list",
		"user-detail", "user:User", "--module", "example",
	)

	app.Scaffold(
		"should prevent creating an singleton type with a typename that already exist",
		true,
		"single",
		"user", "email", "--module", "example",
	)

	app.Scaffold(
		"create an singleton type in a custom module",
		false,
		"single",
		"singleuser", "email", "--module", "example",
	)

	app.EnsureSteady()

	app.RunChainAndSimulateTxs(servers)
}
