package cache

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/xos"
)

// Cache represents a cache for executed scaffold command.
type Cache struct {
	cachePath  string
	cachesPath map[string]string
	mu         sync.RWMutex
}

// New initializes a new Cache instance.
func New(path string) (*Cache, error) {
	return &Cache{
		cachePath:  path,
		cachesPath: make(map[string]string),
	}, os.MkdirAll(path, os.ModePerm)
}

// Save creates a new cache.
func (c *Cache) Save(name, path string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	dstPath := filepath.Join(c.cachePath, name)
	if err := xos.CopyFolder(path, dstPath); err != nil {
		return err
	}

	c.cachesPath[name] = dstPath
	return nil
}

// Has return if the cache exist.
func (c *Cache) Has(name string) bool {
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

// Get return the cache path and copy all files to the destination path.
func (c *Cache) Get(name, dstPath string) error {
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
	if err := xos.CopyFolder(cachePath, dstPath); err != nil {
		return errors.Wrapf(err, "error to copy cache from %s to %s", cachePath, dstPath)
	}
	return nil
}
