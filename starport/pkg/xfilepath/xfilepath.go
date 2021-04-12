package xfilepath

import (
	"os"
	"path/filepath"
)

type PathRetriever func() (path string, err error)

func Path(path string) PathRetriever {
	return func() (string, error) { return path, nil }
}

func PathWithError(path string, err error) PathRetriever {
	return func() (string, error) { return path, err }
}

func Join(paths ...PathRetriever) PathRetriever {
	return func() (string, error) {
		var components []string
		for _, path := range paths {
			component, err := path()
			if err != nil {
				return "", err
			}
			components = append(components, component)
		}
		return filepath.Join(components...), nil
	}
}

func JoinFromHome(paths ...PathRetriever) PathRetriever {
	return Join(append([]PathRetriever{os.UserHomeDir}, paths...)...)
}
