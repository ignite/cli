// Package tsproto provides access to protoc-gen-ts_proto protoc plugin.
package tsproto

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ignite/cli/ignite/pkg/nodetime"
)

const (
	pluginName     = "protoc-gen-ts_proto"
	scriptTemplate = "#!/bin/bash\n%s $@\n"
)

// BinaryPath returns the path to the binary of the ts-proto plugin, so it can be passed to
// protoc via --plugin option.
//
// protoc is very picky about binary names of its plugins. for ts-proto, binary name
// will be protoc-gen-ts_proto.
// see why: https://github.com/stephenh/ts-proto/blob/7f76c05/README.markdown#quickstart.
func BinaryPath() (path string, cleanup func(), err error) {
	// Create binary for the TypeScript protobuf generator
	command, cleanupBin, err := nodetime.Command(nodetime.CommandTSProto)
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			cleanupBin()
		}
	}()

	// Create a random directory for the script that runs the TypeScript protobuf generator.
	// This is required to avoid potential flaky integration tests caused by one concurrent
	// test overwriting the generator script while it is being run in a separate test process.
	tmpDir, err := os.MkdirTemp("", "ts_proto_plugin")
	if err != nil {
		return
	}

	cleanupScriptDir := func() { os.RemoveAll(tmpDir) }

	defer func() {
		if err != nil {
			cleanupScriptDir()
		}
	}()

	cleanup = func() {
		cleanupBin()
		cleanupScriptDir()
	}

	// Wrap the TypeScript protobuf generator in a script with a fixed name
	// located in a random temporary directory.
	script := fmt.Sprintf(scriptTemplate, strings.Join(command, " "))
	path = filepath.Join(tmpDir, pluginName)
	err = os.WriteFile(path, []byte(script), 0o755)

	return path, cleanup, err
}
