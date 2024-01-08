package xos

import (
	"io"
	"os"
	"path/filepath"
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
