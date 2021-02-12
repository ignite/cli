// Package protoc provides high level access to protoc command.
package protoc

import (
	"bytes"
	"context"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/protoanalysis"
)

// Generate generates code into outDir from protoPath and its includePaths by using plugins provided with protocOuts.
func Generate(ctx context.Context, outDir, protoPath string, includePaths, protocOuts []string) error {
	// start preparing the protoc command for execution.
	command := []string{
		"protoc",
	}

	// append third party proto locations to the command.
	for _, importPath := range append([]string{protoPath}, includePaths...) {
		// skip if a third party proto source actually doesn't exist on the filesystem.
		if _, err := os.Stat(importPath); os.IsNotExist(err) {
			continue
		}
		command = append(command, "-I", importPath)
	}

	// find out the list of proto files under the app and generate code for them.
	files, err := protoanalysis.SearchProto(protoPath)
	if err != nil {
		return err
	}

	includesThirdParty := func(filepath string) bool {
		for _, protoThirdPartyPath := range includePaths {
			if strings.HasPrefix(filepath, protoThirdPartyPath) {
				return true
			}
		}
		return false
	}

	for _, file := range files {
		// check if the file belongs to a third party proto. if so, skip it since it should
		// only be included via `-I`.
		if includesThirdParty(file) {
			continue
		}

		// run command for each protocOuts.
		for _, out := range protocOuts {
			command := append(command, out)
			command = append(command, file)

			errb := &bytes.Buffer{}

			err := cmdrunner.
				New(
					cmdrunner.DefaultStderr(errb),
					cmdrunner.DefaultWorkdir(outDir)).
				Run(ctx,
					step.New(step.Exec(command[0], command[1:]...)))

			if err != nil {
				return errors.Wrap(err, errb.String())
			}
		}
	}

	return nil
}
