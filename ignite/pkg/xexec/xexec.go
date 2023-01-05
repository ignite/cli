package xexec

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ignite/cli/ignite/pkg/goenv"
)

// IsExec checks if a file is executable by anyone.
func IsExec(binaryPath string) (bool, error) {
	info, err := os.Stat(binaryPath)
	if err != nil {
		return false, err
	}

	if m := info.Mode(); !m.IsDir() && m&0o111 != 0 {
		return true, nil
	}

	return false, nil
}

// ResolveAbsPath searches for an executable file in the current
// working directory, the directories defined by the PATH environment
// variable and in the Go binary path. Once found returns the absolute
// path to the file.
func ResolveAbsPath(filePath string) (path string, err error) {
	// Check if file exists and it's an executable file
	if path, err = filepath.Abs(filePath); err == nil {
		if ok, _ := IsExec(path); ok {
			return path, nil
		}
	}

	// Search file in the directories defined by the PATH env variable
	path, err = exec.LookPath(filePath)
	if err == nil {
		return path, nil
	}

	// When PATH search fails check if file is located in the Go binary path
	path = filepath.Join(goenv.Bin(), filePath)
	if ok, _ := IsExec(path); ok {
		return path, nil
	}

	return path, err
}

// TryResolveAbsPath searches for an executable file in the current
// working directory, the directories defined by the PATH environment
// variable and in the Go binary path. Once found returns the absolute
// path to the file, or otherwise it returns the file path unmodified.
func TryResolveAbsPath(filePath string) string {
	if path, err := ResolveAbsPath(filePath); err == nil {
		return path
	}

	return filePath
}

// IsCommandAvailable checks if command is available on user's path.
func IsCommandAvailable(name string) bool {
	_, err := ResolveAbsPath(name)
	return err == nil
}
