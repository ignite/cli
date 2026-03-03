package env

import (
	"fmt"
	"os"
	"path"

	"github.com/ignite/cli/v29/ignite/pkg/xfilepath"
)

const (
	DebugEnvVar     = "IGNT_DEBUG"
	ConfigDirEnvVar = "IGNT_CONFIG_DIR"
)

// SetDebug sets the debug environment variable to "1".
// This is used to enable debug mode in the application.
func SetDebug() {
	_ = os.Setenv(DebugEnvVar, "1")
}

// IsDebug checks if the debug environment variable is set to "1".
// This is used to determine if the application is running in debug mode.
func IsDebug() bool {
	return os.Getenv(DebugEnvVar) == "1"
}

func ConfigDir() xfilepath.PathRetriever {
	return func() (string, error) {
		if dir := os.Getenv(ConfigDirEnvVar); dir != "" {
			if !path.IsAbs(dir) {
				panic(fmt.Sprintf("%s must be an absolute path", ConfigDirEnvVar))
			}
			return dir, nil
		}
		return xfilepath.JoinFromHome(xfilepath.Path(".ignite"))()
	}
}

func SetConfigDir(dir string) {
	err := os.Setenv(ConfigDirEnvVar, dir)
	if err != nil {
		panic(fmt.Sprintf("set config dir env: %v", err))
	}
}
