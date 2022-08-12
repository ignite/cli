package xos

import (
	"fmt"
	"io"
	"os"
)

// Rename copy sourcePath to destPath and then delete sourcePath.
// Unlike os.Rename, it doesn't fail when the oldpath and newpath are in
// different partitions (error: invalid cross-device link).
func Rename(oldpath, newpath string) error {
	inputFile, err := os.Open(oldpath)
	if err != nil {
		return fmt.Errorf("rename %s %s: couldn't open oldpath: %w", oldpath, newpath, err)
	}
	defer inputFile.Close()
	outputFile, err := os.Create(newpath)
	if err != nil {
		return fmt.Errorf("rename %s %s: couldn't open dest file: %w", oldpath, newpath, err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		return fmt.Errorf("rename %s %s: writing to output file failed: %w", oldpath, newpath, err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(oldpath)
	if err != nil {
		return fmt.Errorf("rename %s %s: failed removing original file: %w", oldpath, newpath, err)
	}
	return nil
}
