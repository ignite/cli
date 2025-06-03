// Package xfilepath defines functions to define path retrievers that support error handling
package xfilepath

import (
	"os"
	"path/filepath"
)

// PathRetriever is a function that retrieves the contained path or an error.
type PathRetriever func() (path string, err error)

// PathsRetriever is a function that retrieves the contained list of paths or an error.
type PathsRetriever func() (path []string, err error)

// MustInvoke invokes the PathsRetriever func and panics if it returns an error.
func MustInvoke(p PathRetriever) string {
	path, err := p()
	if err != nil {
		panic(err)
	}
	return path
}

// Path returns a path retriever from the provided path.
func Path(path string) PathRetriever {
	return func() (string, error) { return path, nil }
}

// PathWithError returns a path retriever from the provided path and error.
func PathWithError(path string, err error) PathRetriever {
	return func() (string, error) { return path, err }
}

// Join returns a path retriever from the join of the provided path retrievers.
// The returned path retriever eventually returns the error from the first provided path retrievers
// that returns a non-nil error.
func Join(paths ...PathRetriever) PathRetriever {
	return func() (string, error) {
		var components []string
		var err error
		for _, path := range paths {
			var component string
			component, err = path()
			if err != nil {
				break
			}
			components = append(components, component)
		}
		path := filepath.Join(components...)
		return path, err
	}
}

// JoinFromHome returns a path retriever from the join of the user home and the provided path retrievers.
// The returned path retriever eventually returns the error from the first provided path retrievers that returns a non-nil error.
func JoinFromHome(paths ...PathRetriever) PathRetriever {
	return Join(append([]PathRetriever{os.UserHomeDir}, paths...)...)
}

// List returns a paths retriever from a list of path retrievers.
// The returned paths retriever eventually returns the error from the first provided path retrievers that returns a non-nil error.
func List(paths ...PathRetriever) PathsRetriever {
	return func() ([]string, error) {
		var list []string
		var err error
		for _, path := range paths {
			var resolved string
			resolved, err = path()
			if err != nil {
				break
			}
			list = append(list, resolved)
		}

		return list, err
	}
}

// Mkdir ensure path exists before returning it.
func Mkdir(path PathRetriever) PathRetriever {
	return func() (string, error) {
		p, err := path()
		if err != nil {
			return "", err
		}
		return p, os.MkdirAll(p, 0o755)
	}
}

// RelativePath return the relative app path from the current directory.
func RelativePath(appPath string) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	path, err := filepath.Rel(pwd, appPath)
	if err != nil {
		return "", err
	}
	return path, nil
}

// IsDir returns true if the path is a local directory.
func IsDir(path string) bool {
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		return true
	}
	return false
}

// MustAbs returns an absolute representation of path
// if the path is not absolute.
func MustAbs(path string) (string, error) {
	if !filepath.IsAbs(path) {
		return filepath.Abs(path)
	}
	return path, nil
}
