package xos

import (
	"io"
	"os"
	"path/filepath"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

// CopyFolder copy the source folder to the destination folder.
func CopyFolder(srcPath, dstPath string) error {
	return filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root folder
		if path == srcPath {
			return nil
		}

		// Get the relative path within the source folder
		relativePath, err := filepath.Rel(srcPath, path)
		if err != nil {
			return err
		}

		// Create the corresponding destination path
		destPath := filepath.Join(dstPath, relativePath)

		if info.IsDir() {
			// Create the directory in the destination
			err = os.MkdirAll(destPath, 0o755)
			if err != nil {
				return err
			}
		} else {
			// Copy the file content
			err = CopyFile(path, destPath)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// ValidateFolderCopy validates that all files in source folder exist in destination folder
// with same name and relative path.
func ValidateFolderCopy(srcPath, dstPath string, exclude ...string) ([]string, error) {
	if srcPath == dstPath {
		return nil, errors.Errorf("source and destination paths are the same %s", srcPath)
	}

	// Check if the destination path exists
	if _, err := os.Stat(dstPath); errors.Is(err, os.ErrNotExist) {
		return nil, errors.Errorf("destination path does not exist: %s", dstPath)
	} else if err != nil {
		return nil, err
	}

	excludeMap := make(map[string]struct{}, len(exclude))
	for _, ex := range exclude {
		excludeMap[ex] = struct{}{}
	}

	var sameFiles []string
	err := filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
		if errors.Is(err, os.ErrNotExist) {
			return errors.Errorf("source path does not exist: %s", path)
		}
		if err != nil {
			return err
		}

		// Skip dirs
		if info.IsDir() {
			return nil
		}

		// Get the relative path within the source folder
		relativePath, err := filepath.Rel(srcPath, path)
		if err != nil {
			return err
		}

		// Skip excluded files
		if _, ok := excludeMap[relativePath]; ok {
			return nil
		}

		// Create the corresponding destination path
		destPath := filepath.Join(dstPath, relativePath)

		// Check if the destination path exists
		destInfo, err := os.Stat(destPath)
		if os.IsNotExist(err) {
			return nil
		} else if err != nil {
			return err
		}

		// Verify if directory/file types match
		if info.IsDir() != destInfo.IsDir() {
			return os.ErrInvalid
		}

		sameFiles = append(sameFiles, relativePath)
		return nil
	})
	return sameFiles, err
}

// CopyFile copy the source file to the destination file.
func CopyFile(srcPath, dstPath string) error {
	srcFile, err := os.OpenFile(srcPath, os.O_RDONLY, 0o666)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	return err
}
