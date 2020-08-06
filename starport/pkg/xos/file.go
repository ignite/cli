package xos

import "os"

// OpenFirst finds and opens the first found file within names.
func OpenFirst(names ...string) (file *os.File, err error) {
	for _, name := range names {
		file, err = os.Open(name)
		if err == nil {
			break
		}
	}
	return file, err
}
