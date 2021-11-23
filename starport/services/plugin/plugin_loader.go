package plugin

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"plugin"
	"reflect"

	"github.com/tendermint/starport/starport/chainconfig"
)

// Loader provides managing features for plugin config.
type Loader interface {
	IsInstalled(config chainconfig.Plugin) bool
	LoadPlugin(config chainconfig.Plugin, pluginPath string) (StarportPlugin, error)
}

type configLoader struct {
}

func IsExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func find(root, ext string) []string {
	var a []string
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			a = append(a, s)
		}
		return nil
	})
	return a
}

// IsInstalled checks whether the given plugin is installed.
func (l *configLoader) IsInstalled(plugin chainconfig.Plugin) bool {
	isExists := false
	// TODO: D.K: Check plugin file exist on home.
	defaultPath, _ := chainconfig.ConfigDirPath()
	var pluginsPath = filepath.Join(defaultPath, "plugins")
	pluginsDirectory, _ := IsExists(pluginsPath)
	if pluginsDirectory {
		var pluginPath = filepath.Join(pluginsPath, plugin.Name)
		selectedPluginPath, _ := IsExists(pluginPath)
		if selectedPluginPath {
			fileList := find(pluginPath, ".so")
			if len(fileList) > 0 {
				isExists = true
			}
		}
	}
	return isExists
}

func (l *configLoader) LoadPlugin(config chainconfig.Plugin, pluginPath string) (StarportPlugin, error) {
	pluginSymbol := fmt.Sprintf("%s/%s/%s.so", pluginPath, config.Name, config.Name)
	specs, err := l.loadSymbol(pluginSymbol)
	if err != nil {
		return nil, err
	}

	p := starportplugin{
		name:      config.Name,
		funcSpecs: specs,
	}

	return &p, nil
}

func (l *configLoader) loadSymbol(symbolName string) (map[string]FuncSpec, error) {
	p, err := plugin.Open(symbolName)
	if err != nil {
		log.Println(err)
		return nil, ErrSymbolNotExist
	}

	sym, err := p.Lookup(PluginSymbolName)
	if err != nil {
		return nil, err
	}

	funcCallSpecs := map[string]FuncSpec{}

	symType := reflect.TypeOf(sym)
	symVal := reflect.ValueOf(sym)

	for i := 0; i < symType.NumMethod(); i++ {
		method := symType.Method(i)

		callSpec := FuncSpec{
			Name:       method.Name,
			ParamTypes: make([]reflect.Type, method.Type.NumIn()-1),
			Func:       symVal.Method(i),
		}

		// First element of parameter is self instance. ignore it.
		for j := 0; j < method.Type.NumIn()-1; j++ {
			callSpec.ParamTypes[j] = method.Type.In(j + 1)
		}

		funcCallSpecs[method.Name] = callSpec
	}

	// TODO: jkkim: Check mandatory functions.

	return funcCallSpecs, nil
}

// NewLoader creates loader for plugin.
func NewLoader() (Loader, error) {
	return &configLoader{}, nil
}
