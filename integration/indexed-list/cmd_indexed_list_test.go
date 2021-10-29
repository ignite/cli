//go:build !relayer
// +build !relayer

package indexed_list_test

import (
	"testing"

	envtest "github.com/tendermint/starport/integration"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
)

func TestGenerateAnAppWithStargateWithListAndVerify(t *testing.T) {
	var (
		env  = envtest.New(t)
		path = env.Scaffold("blog")
	)

	env.Must(env.Exec("create an indexed list no index provided",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "indexed-list", "post", "content"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create an indexed list with a custom index",
		step.NewSteps(step.New(
			step.Exec("starport",
				"s",
				"indexed-list",
				"employee",
				"numInt:int",
				"numsInt:array.int",
				"numUint:uint",
				"numsUint:array.uint",
				"textString:string",
				"textStrings:array.string",
				"textCoin:coin",
				"textCoins:array.coin",
				"--index",
				"s,i:int,u:uint,b:bool",
			),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create an indexed list with custom field type",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "indexed-list", "custom", "employee:Employee"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("should prevent creating an indexed list with duplicated fields",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "indexed-list", "company", "name", "name"),
			step.Workdir(path),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating an indexed list with unrecognized field type",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "indexed-list", "company", "level:itn"),
			step.Workdir(path),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating an existing indexed list",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "indexed-list", "employee", "email"),
			step.Workdir(path),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating an indexed list whose name is a reserved word",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "indexed-list", "map", "size:int"),
			step.Workdir(path),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating an indexed list containing a field with a reserved word",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "indexed-list", "document", "type:int"),
			step.Workdir(path),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating an indexed list containing an invalid index",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "indexed-list", "document", "--index", "foo,foo"),
			step.Workdir(path),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("create an indexed list with no interaction message",
		step.NewSteps(step.New(
			step.Exec("starport", "s", "indexed-list", "no-message", "email", "--no-message"),
			step.Workdir(path),
		)),
	))

	env.EnsureAppIsSteady(path)
}

