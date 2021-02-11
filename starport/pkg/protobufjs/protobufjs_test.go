package protobufjs

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/tendermint/starport/starport/pkg/gomodule"
	"github.com/tendermint/starport/starport/pkg/protopath"
)

func ExampleGenerate() {
	projectPath := "/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/test"

	modfile, err := gomodule.ParseAt(projectPath)
	if err != nil {
		panic(err)
	}

	resolved, err := protopath.ResolveDependencyPaths(modfile.Require,
		protopath.NewModule("github.com/cosmos/cosmos-sdk", "proto"),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(
		Generate(
			context.Background(),
			".",
			"types",
			filepath.Join(projectPath, "proto"),
			resolved,
		),
	)
}
