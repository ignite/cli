// Package sta provides access to swagger-typescript-api CLI.
package sta

import (
	"context"
	"path/filepath"

	"github.com/ignite/cli/ignite/pkg/cmdrunner/exec"
	"github.com/ignite/cli/ignite/pkg/nodetime"
)

// Generate generates client code and TS types to outPath from an OpenAPI spec that resides at specPath.
func Generate(ctx context.Context, outPath, specPath, moduleNameIndex string) error {
	command, cleanup, err := nodetime.Command(nodetime.CommandSTA)
	if err != nil {
		return err
	}
	defer cleanup()

	dir := filepath.Dir(outPath)
	file := filepath.Base(outPath)

	// command constructs the sta command.
	command = append(command, []string{
		"--module-name-index",
		moduleNameIndex,
		"-p",
		specPath,
		"-o",
		dir,
		"-n",
		file,
	}...)

	// execute the command.
	return exec.Exec(ctx, command, exec.IncludeStdLogsToError())
}
