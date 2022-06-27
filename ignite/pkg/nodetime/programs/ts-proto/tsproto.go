// Package tsproto provides access to protoc-gen-ts_proto protoc plugin.
package tsproto

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ignite/cli/ignite/pkg/nodetime"
)

const pluginName = "protoc-gen-ts_proto"

// BinaryPath returns the path to the binary of the ts-proto plugin so it can be passed to
// protoc via --plugin option.
//
// protoc is very picky about binary names of its plugins. for ts-proto, binary name
// will be protoc-gen-ts_proto.
// see why: https://github.com/stephenh/ts-proto/blob/7f76c05/README.markdown#quickstart.
func BinaryPath() (path string, cleanup func(), err error) {
	var command []string

	command, cleanup, err = nodetime.Command(nodetime.CommandTSProto)
	if err != nil {
		return
	}

	tmpdir := os.TempDir()
	path = filepath.Join(tmpdir, pluginName)

	// comforting protoc by giving protoc-gen-ts_proto name to the plugin's binary.
	script := fmt.Sprintf(`#!/bin/bash
%s "$@"
`, strings.Join(command, " "))

	err = os.WriteFile(path, []byte(script), 0755)

	return
}
