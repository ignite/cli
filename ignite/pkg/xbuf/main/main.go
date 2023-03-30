package main

import (
	"context"
	"github.com/ignite/cli/ignite/pkg/xbuf"
)

func main() {
	err := xbuf.Generate(
		context.Background(),
		"/Users/danilopantani/Desktop/go/src/github.com/ignite/earth/proto/buf.gen.gogo.yaml",
		"/Users/danilopantani/Desktop/go/src/github.com/ignite/earth/proto",
		".",
	)
	if err != nil {
		panic(err)
	}
}
