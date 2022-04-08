package xexec

import "os/exec"

// IsCommandAvailable checks if command is available on user's path.
func IsCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
