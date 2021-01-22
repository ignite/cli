package dirchange

import (
	"bytes"
	"crypto/md5"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	checksumFile = "source_checksum.txt"
)

// SaveDirChecksum saves the md5 checksum of the provided paths (directories or files) in the specified directory
// If checksumSavePath directory doesn't exist, it is created
// paths are relative to workdir, if workdir is empty string paths are absolute
func SaveDirChecksum(workdir string, paths []string, checksumSavePath string) error {
	checksum, err := checksumFromPaths(workdir, paths)
	if err != nil {
		return err
	}

	// create directory if needed
	if err := os.MkdirAll(checksumSavePath, 0700); err != nil && !os.IsExist(err) {
		return err
	}

	// save checksum
	checksumFilePath := filepath.Join(checksumSavePath, checksumFile)
	return ioutil.WriteFile(checksumFilePath, checksum, 0644)
}

// HasDirChecksumChanged computes the md5 checksum of the provided paths (directories or files)
// and compare it with the current saved checksum
// If the checksum is different, the new checksum is saved
// Return true if the checksum file doesn't exist yet and if checksumSavePath directory doesn't exist, it is created
// paths are relative to workdir, if workdir is empty string paths are absolute
func HasDirChecksumChanged(workdir string, paths []string, checksumSavePath string) (bool, error) {
	checksum, err := checksumFromPaths(workdir, paths)
	if err != nil {
		return false, err
	}

	// create directory if needed
	if err := os.MkdirAll(checksumSavePath, 0700); err != nil && !os.IsExist(err) {
		return false, err
	}

	checksumFilePath := filepath.Join(checksumSavePath, checksumFile)
	if _, err := os.Stat(checksumFilePath); os.IsNotExist(err) {
		return true, ioutil.WriteFile(checksumFilePath, checksum, 0644)
	}

	// Compare checksums
	savedChecksum, err := ioutil.ReadFile(checksumFilePath)
	if err != nil {
		return false, err
	}
	if bytes.Equal(checksum, savedChecksum) {
		return false, nil
	}

	// The checksum has changed
	return true, ioutil.WriteFile(checksumFilePath, checksum, 0644)
}

// checksumFromPaths computes the md5 checksum from the provided paths
// paths are relative to workdir, if workdir is empty string paths are absolute
func checksumFromPaths(workdir string, paths []string) ([]byte, error) {
	hash := md5.New()

	// Can't compute hash if no file present
	noFile := true

	// read files
	for _, path := range paths {
		if workdir != "" {
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
			content, err := ioutil.ReadFile(subPath)
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
		return []byte{}, errors.New("no file in specified paths")
	}

	// compute checksum
	return hash.Sum(nil), nil
}
