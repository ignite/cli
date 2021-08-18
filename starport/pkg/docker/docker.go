package docker

import (
	"os"
	"strings"
)

// IsInDocker reports whether if running inside a Docker container or not.
func IsInDocker() bool {
	content, err := os.ReadFile("/proc/1/cgroup")
	if err != nil {
		return false
	}

	return strings.Contains(string(content), "docker")
}
