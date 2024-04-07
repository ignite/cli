package main

import (
	"fmt"
	"os"

	"github.com/ignite/cli/ignite/internal/tools/gen-mig-diffs/cmd"
)

func main() {
	if err := cmd.NewRootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
