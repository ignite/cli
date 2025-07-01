//go:build !relayer

package other_components_test

import (
	"testing"

	envtest "github.com/ignite/cli/v29/integration"
)

func TestGenerateAnAppWithMessage(t *testing.T) {
	var (
		env = envtest.New(t)
		app = env.ScaffoldApp("github.com/test/blog")
	)

	app.Scaffold(
		"create a message",
		false,
		"message",
		"do-foo",
		"text",
		"vote:int",
		"like:bool",
		"from:address",
		"-r",
		"foo,bar:int,foobar:bool",
	)

	app.Scaffold(
		"create a message with custom path",
		false,
		"message",
		"app-path",
		"text",
		"vote:int",
		"like:bool",
		"-r",
		"foo,bar:int,foobar:bool",
		"--path",
		"blog",
		"--no-simulation",
	)

	app.Scaffold(
		"should prevent creating an existing message",
		true,
		"message", "do-foo", "bar",
	)

	app.Scaffold(
		"should prevent creating a message whose name only differs in capitalization",
		true,
		"message", "do-Foo", "bar",
	)

	app.Scaffold(
		"create a message with a custom signer name",
		false,
		"message", "--yes", "do-bar", "bar", "--signer", "bar-doer",
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
		"create a message with the custom field type",
		false,
		"message", "foo-baz", "customField:CustomType", "textCoinsAlias:coins",
	)

	app.Scaffold(
		"create a module",
		false,
		"module", "foo", "--require-registration",
	)

	app.Scaffold(
		"create a message in a module",
		false,
		"message",
		"do-foo",
		"text",
		"userIds:array.uint",
		"--module",
		"foo",
		"--desc",
		"foo bar foobar",
		"--response",
		"foo,bar:int,foobar:bool",
	)

	app.EnsureSteady()
}
