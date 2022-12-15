package env

import (
	"fmt"
	"os"
	"path"

	"github.com/ignite/cli/ignite/pkg/xfilepath"
)

const (
	debug     = "IGNT_DEBUG"
	configDir = "IGNT_CONFIG_DIR"
)

func DebugEnabled() bool {
	return os.Getenv(debug) == "1"
}

func ConfigDir() xfilepath.PathRetriever {
	return func() (string, error) {
		if dir := os.Getenv(configDir); dir != "" {
			if !path.IsAbs(dir) {
				panic(fmt.Sprintf("%s must be an absolute path", configDir))
			}
			return dir, nil
		}
		return xfilepath.JoinFromHome(xfilepath.Path(".ignite"))()
	}
}

func SetConfigDir(dir string) {
	err := os.Setenv(configDir, dir)
	if err != nil {
		panic(fmt.Sprintf("set config dir env: %v", err))
	}
}
