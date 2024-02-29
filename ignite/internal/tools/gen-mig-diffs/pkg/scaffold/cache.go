package scaffold

import (
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/ignite/cli/v28/ignite/pkg/errors"
)

// cache represents a cache for executed scaffold command.
type cache struct {
	cachePath  string
	cachesPath map[string]string
	mu         sync.RWMutex
}

// newCache initializes a new Cache instance.
func newCache(path string) (*cache, error) {
	return &cache{
		cachePath:  path,
		cachesPath: make(map[string]string),
	}, os.MkdirAll(path, os.ModePerm)
}

// save creates a new cache.
func (c *cache) save(name, path string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	dstPath := filepath.Join(c.cachePath, name)
	if err := copyFiles(path, dstPath); err != nil {
		return err
	}

	c.cachesPath[name] = dstPath
	return nil
}

// has return if the cache exist.
func (c *cache) has(name string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cachePath, ok := c.cachesPath[name]
	if !ok {
		return false
	}
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		return false
	}

	return true
}

// get return the cache path and copy all files to the destination path.
func (c *cache) get(name, dstPath string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cachePath, ok := c.cachesPath[name]
	if !ok {
		return errors.Errorf("command %s not exist in the cache list", name)
	}
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		return errors.Wrapf(err, "cache %s not exist in the path", name)
	}
	dstPath, err := filepath.Abs(dstPath)
	if err != nil {
		return err
	}
	if err := copyFiles(cachePath, dstPath); err != nil {
		return errors.Wrapf(err, "error to copy cache from %s to %s", cachePath, dstPath)
	}
	return nil
}

// copyFiles copy all files from the source path to the destination path.
func copyFiles(srcPath, dstPath string) error {
	srcInfo, err := os.Stat(srcPath)
	switch {
	case os.IsNotExist(err):
		return errors.Wrapf(err, "cache %s not exist in the path", srcPath)
	case err != nil:
		return err
	case !srcInfo.IsDir():
		return errors.Wrapf(err, "cache %s is not a directory", srcPath)
	}

	// Walk through the original path and copy all content to the cache path.
	return filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(srcPath, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dstPath, relPath)
		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()
		dstFile, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer dstFile.Close()
		_, err = io.Copy(dstFile, srcFile)
		return err
	})
}
