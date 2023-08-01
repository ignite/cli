package dirchange

import (
	"bytes"
	"crypto/md5"
	"errors"
	"os"
	"path/filepath"

	"github.com/ignite/cli/ignite/pkg/cache"
)

var ErrNoFile = errors.New("no file in specified paths")

// SaveDirChecksum saves the md5 checksum of the provided paths (directories or files) in the provided cache.
// If checksumSavePath directory doesn't exist, it is created.
// Paths are relative to workdir. If workdir is empty, string paths are absolute.
func SaveDirChecksum(checksumCache cache.Cache[[]byte], cacheKey string, workdir string, paths ...string) error {
	checksum, err := ChecksumFromPaths(workdir, paths...)
	if err != nil {
		return err
	}

	// save checksum
	return checksumCache.Put(cacheKey, checksum)
}

// HasDirChecksumChanged computes the md5 checksum of the provided paths (directories or files)
// and compares it with the current cached checksum.
// Return true if the checksum doesn't exist yet.
// paths are relative to workdir, if workdir is empty string paths are absolute.
func HasDirChecksumChanged(checksumCache cache.Cache[[]byte], cacheKey string, workdir string, paths ...string) (bool, error) {
	savedChecksum, err := checksumCache.Get(cacheKey)
	if errors.Is(err, cache.ErrorNotFound) {
		return true, nil
	}
	if err != nil {
		return false, err
	}

	// Compute checksum
	checksum, err := ChecksumFromPaths(workdir, paths...)
	if errors.Is(err, ErrNoFile) {
		// Checksum cannot be saved with no file
		// Therefore if no file are found, this means these have been deleted, then the directory has been changed
		return true, nil
	} else if err != nil {
		return false, err
	}

	// Compare checksums
	if bytes.Equal(checksum, savedChecksum) {
		return false, nil
	}

	// The checksum has changed
	return true, nil
}

// ChecksumFromPaths computes the md5 checksum from the provided paths.
// Relative paths to the workdir are used. If workdir is empty, string paths are absolute.
func ChecksumFromPaths(workdir string, paths ...string) ([]byte, error) {
	hash := md5.New()

	// Can't compute hash if no file present
	noFile := true

	// read files
	for _, path := range paths {
		if !filepath.IsAbs(path) {
			path = filepath.Join(workdir, path)
		}

		// non-existent paths are ignored
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		} else if err != nil {
			return []byte{}, err
		}

		err := filepath.Walk(path, func(subPath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// ignore directory
			if info.IsDir() {
				return nil
			}

			noFile = false

			// write file content
			content, err := os.ReadFile(subPath)
			if err != nil {
				return err
			}
			_, err = hash.Write(content)
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return []byte{}, err
		}
	}

	if noFile {
		return []byte{}, ErrNoFile
	}

	// compute checksum
	return hash.Sum(nil), nil
}
