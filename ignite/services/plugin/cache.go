package plugin

import (
	"encoding/gob"
	"fmt"
	"net"
	"path"
	"path/filepath"

	hplugin "github.com/hashicorp/go-plugin"

	chainconfig "github.com/ignite/cli/ignite/config"
	"github.com/ignite/cli/ignite/pkg/cache"
)

const (
	cacheFileName  = "ignite_plugin_cache.db"
	cacheNamespace = "plugin.rpc.context"
)

var (
	storage      *cache.Storage
	storageCache *cache.Cache[ConfigContext]
)

func init() {
	gob.Register(hplugin.ReattachConfig{})
	gob.Register(net.UnixAddr{})
}

type ConfigContext struct {
	Plugin hplugin.ReattachConfig
	Addr   net.UnixAddr
}

func WritePluginConfig(pluginPath string, conf hplugin.ReattachConfig) error {
	name := path.Base(pluginPath)

	if name == "" {
		return fmt.Errorf("provided path is invalid: %s", pluginPath)
	}

	confCont := ConfigContext{}

	// TODO: figure out a better way of resolving the type of network connection is established between plugin server and host
	// currently this will always be a unix network socket. but this might not be the case moving forward.
	ua, err := net.ResolveUnixAddr(conf.Addr.Network(), conf.Addr.String())
	if err != nil {
		return err
	}

	confCont.Addr = *ua
	conf.Addr = nil
	confCont.Plugin = conf

	cache, err := newCache()
	if err != nil {
		return err
	}

	cache.Put(name, confCont)

	return err
}

func ReadPluginConfig(pluginPath string, ref *hplugin.ReattachConfig) error {
	name := path.Base(pluginPath)

	if name == "" {
		return fmt.Errorf("provided path is invalid: %s", pluginPath)
	}

	cache, err := newCache()
	if err != nil {
		return err
	}

	confCont, err := cache.Get(name)
	if err != nil {
		return err
	}

	*ref = confCont.Plugin
	ref.Addr = &confCont.Addr

	return nil
}

func CheckPluginConf(pluginPath string) bool {
	name := path.Base(pluginPath)

	if name == "." {
		return false
	}

	cache, err := newCache()
	if err != nil {
		return false
	}
	if _, err := cache.Get(name); err != nil {
		return false
	}
	return true
}

func DeletePluginConf(pluginPath string) error {
	name := path.Base(pluginPath)

	if name == "" {
		return fmt.Errorf("provided path is invalid: %s", pluginPath)
	}
	cache, err := newCache()
	if err != nil {
		return err
	}

	if err := cache.Delete(name); err != nil {
		return err
	}

	return nil
}

func newCache() (*cache.Cache[ConfigContext], error) {
	cacheRootDir, err := chainconfig.DirPath()
	if err != nil {
		return nil, err
	}
	if storage == nil {
		storageTmp, err := cache.NewStorage(filepath.Join(cacheRootDir, cacheFileName))
		if err != nil {
			return nil, err
		}
		storage = &storageTmp
		cacheTmp := cache.New[ConfigContext](*storage, cacheNamespace)
		storageCache = &cacheTmp
	}

	return storageCache, nil
}
