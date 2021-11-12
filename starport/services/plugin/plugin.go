package plugin

import (
	"errors"
	"fmt"
	"log"
	"plugin"
	"reflect"
	"strconv"
)

// Constants
const (
	PluginSymbolName = "Plugin"
)

// Errors
var (
	ErrSymbolNotExist = errors.New("not exist symbol")
	ErrNotInitilized  = errors.New("not initialized")
)

// Plugin provides interfaces for starport plugin.
type Plugin interface {
	Execute(name string, args []string) error
	List() error
	Usage(name string) error

	Name() string
}

type starportplugin struct {
	name      string
	funcSpecs map[string]FuncSpec
}

func (p *starportplugin) Execute(name string, args []string) error {
	spec, ok := p.funcSpecs[name]
	if !ok {
		return ErrSymbolNotExist
	}

	paramValues := make([]reflect.Value, len(spec.ParamTypes))
	for i, paramType := range spec.ParamTypes {
		val, err := convert(args[i], paramType)
		if err != nil {
			return err
		}

		paramValues[i] = val
	}

	// TODO: Any err required?
	_ = spec.Func.Call(paramValues)

	return nil
}

func convert(in string, expectType reflect.Type) (reflect.Value, error) {
	switch expectType.Kind() {
	case reflect.Bool:
		v, err := strconv.ParseBool(in)
		return reflect.ValueOf(v), err

	case reflect.Int:
		v, err := strconv.ParseInt(in, 10, 64)
		return reflect.ValueOf(int(v)), err

	case reflect.Int8:
		v, err := strconv.ParseInt(in, 10, 8)
		return reflect.ValueOf(int8(v)), err

	case reflect.Int16:
		v, err := strconv.ParseInt(in, 10, 16)
		return reflect.ValueOf(int16(v)), err

	case reflect.Int32:
		v, err := strconv.ParseInt(in, 10, 32)
		return reflect.ValueOf(int32(v)), err

	case reflect.Int64:
		v, err := strconv.ParseInt(in, 10, 64)
		return reflect.ValueOf(v), err

	case reflect.Uint8:
		v, err := strconv.ParseUint(in, 10, 8)
		return reflect.ValueOf(uint8(v)), err

	case reflect.Uint16:
		v, err := strconv.ParseUint(in, 10, 16)
		return reflect.ValueOf(uint16(v)), err

	case reflect.Uint32:
		v, err := strconv.ParseUint(in, 10, 32)
		return reflect.ValueOf(uint32(v)), err

	case reflect.Uint64:
		v, err := strconv.ParseUint(in, 10, 64)
		return reflect.ValueOf(v), err

	case reflect.Float32:
		v, err := strconv.ParseFloat(in, 32)
		return reflect.ValueOf(float32(v)), err

	case reflect.Float64:
		v, err := strconv.ParseFloat(in, 64)
		return reflect.ValueOf(v), err

	case reflect.Complex64:
		v, err := strconv.ParseComplex(in, 64)
		return reflect.ValueOf(complex64(v)), err

	case reflect.Complex128:
		v, err := strconv.ParseComplex(in, 128)
		return reflect.ValueOf(v), err

	case reflect.String:
		return reflect.ValueOf(in), nil

	default:
		return reflect.ValueOf(1), nil
	}
}

func (p *starportplugin) List() error {
	// TODO:

	if len(p.funcSpecs) == 0 {
		return ErrNotInitilized
	}

	for k, spec := range p.funcSpecs {
		fmt.Printf("%+v %+v\n", k, spec)
	}

	return nil
}

func (p *starportplugin) Usage(name string) error {
	// TODO: How to provide help?
	return nil
}

func (p *starportplugin) Name() string {
	return p.name
}

// FuncSpec describes function spec of reflection to be called.
type FuncSpec struct {
	Name       string
	ParamTypes []reflect.Type
	Func       reflect.Value
}

// LoadPlugin loads received parametered symbols and creates plugin instance.
// This function will injected into PluginLoader.
func LoadPlugin(name string, pluginSymbol string) (Plugin, error) {
	specs, err := loadSymbol(pluginSymbol)
	if err != nil {
		return nil, err
	}

	p := starportplugin{
		name:      name,
		funcSpecs: specs,
	}

	return &p, nil
}

func loadSymbol(symbolName string) (map[string]FuncSpec, error) {
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

	// TODO: Check mandatory functions.

	return funcCallSpecs, nil
}
