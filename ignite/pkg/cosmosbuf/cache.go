package cosmosbuf

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ignite/cli/v29/ignite/pkg/cache"
	"github.com/ignite/cli/v29/ignite/pkg/dirchange"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/xos"
)

func cacheKey(src, file, template string) (string, error) {
	relPath, err := filepath.Rel(src, file)
	if err != nil {
		return "", err
	}

	checksum, err := dirchange.ChecksumFromPaths(src, relPath)
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
	key := fmt.Sprintf("%x", h.Sum(nil))
	return key, nil
}

func (b Buf) getFileCache(src, file, template, output string) (bool, error) {
	key, err := cacheKey(src, file, template)
	if err != nil {
		return false, err
	}

	existingFile, err := b.storageCache.Get(key)
	if errors.Is(err, cache.ErrorNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	relPath, err := filepath.Rel(src, file)
	if err != nil {
		return false, err
	}

	filePath := filepath.Join(output, relPath)
	if err := os.WriteFile(filePath, existingFile, 0o644); err != nil {
		return false, err
	}
	return true, nil
}

func (b Buf) getDirCache(src, output, template string) (map[string]struct{}, error) {
	result := make(map[string]struct{})
	files, err := xos.FindFiles(src)
	if err != nil {
		return result, err
	}

	for _, file := range files {
		ok, err := b.getFileCache(src, file, template, output)
		if err != nil {
			return result, err
		}
		if ok {
			result[file] = struct{}{}
		}
	}
	return result, nil
}

func (b Buf) saveFileCache(src, file, template string) error {
	key, err := cacheKey(src, file, template)
	if err != nil {
		return err
	}

	f, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return b.storageCache.Put(key, f)
}

func (b Buf) saveDirCache(src, template string) error {
	files, err := xos.FindFiles(src)
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := b.saveFileCache(src, file, template); err != nil {
			return err
		}
	}
	return nil
}
