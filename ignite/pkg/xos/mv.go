package xos

import (
	"fmt"
	"io"
	"os"
)

// Rename copy oldPath to newPath and then delete oldPath.
// Unlike os.Rename, it doesn't fail when the oldPath and newPath are in
// different partitions (error: invalid cross-device link).
func Rename(oldPath, newPath string) error {
	inputFile, err := os.Open(oldPath)
	if err != nil {
		return fmt.Errorf("rename %s %s: couldn't open oldpath: %w", oldPath, newPath, err)
	}
	defer inputFile.Close()
	outputFile, err := os.Create(newPath)
	if err != nil {
		return fmt.Errorf("rename %s %s: couldn't open dest file: %w", oldPath, newPath, err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		return fmt.Errorf("rename %s %s: writing to output file failed: %w", oldPath, newPath, err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(oldPath)
	if err != nil {
		return fmt.Errorf("rename %s %s: failed removing original file: %w", oldPath, newPath, err)
	}
	return nil
}
