//go:build !relayer

package list_test

import (
	"context"
	"fmt"
	"testing"

	envtest "github.com/ignite/cli/v29/integration"
)

func TestGenerateAnAppWithListAndVerify(t *testing.T) {
	var (
		name      = "blog"
		namespace = "github.com/test/" + name

		env     = envtest.New(t)
		app     = env.Scaffold(namespace)
		servers = app.RandomizeServerPorts()
	)

	app.Scaffold(
		"create a module",
		false,
		"s", "module", "example", "--require-registration",
	)

	app.Scaffold(
		"create a list",
		false,
		"s", "list", "user", "email",
	)

	app.Scaffold(
		"create a list with custom path and module",
		false,
		"s",
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
		"s",
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
		"textCoinsAlias:coins",
		"--no-simulation",
	)

	app.Scaffold(
		"create a list with bool",
		false,
		"s",
		"list",
		"document",
		"signed:bool",
		"--module",
		"example",
	)

	app.Scaffold(
		"create a list with custom field type",
		false,
		"s",
		"list",
		"custom",
		"document:Document",
		"--module",
		"example",
	)

	app.Scaffold(
		"should prevent creating a list with duplicated fields",
		true,
		"s", "list", "company", "name", "name",
	)

	app.Scaffold(
		"should prevent creating a list with unrecognized field type",
		true,
		"s", "list", "employee", "level:itn",
	)

	app.Scaffold(
		"should prevent creating an existing list",
		true,
		"s", "list", "user", "email",
	)

	app.Scaffold(
		"should prevent creating a list whose name is a reserved word",
		true,
		"s", "list", "map", "size:int",
	)

	app.Scaffold(
		"should prevent creating a list containing a field with a reserved word",
		true,
		"s", "list", "document", "type:int",
	)

	app.Scaffold(
		"create a list with no interaction message",
		false,
		"s", "list", "nomessage", "email", "--no-message",
	)

	app.Scaffold(
		"should prevent creating a list in a non existent module",
		true,
		"s", "list", "user", "email", "--module", "idontexist",
	)

	app.EnsureSteady()

	ctx, cancel := context.WithCancel(env.Ctx())
	defer cancel()

	go func() {
		app.MustServe(ctx)
	}()

	app.WaitChainUp(ctx, servers.API)

	txReponse := app.CLITx(
		ctx,
		servers.RPC,
		"blog",
		"create-user",
		"test@user.com",
	)

	txReponse = app.CLIQueryTx(
		ctx,
		servers.RPC,
		txReponse.TxHash,
	)

	apiReponse := app.APIQuery(
		ctx,
		servers.API,
		namespace,
		name,
		"user",
	)
	fmt.Println(apiReponse)
}

func TestGen(t *testing.T) {
	var (
		name      = "blog"
		namespace = "github.com/test/" + name

		env     = envtest.New(t)
		app     = env.Scaffold(namespace)
		servers = app.RandomizeServerPorts()
	)

	app.Scaffold(
		"create a module",
		false,
		"s", "module", "example", "--require-registration",
	)

	app.Scaffold(
		"create a list",
		false,
		"s", "list", "user", "email",
	)

	ctx, cancel := context.WithCancel(env.Ctx())
	defer cancel()

	go func() {
		app.MustServe(ctx)
	}()

	app.WaitChainUp(ctx, servers.API)

	txResponse := app.CLITx(
		ctx,
		servers.RPC,
		name,
		"create-user",
		"test@user.com",
	)
	fmt.Println(txResponse)

	txResponse = app.CLIQueryTx(
		ctx,
		servers.RPC,
		txResponse.TxHash,
	)
	fmt.Println(txResponse)

	queryReponse := app.CLIQuery(
		ctx,
		servers.RPC,
		name,
		"list-user",
	)
	fmt.Println(queryReponse)

	queryReponse = app.CLIQuery(
		ctx,
		servers.RPC,
		name,
		"get-user",
		"0",
	)
	fmt.Println(queryReponse)

	apiReponse := app.APIQuery(
		ctx,
		servers.API,
		namespace,
		name,
		"user",
	)
	fmt.Println(apiReponse)

	apiReponse = app.APIQuery(
		ctx,
		servers.API,
		namespace,
		name,
		"user",
		"0",
	)
	fmt.Println(apiReponse)
}
