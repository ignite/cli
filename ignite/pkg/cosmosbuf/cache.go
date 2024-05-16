package cosmosbuf

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

func ClearCache() error {
	path, err := cachePath()
	if err != nil {
		return err
	}
	return os.RemoveAll(path)
}

func cachePath() (string, error) {
	globalPath, err := config.DirPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(globalPath, "buf"), nil
}

func cacheKey(src, template string) (string, error) {
	checksum, err := dirchange.ChecksumFromPaths(src, "")
	if err != nil {
		return "", err
	}

	h := sha256.New()
	if _, err := h.Write(checksum); err != nil {
		return "", err
	}
	if _, err := h.Write([]byte(template)); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func (b Buf) copyCache(src, template, output string) (string, error) {
	key, err := cacheKey(src, template)
	if err != nil {
		return key, err
	}

	cachedPath, err := b.storageCache.Get(key)
	if errors.Is(err, cache.ErrorNotFound) {
		return key, ErrCacheNotFound
	} else if err != nil {
		return key, err
	}

	if err := copy.Copy(cachedPath, output); err != nil {
		return "", errors.Wrapf(err, "buf get cache cannot copy path %s to %s", cachedPath, output)
	}
	return key, nil
}

func (b Buf) saveCache(key, src string) error {
	cachePath := filepath.Join(b.bufCachePath, key)
	if err := os.Mkdir(cachePath, 0o700); os.IsExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	if err := copy.Copy(src, cachePath); err != nil {
		return errors.Wrapf(err, "buf save cache cannot copy path %s to %s", src, cachePath)
	}
	return b.storageCache.Put(key, cachePath)
}
