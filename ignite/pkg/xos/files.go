package xos

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	JSONFile  = "json"
	ProtoFile = "proto"
	YAMLFile  = "yaml"
	YMLFile   = "yml"
)

type findFileOptions struct {
	extension []string
	prefix    string
}

type FindFileOptions func(o *findFileOptions)

// WithExtension adds a file extension to the search options.
// It can be called multiple times to add multiple extensions.
func WithExtension(extension string) FindFileOptions {
	return func(o *findFileOptions) {
		o.extension = append(o.extension, extension)
	}
}

func WithPrefix(prefix string) FindFileOptions {
	return func(o *findFileOptions) {
		o.prefix = prefix
	}
}

// FindFiles searches for files in the specified directory based on the given options.
// It supports filtering files by extension and prefix. Returns a list of matching files or an error.
func FindFiles(directory string, options ...FindFileOptions) ([]string, error) {
	opts := findFileOptions{}
	for _, apply := range options {
		apply(&opts)
	}

	files := make([]string, 0)
	return files, filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Filter by file extension if provided
		var matched bool
		for _, ext := range opts.extension {
			if filepath.Ext(path) == fmt.Sprintf(".%s", ext) {
				matched = true
				break
			}
		}

		if len(opts.extension) > 0 && !matched {
			return nil // Skip files that don't match the extension
		}

		// Filter by file prefix if provided
		if opts.prefix != "" && !strings.HasPrefix(filepath.Base(path), opts.prefix) {
			return nil // Skip files that don't match the prefix
		}

		// Add file to the result list if it is not a directory
		if !info.IsDir() {
			files = append(files, path)
		}

		return nil
	})
}

// FileExists check if a file from a given path exists.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
