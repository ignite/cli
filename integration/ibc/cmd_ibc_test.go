//go:build !relayer

package ibc_test

import (
	"testing"

	envtest "github.com/ignite/cli/v29/integration"
)

func TestCreateModuleWithIBC(t *testing.T) {
	var (
		env = envtest.New(t)
		app = env.ScaffoldApp("github.com/test/blogibc")
	)

	app.Scaffold(
		"create an IBC module",
		false,
		"module", "foo", "--ibc", "--require-registration",
	)

	app.Scaffold(
		"create an IBC module with custom path",
		false,
		"module",
		"appPath",
		"--ibc",
		"--require-registration",
		"--path",
		"./blogibc",
	)

	app.Scaffold(
		"create a type in an IBC module",
		false,
		"list", "user", "email", "--module", "foo",
	)

	app.Scaffold(
		"create an IBC module with an ordered channel",
		false,
		"module",
		"orderedfoo",
		"--ibc",
		"--ordering",
		"ordered",
		"--require-registration",
	)

	app.Scaffold(
		"create an IBC module with an unordered channel",
		false,
		"module",
		"unorderedfoo",
		"--ibc",
		"--ordering",
		"unordered",
		"--require-registration",
	)

	app.Scaffold(
		"create a non IBC module",
		false,
		"module", "non_ibc", "--require-registration",
	)

	app.Scaffold(
		"create an IBC module with dependencies",
		false,
		"module",
		"with_dep",
		"--ibc",
		"--dep",
		"auth,bank,staking,slashing",
		"--require-registration",
	)

	app.EnsureSteady()
}

func TestCreateIBCPacket(t *testing.T) {
	var (
		env = envtest.New(t)
		app = env.ScaffoldApp("github.com/test/blogibcb")
	)

	app.Scaffold(
		"create an IBC module",
		false,
		"module", "foo", "--ibc", "--require-registration",
	)

	app.Scaffold(
		"create a packet",
		false,
		"packet",
		"bar",
		"text",
		"texts:strings",
		"--module",
		"foo",
		"--ack",
		"foo:string,bar:int,baz:bool",
	)

	app.Scaffold(
		"should prevent creating a packet with no module specified",
		true,
		"packet", "bar", "text",
	)

	app.Scaffold(
		"should prevent creating a packet in a non existent module",
		true,
		"packet", "bar", "text", "--module", "nomodule",
	)

	app.Scaffold(
		"should prevent creating an existing packet",
		true,
		"packet", "bar", "post", "--module", "foo",
	)

	app.Scaffold(
		"create a packet with custom type fields",
		false,
		"packet",
		"ticket",
		"numInt:int",
		"numsInt:array.int",
		"numsIntAlias:ints",
		"numUint:uint",
		"numsUint:array.uint",
		"numsUintAlias:uints",
		"textString:string",
		"textStrings:array.string",
		"textStringsAlias:strings",
		"victory:bool",
		"textCoin:coin",
		"textCoins:array.coin",
		"--module",
		"foo",
	)

	app.Scaffold(
		"create a custom field type",
		false,
		"type", "custom-type", "customField:uint", "textCoinsAlias:coins", "--module", "foo",
	)

	app.Scaffold(
		"create a packet with a custom field type",
		false, "packet", "foo-baz", "customField:CustomType", "--module", "foo",
	)

	app.Scaffold(
		"create a packet with no send message",
		false, "packet", "nomessage", "foo", "--no-message", "--module", "foo",
	)

	app.Scaffold(
		"create a packet with no field",
		false, "packet", "empty", "--module", "foo",
	)

	app.Scaffold(
		"create a non-IBC module",
		false, "module", "bar", "--require-registration",
	)

	app.Scaffold(
		"should prevent creating a packet in a non IBC module",
		true, "packet", "foo", "text", "--module", "bar",
	)

	app.EnsureSteady()
}
