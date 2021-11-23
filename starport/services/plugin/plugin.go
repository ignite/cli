package plugin

import (
	"errors"
	"log"
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

// StarportPlugin provides interfaces for starport plugin.
type StarportPlugin interface {
	Execute(name string, args []string) error
	List() []FuncSpec

	// TODO: Need?
	// Cause single plugin can have multiple functions.
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
		log.Println(ErrSymbolNotExist.Error())
		return ErrSymbolNotExist
	}

	paramValues := make([]reflect.Value, len(spec.ParamTypes))
	for i, paramType := range spec.ParamTypes {
		val, err := convert(args[i], paramType)
		if err != nil {
			log.Println(err)
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

func (p *starportplugin) List() []FuncSpec {
	specs := make([]FuncSpec, len(p.funcSpecs))

	i := 0
	for _, v := range p.funcSpecs {
		specs[i] = v
		i++
	}

	return specs
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
