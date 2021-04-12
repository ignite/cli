package services

import (
	"os"
	"path/filepath"
)

const (
	starportConfDir = ".starport"
)

// StarportConfPath returns the Starport Configuration directory
func StarportConfPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, starportConfDir), nil
}