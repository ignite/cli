package main

import (
	"context"
	"github.com/ignite/cli/ignite/pkg/buf"
	"github.com/ignite/cli/ignite/pkg/cosmosgen"
	"path/filepath"
)

func main() {
	var (
		ctx       = context.Background()
		appPath   = "/Users/danilopantani/Desktop/go/src/github.com/ignite/mars"
		protoPath = filepath.Join(appPath, "proto")
	)
	if err := cosmosgen.InstallDepTools(ctx, appPath); err != nil {
		panic(err)
	}

	b, err := buf.New()
	if err != nil {
		panic(err)
	}

	err = b.Generate(
		ctx,
		protoPath,
		"/Users/danilopantani/Desktop",
		"buf.gen.gogo.yaml",
	)
	if err != nil {
		panic(err)
	}
}
