package xos

import (
	"errors"
	"os"
	"strings"
)

// IsInPath checks if binpath is in system path.
func IsInPath(binpath string) error {
	paths := strings.Split(os.Getenv("PATH"), ":")
	for _, path := range paths {
		if path == binpath {
			return nil
		}
	}
	return errors.New("bin path is not in PATH")
}
