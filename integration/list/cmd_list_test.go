//go:build !relayer

package list_test

import (
	"testing"

	envtest "github.com/ignite/cli/v29/integration"
)

func TestGenerateAnAppWithListAndVerify(t *testing.T) {
	var (
		env     = envtest.New(t)
		app     = env.ScaffoldApp("github.com/test/blog")
		servers = app.RandomizeServerPorts()
	)

	app.Scaffold(
		"create a module",
		false,
		"module", "example", "--require-registration",
	)

	app.Scaffold(
		"create a list",
		false,
		"list", "user", "email",
	)

	app.Scaffold(
		"create a list with custom path and module",
		false,
		"list",
		"AppPath",
		"email",
		"--path",
		"blog",
		"--module",
		"example",
	)

	app.Scaffold(
		"create a custom type fields",
		false,
		"list",
		"employee",
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
		"--no-simulation",
	)

	app.Scaffold(
		"create a list with bool",
		false,
		"list",
		"document",
		"signed:bool",
		"textCoinsAlias:coins",
		"--module",
		"example",
	)

	app.Scaffold(
		"create a list with decimal coin",
		false,
		"list",
		"decimal",
		"deccointype:dec.coin",
		"deccoins:dec.coins",
		"--module",
		"example",
	)

	app.Scaffold(
		"create a list with custom field type",
		false,
		"list",
		"custom",
		"document:Document",
		"--module",
		"example",
	)

	app.Scaffold(
		"should prevent creating a list with duplicated fields",
		true,
		"list", "company", "name", "name",
	)

	app.Scaffold(
		"should prevent creating a list with unrecognized field type",
		true,
		"list", "invalidField", "level:itn",
	)

	app.Scaffold(
		"should prevent creating an existing list",
		true,
		"list", "user", "email",
	)

	app.Scaffold(
		"should prevent creating a list whose name is a reserved word",
		true,
		"list", "map", "size:int",
	)

	app.Scaffold(
		"should prevent creating a list containing a field with a reserved word",
		true,
		"list", "document", "type:int",
	)

	app.Scaffold(
		"create a list with no interaction message",
		false,
		"list", "nomessage", "email", "--no-message",
	)

	app.Scaffold(
		"should prevent creating a list in a non existent module",
		true,
		"list", "user", "email", "--module", "idontexist",
	)

	app.EnsureSteady()

	app.RunChainAndSimulateTxs(servers)
}
