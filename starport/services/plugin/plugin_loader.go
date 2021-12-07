package plugin

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"plugin"
	"reflect"
	"strings"

	"github.com/tendermint/starport/starport/chainconfig"
)

var (
	mandatories = map[string][]reflect.Kind{
		"Init": {},
		"Help": {reflect.String},
	}
)

// Loader provides managing features for plugin config.
type Loader interface {
	IsInstalled(config chainconfig.Plugin) bool
	LoadPlugin(config chainconfig.Plugin, pluginPath string) (StarportPlugin, error)

	LoadSymbol(symbol string) (map[string]FuncSpec, error)
}

type configLoader struct {
	chainID    string
	pluginSpec *starportplugin
}

// IsExists return true(bool) if given path does exist or doesn't
func (l *configLoader) IsExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return true, err
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
	defaultPath, err := chainconfig.ConfigDirPath()
	if err != nil {
		panic(err)
	}

	pluginsHome := filepath.Join(defaultPath, "plugins")

	tokens := strings.Split(plugin.RepositoryURL, "/")
	repoName := tokens[len(tokens)-1]

	pluginPath := fmt.Sprintf("%s/%s/%s/%s", pluginsHome, l.chainID, repoName, plugin.Name)

	isExist, err := l.IsExists(pluginPath)
	if err != nil {
		return false
	}

	if isExist {
		libList := l.Find(pluginPath, ".so")
		if len(libList) > 0 {
			return true
		}
	}
	return false
}

func (l *configLoader) LoadPlugin(config chainconfig.Plugin, pluginPath string) (StarportPlugin, error) {
	tokens := strings.Split(config.RepositoryURL, "/")
	repoName := tokens[len(tokens)-1]

	pluginSymbol := fmt.Sprintf("%s/%s/%s/%s/%s.so", pluginPath, l.chainID, repoName, config.Name, config.Name)
	specs, err := l.LoadSymbol(pluginSymbol)
	if err != nil {
		return nil, err
	}

	p := starportplugin{
		name:      config.Name,
		funcSpecs: specs,
	}

	l.pluginSpec = &p

	return &p, nil
}

func (l *configLoader) LoadSymbol(symbolName string) (map[string]FuncSpec, error) {
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

	err = l.checkMandatoryFunctions(funcCallSpecs)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return funcCallSpecs, err
}

func (l *configLoader) checkMandatoryFunctions(spec map[string]FuncSpec) error {
	for funcName, paramTypes := range mandatories {
		loadSpec, ok := spec[funcName]
		if !ok {
			log.Println("Not exist func ", funcName)
			return ErrPluginWrongSpec
		}

		if len(loadSpec.ParamTypes) != len(paramTypes) {
			log.Println("Invalid param ", funcName)
			return ErrPluginWrongSpec
		}

		for i, t := range paramTypes {
			if loadSpec.ParamTypes[i].Kind() != t {
				return ErrPluginWrongSpec
			}
		}
	}

	return nil
}

// NewLoader creates loader for plugin.
func NewLoader(chainID string) (Loader, error) {
	return &configLoader{chainID: chainID}, nil
}
