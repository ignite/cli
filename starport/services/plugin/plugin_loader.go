package plugin

import (
	"fmt"
	"log"
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

// IsInstalled checks whether the plugins are installed.
func (l *configLoader) IsInstalled(config chainconfig.Plugin) bool {
	// TODO: jkkim: Check plugin file exist on home.

	return false
}

func (l *configLoader) LoadPlugin(config chainconfig.Plugin, pluginPath string) (StarportPlugin, error) {
	// TODO: jkkim: How do I get path of real ".so" symbol?
	// configDir, err := chainconfig.ConfigDirPath()
	// if err != nil {
	// 	return nil, err
	// }

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
