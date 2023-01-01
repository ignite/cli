package plugin

import (
	"encoding/gob"
	"fmt"
	"net"
	"path"

	hplugin "github.com/hashicorp/go-plugin"

	"github.com/ignite/cli/ignite/pkg/cache"
)

const (
	cacheFileName  = "ignite_plugin_cache.db"
	cacheNamespace = "plugin.rpc.context"
)

var storageCache *cache.Cache[hplugin.ReattachConfig]

func init() {
	gob.Register(hplugin.ReattachConfig{})
	gob.Register(&net.UnixAddr{})
}

func writeConfigCache(pluginPath string, conf hplugin.ReattachConfig) error {
	if pluginPath == "" {
		return fmt.Errorf("provided path is invalid: %s", pluginPath)
	}
	if conf.Addr == nil {
		return fmt.Errorf("plugin Address info cannot be empty")
	}
	cache, err := newCache()
	if err != nil {
		return err
	}
	return cache.Put(pluginPath, conf)
}

func readConfigCache(pluginPath string) (hplugin.ReattachConfig, error) {
	if pluginPath == "" {
		return hplugin.ReattachConfig{}, fmt.Errorf("provided path is invalid: %s", pluginPath)
	}
	cache, err := newCache()
	if err != nil {
		return hplugin.ReattachConfig{}, err
	}
	return cache.Get(pluginPath)
}

func checkConfCache(pluginPath string) bool {
	if pluginPath == "" {
		return false
	}
	cache, err := newCache()
	if err != nil {
		return false
	}
	_, err = cache.Get(pluginPath)
	return err == nil
}

func deleteConfCache(pluginPath string) error {
	if pluginPath == "" {
		return fmt.Errorf("provided path is invalid: %s", pluginPath)
	}
	cache, err := newCache()
	if err != nil {
		return err
	}
	return cache.Delete(pluginPath)
}

func newCache() (*cache.Cache[hplugin.ReattachConfig], error) {
	cacheRootDir, err := PluginsPath()
	if err != nil {
		return nil, err
	}
	if storageCache == nil {
		storage, err := cache.NewStorage(path.Join(cacheRootDir, cacheFileName))
		if err != nil {
			return nil, err
		}
		c := cache.New[hplugin.ReattachConfig](storage, cacheNamespace)
		storageCache = &c
	}
	return storageCache, nil
}
