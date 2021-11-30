package plugin

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"plugin"
	"reflect"

	"github.com/tendermint/starport/starport/chainconfig"
)

// Errors
var (
	ErrPluginWrongSpec = errors.New("plugin should follow basic specs")
)

// Loader provides managing features for plugin config.
type Loader interface {
	IsInstalled(config chainconfig.Plugin) bool
	LoadPlugin(config chainconfig.Plugin, pluginPath string) (StarportPlugin, error)
}

type configLoader struct {
	pluginSpec *starportplugin
}

// IsExists return true(bool) if given path does exist or doesn't
func (l *configLoader) IsExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Find return name of all files that does exist on given path, and given extension
func (l *configLoader) Find(root, ext string) []string {
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
	pluginsDirectory, _ := l.IsExists(pluginsPath)
	if pluginsDirectory {
		var pluginPath = filepath.Join(pluginsPath, plugin.Name)
		selectedPluginPath, _ := l.IsExists(pluginPath)
		if selectedPluginPath {
			fileList := l.Find(pluginPath, ".so")
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

	l.pluginSpec = &p

	err = l.checkMandatoryFunctions()
	if err != nil {
		log.Println(err)
		return nil, err
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

	return funcCallSpecs, err
}

func (l *configLoader) checkMandatoryFunctions() error {
	mandatories := []string{
		"Init",
		"Help",
	}

	marks := map[string]bool{}

	for _, v := range mandatories {
		marks[v] = false
	}

	for k, v := range l.pluginSpec.funcSpecs {
		marks[k] = true
		_ = v
	}

	for _, v := range marks {
		if !v {
			return ErrPluginWrongSpec
		}
	}

	// TODO: Check parameters.

	return nil
}

// NewLoader creates loader for plugin.
func NewLoader() (Loader, error) {
	return &configLoader{}, nil
}
