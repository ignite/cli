package docker

import (
	"os"
)

// IsInDocker reports whether if running inside a Docker container or not.
func IsInDocker() bool {
	_, err := os.Stat(os.ExpandEnv("$HOME/.dockerenv"))
	return err == nil
}
