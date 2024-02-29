package scaffold

import (
	"io"
	"os"
	"path/filepath"
	"sync"
)

// cache represents a cache for executed scaffold commandList.
type cache struct {
	cachePath  string
	cachesPath map[string]string
	mu         sync.RWMutex
}

// newCache initializes a new Cache instance.
func newCache(path string) *cache {
	return &cache{
		cachePath:  path,
		cachesPath: make(map[string]string),
	}
}

// saveCache save a new cache
func (c *cache) saveCache(name, path string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	cachePath := filepath.Join(path, name)
	// Walk through the original path and copy all content to the cache path.
	err := filepath.Walk(path, func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(path, srcPath)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(cachePath, relPath)
		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}
		srcFile, err := os.Open(srcPath)
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
	if err != nil {
		return err
	}

	c.cachesPath[name] = cachePath
	return nil
}

// get return the cache path
func (c *cache) get(name string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	cachePath, ok := c.cachesPath[name]
	if !ok {
		return "", false
	}
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		return "", false
	}
	return cachePath, true
}
