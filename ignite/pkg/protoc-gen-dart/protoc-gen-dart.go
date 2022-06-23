package protocgendart

import (
	"fmt"

	"github.com/ignite/cli/ignite/pkg/localfs"
	"github.com/ignite/cli/ignite/pkg/protoc-gen-dart/data"
)

// Name of the plugin.
const Name = "protoc-gen-dart"

// BinaryPath returns the binary path for the plugin.
func BinaryPath() (path string, cleanup func(), err error) {
	return localfs.SaveBytesTemp(data.Binary(), Name, 0755)
}

// Flag returns the binary name-binary path format to pass to protoc --plugin.
func Flag() (flag string, cleanup func(), err error) {
	path, cleanup, err := BinaryPath()
	flag = fmt.Sprintf("%s=%s", Name, path)
	return
}
