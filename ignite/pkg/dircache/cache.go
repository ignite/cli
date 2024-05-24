package dircache

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"

	"github.com/otiai10/copy"

	"github.com/ignite/cli/v29/ignite/config"
	"github.com/ignite/cli/v29/ignite/pkg/cache"
	"github.com/ignite/cli/v29/ignite/pkg/dirchange"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

var ErrCacheNotFound = errors.New("cache not found")

type Cache struct {
	path         string
	storageCache cache.Cache[string]
}

// New creates a new Buf based on the installed binary.
func New(cacheStorage cache.Storage, dir, specNamespace string) (Cache, error) {
	path, err := cachePath()
	if err != nil {
		return Cache{}, err
	}
	path = filepath.Join(path, dir)
	if err := os.MkdirAll(path, 0o755); err != nil && !os.IsExist(err) {
		return Cache{}, err
	}

	return Cache{
		path:         path,
		storageCache: cache.New[string](cacheStorage, specNamespace),
	}, nil
}

// ClearCache remove the cache path.
func ClearCache() error {
	path, err := cachePath()
	if err != nil {
		return err
	}
	return os.RemoveAll(path)
}

// cachePath returns the cache path.
func cachePath() (string, error) {
	globalPath, err := config.DirPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(globalPath, "cache"), nil
}

// cacheKey create the cache key.
func cacheKey(src string, keys ...string) (string, error) {
	checksum, err := dirchange.ChecksumFromPaths(src, "")
	if err != nil {
		return "", err
	}

	h := sha256.New()
	if _, err := h.Write(checksum); err != nil {
		return "", err
	}
	for _, key := range keys {
		if _, err := h.Write([]byte(key)); err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// CopyTo gets the cache folder based on the cache key from the storage and copies the folder to the output.
func (c Cache) CopyTo(src, output string, keys ...string) (string, error) {
	key, err := cacheKey(src, keys...)
	if err != nil {
		return key, err
	}

	cachedPath, err := c.storageCache.Get(key)
	if errors.Is(err, cache.ErrorNotFound) {
		return key, ErrCacheNotFound
	} else if err != nil {
		return key, err
	}

	if err := copy.Copy(cachedPath, output); err != nil {
		return "", errors.Wrapf(err, "get dir cache cannot copy path %s to %s", cachedPath, output)
	}
	return key, nil
}

// Save copies the source to the cache folder and saves the path into the storage based on the key.
func (c Cache) Save(src, key string) error {
	path := filepath.Join(c.path, key)
	if err := os.Mkdir(path, 0o700); os.IsExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	if err := copy.Copy(src, path); err != nil {
		return errors.Wrapf(err, "save dir cache cannot copy path %s to %s", src, path)
	}
	return c.storageCache.Put(key, path)
}
