package xfilepath

import (
	"os"
	"path/filepath"
)

type PathRetriever func() (path string, err error)

type PathsRetriever func() (path []string, err error)

func Path(path string) PathRetriever {
	return func() (string, error) { return path, nil }
}

func PathWithError(path string, err error) PathRetriever {
	return func() (string, error) { return path, err }
}

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

func JoinFromHome(paths ...PathRetriever) PathRetriever {
	return Join(append([]PathRetriever{os.UserHomeDir}, paths...)...)
}

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
