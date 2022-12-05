package config

import (
	"os"

	"github.com/ignite/cli/ignite/pkg/xfilepath"
)

// DirPath returns the path of configuration directory of Ignite.
var DirPath = xfilepath.JoinFromHome(xfilepath.Path(".ignite"))

// CreateConfigDir creates config directory if it is not created yet.
func CreateConfigDir() error {
	path, err := DirPath()
	if err != nil {
		return err
	}

	return os.MkdirAll(path, 0o755)
}
