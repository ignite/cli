//go:build !relayer

package other_components_test

import (
	"testing"

	envtest "github.com/ignite/cli/v29/integration"
)

func TestGenerateAnAppWithQuery(t *testing.T) {
	var (
		env = envtest.New(t)
		app = env.ScaffoldApp("github.com/test/blog")
	)

	app.Scaffold(
		"create a query",
		false,
		"query",
		"foo",
		"text",
		"vote:int",
		"like:bool",
		"-r",
		"foo,bar:int,foobar:bool",
	)

	app.Scaffold(
		"create a query with custom path",
		false,
		"query",
		"AppPath",
		"text",
		"vote:int",
		"like:bool",
		"-r",
		"foo,bar:int,foobar:bool",
		"--path",
		"./blog",
	)

	app.Scaffold(
		"create a paginated query",
		false,
		"query",
		"bar",
		"text",
		"vote:int",
		"like:bool",
		"-r",
		"foo,bar:int,foobar:bool",
		"--paginated",
	)

	app.Scaffold(
		"create a custom field type",
		false,
		"type",
		"custom-type",
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
	)

	app.Scaffold(
		"create a query with the custom field type as a response",
		false,
		"query", "foobaz", "-r", "bar:CustomType",
	)

	app.Scaffold(
		"create a query with lowercase custom field type as a response",
		false,
		"query", "foobaz-lower", "-r", "bar:customType",
	)

	app.Scaffold(
		"create a query with the custom array field type as a response",
		false,
		"query", "foobaz-array", "-r", "bars:array.CustomType",
	)

	app.Scaffold(
		"create a query with lowercase custom array field type as a response",
		false,
		"query", "foobaz-array-lower", "-r", "bars:array.customType",
	)

	app.Scaffold(
		"should prevent creating a query with invalid custom array field type as a response",
		true,
		"query", "foobaz-invalid-array", "-r", "bars:array.UnknownType",
	)

	app.Scaffold(
		"should prevent using custom type in request params",
		true,
		"query", "bur", "bar:CustomType",
	)

	app.Scaffold(
		"should prevent using custom array type in request params",
		true,
		"query", "bur-array", "bar:array.CustomType",
	)

	app.Scaffold(
		"create an empty query",
		false,
		"query", "foobar",
	)

	app.Scaffold(
		"should prevent creating an existing query",
		true,
		"query", "foo", "bar",
	)

	app.Scaffold(
		"create a module",
		false,
		"module", "foo", "--require-registration",
	)

	app.Scaffold(
		"create a query in a module",
		false,
		"query",
		"foo",
		"text",
		"--module",
		"foo",
		"--desc",
		"foo bar foobar",
		"--response",
		"foo,bar:int,foobar:bool",
	)

	app.EnsureSteady()
}
