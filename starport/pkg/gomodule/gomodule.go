package gomodule

import (
	"io/ioutil"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

// ParseAt finds and parses go.mod at app's path.
func ParseAt(path string) (*modfile.File, error) {
	gomod, err := ioutil.ReadFile(filepath.Join(path, "go.mod"))
	if err != nil {
		return nil, err
	}
	return modfile.Parse("", gomod, nil)
}
