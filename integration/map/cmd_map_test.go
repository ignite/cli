//go:build !relayer

package map_test

import (
	"path/filepath"
	"testing"

	envtest "github.com/ignite/cli/v29/integration"
)

func TestCreateMap(t *testing.T) {
	var (
		env     = envtest.New(t)
		app     = env.ScaffoldApp("github.com/test/blog")
		servers = app.RandomizeServerPorts()
	)

	app.Scaffold(
		"create a map",
		false,
		"map", "user", "user-id", "email",
	)

	app.Scaffold(
		"create a map with custom path",
		false,
		"map", "appPath", "email", "--path", filepath.Join(app.SourcePath(), "app"),
	)

	app.Scaffold(
		"create a map with no message",
		false,
		"map", "nomessage", "email", "--no-message",
	)

	app.Scaffold(
		"create a module",
		false,
		"module", "example", "--require-registration",
	)

	app.Scaffold(
		"create a list",
		false,
		"list",
		"user",
		"email",
		"--module",
		"example",
		"--no-simulation",
	)

	app.Scaffold(
		"create a map with decimal coin",
		false,
		"map",
		"decimal",
		"deccointype:dec.coin",
		"deccoins:dec.coins",
		"--module",
		"example",
	)

	app.Scaffold(
		"should prevent creating a map with a typename that already exist",
		true,
		"map", "user", "email", "--module", "example",
	)

	app.Scaffold(
		"create a map in a custom module",
		false,
		"map", "mapUser", "email", "--module", "example",
	)

	app.Scaffold(
		"create a map with a custom field type",
		false,
		"map", "mapDetail", "user:MapUser", "--module", "example",
	)

	app.Scaffold(
		"create a map with Coin and []Coin",
		false,
		"map",
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
		"textCoinsAlias:coins",
		"--module",
		"example",
	)

	app.Scaffold(
		"create a map with Coin and Coins",
		false,
		"map",
		"budget",
		"textCoin:coin",
		"textCoins:array.coin",
		"--module",
		"example",
	)

	app.Scaffold(
		"create a map with index",
		false,
		"map",
		"map_with_index",
		"email",
		"emailIds:ints",
		"--index",
		"bar:int",
		"--module",
		"example",
	)

	app.Scaffold(
		"create a map with invalid index (multi-index)",
		true,
		"map",
		"map_with_invalid_index",
		"email",
		"--index",
		"foo:strings,bar:int",
		"--module",
		"example",
	)

	app.Scaffold(
		"create a map with invalid index (invalid type)",
		true,
		"map",
		"map_with_invalid_index",
		"email",
		"--index",
		"foo:unknown",
		"--module",
		"example",
	)

	app.Scaffold(
		"create a message and a map with no-message flag to check conflicts",
		false,
		"message", "create-scavenge", "description",
	)

	app.Scaffold(
		"create a message and a map with no-message flag to check conflicts",
		false,
		"map", "scavenge", "description", "--no-message",
	)

	app.Scaffold(
		"should prevent creating a map with an index present in fields",
		true,
		"map", "map_with_invalid_index", "email", "--index", "email",
	)

	app.EnsureSteady()

	app.RunChainAndSimulateTxs(servers)
}
