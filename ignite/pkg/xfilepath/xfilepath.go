// Package xfilepath defines functions to define path retrievers that support error handling
package xfilepath

import (
	"os"
	"path/filepath"
)

// PathRetriever is a function that retrieve the contained path or an error
type PathRetriever func() (path string, err error)

// PathsRetriever is a function that retrieve the contained list of paths or an error
type PathsRetriever func() (path []string, err error)

// Path returns a path retriever from the provided path
func Path(path string) PathRetriever {
	return func() (string, error) { return path, nil }
}

// PathWithError returns a path retriever from the provided path and error
func PathWithError(path string, err error) PathRetriever {
	return func() (string, error) { return path, err }
}

// Join returns a path retriever from the join of the provided path retrievers
// The returned path retriever eventually returns the error from the first provided path retrievers that returns a non-nil error
func Join(paths ...PathRetriever) PathRetriever {
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

	return func() (string, error) {
		return path, err
	}
}

// JoinFromHome returns a path retriever from the join of the user home and the provided path retrievers
// The returned path retriever eventually returns the error from the first provided path retrievers that returns a non-nil error
func JoinFromHome(paths ...PathRetriever) PathRetriever {
	return Join(append([]PathRetriever{os.UserHomeDir}, paths...)...)
}

// List returns a paths retriever from a list of path retrievers
// The returned paths retriever eventually returns the error from the first provided path retrievers that returns a non-nil error
func List(paths ...PathRetriever) PathsRetriever {
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

	return func() ([]string, error) {
		return list, err
	}
}
