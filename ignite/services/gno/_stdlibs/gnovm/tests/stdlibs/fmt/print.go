package fmt

import (
	"fmt"
	"math"
	"sort"
	"unsafe"

	"github.com/gnolang/gno/gnovm/pkg/gnolang"
)

func X_typeString(v gnolang.TypedValue) string {
	if v.IsUndefined() {
		return "<nil>"
	}
	return v.T.String()
}

func X_valueOfInternal(v gnolang.TypedValue) (
	kind, declaredName string,
	bytes uint64,
	base gnolang.TypedValue,
	xlen int,
) {
	if v.IsUndefined() {
		kind = "nil"
		return
	}
	if dt, ok := v.T.(*gnolang.DeclaredType); ok {
		declaredName = dt.String()
	}
	baseT := gnolang.BaseOf(v.T)
	base = gnolang.TypedValue{
		T: baseT,
		V: v.V,
		N: v.N,
	}
	switch baseT.Kind() {
	case gnolang.BoolKind:
		kind = "bool"
		if v.GetBool() {
			bytes = 1
		}
	case gnolang.StringKind:
		kind, xlen = "string", v.GetLength()
	case gnolang.IntKind:
		kind, bytes = "int", uint64(v.GetInt())
	case gnolang.Int8Kind:
		kind, bytes = "int8", uint64(v.GetInt8())
	case gnolang.Int16Kind:
		kind, bytes = "int16", uint64(v.GetInt16())
	case gnolang.Int32Kind:
		kind, bytes = "int32", uint64(v.GetInt32())
	case gnolang.Int64Kind:
		kind, bytes = "int64", uint64(v.GetInt64())
	case gnolang.UintKind:
		kind, bytes = "uint", v.GetUint()
	case gnolang.Uint8Kind:
		kind, bytes = "uint8", uint64(v.GetUint8())
	case gnolang.Uint16Kind:
		kind, bytes = "uint16", uint64(v.GetUint16())
	case gnolang.Uint32Kind:
		kind, bytes = "uint32", uint64(v.GetUint32())
	case gnolang.Uint64Kind:
		kind, bytes = "uint64", v.GetUint64()
	case gnolang.Float32Kind:
		kind, bytes = "float32", uint64(v.GetFloat32())
	case gnolang.Float64Kind:
		kind, bytes = "float64", v.GetFloat64()
	case gnolang.ArrayKind:
		kind, xlen = "array", v.GetLength()
	case gnolang.SliceKind:
		kind, xlen = "slice", v.GetLength()
	case gnolang.PointerKind:
		kind = "pointer"
	case gnolang.StructKind:
		kind, xlen = "struct", len(baseT.(*gnolang.StructType).Fields)
	case gnolang.InterfaceKind:
		kind, xlen = "interface", len(baseT.(*gnolang.InterfaceType).Methods)
	case gnolang.FuncKind:
		kind = "func"
	case gnolang.MapKind:
		kind, xlen = "map", v.GetLength()
	default:
		panic("unexpected gnolang.Kind")
	}
	return
}

func X_getAddr(m *gnolang.Machine, v gnolang.TypedValue) uint64 {
	switch v.T.Kind() {
	case gnolang.FuncKind, gnolang.MapKind, gnolang.SliceKind, gnolang.PointerKind:
	default:
		panic("invalid type")
	}
	switch v := v.V.(type) {
	case nil:
		return 0
	case gnolang.PointerValue:
		if v.TV == nil {
			return 0
		}
		return uint64(uintptr(unsafe.Pointer(v.TV))) ^
			uint64(uintptr(unsafe.Pointer(&v.Base))) ^
			uint64(v.Index)
	case *gnolang.FuncValue:
		return uint64(uintptr(unsafe.Pointer(v)))
	case *gnolang.MapValue:
		return uint64(uintptr(unsafe.Pointer(v)))
	case *gnolang.SliceValue:
		return uint64(uintptr(unsafe.Pointer(v.GetBase(m.Store))))
	default:
		panic(fmt.Sprintf("unexpected value in getAddr: %T", v))
	}
}

func X_getPtrElem(v gnolang.TypedValue) gnolang.TypedValue {
	return v.V.(gnolang.PointerValue).Deref()
}

var gSliceOfAny = &gnolang.SliceType{
	Elt: &gnolang.InterfaceType{},
}

func X_mapKeyValues(v gnolang.TypedValue) (keys, values gnolang.TypedValue) {
	if v.T.Kind() != gnolang.MapKind {
		panic(fmt.Sprintf("invalid arg to mapKeyValues of kind: %s", v.T.Kind()))
	}
	keys.T = gSliceOfAny
	values.T = gSliceOfAny
	if v.V == nil {
		return
	}

	mv := v.V.(*gnolang.MapValue)
	ks, vs := make([]gnolang.TypedValue, 0, mv.GetLength()), make([]gnolang.TypedValue, 0, mv.GetLength())
	for el := mv.List.Head; el != nil; el = el.Next {
		ks = append(ks, el.Key)
		vs = append(vs, el.Value)
	}

	// use stable to maintain the same order when we have weird map keys,
	// like interfaces allowing for different concrete values.
	sort.Stable(mapKV{ks, vs})

	keys.V = &gnolang.SliceValue{
		Base: &gnolang.ArrayValue{
			List: ks,
		},
		Length: len(ks),
		Maxcap: len(ks),
	}
	values.V = &gnolang.SliceValue{
		Base: &gnolang.ArrayValue{
			List: vs,
		},
		Length: len(vs),
		Maxcap: len(vs),
	}
	return
}

// mapKV is the map key values for sorting them, similarly to internal/fmtsort.
// it implements sort.Interface, as it has to work on two different slices.
type mapKV struct {
	keys   []gnolang.TypedValue
	values []gnolang.TypedValue
}

func (m mapKV) Len() int { return len(m.keys) }
func (m mapKV) Swap(i, j int) {
	m.keys[i], m.keys[j] = m.keys[j], m.keys[i]
	m.values[i], m.values[j] = m.values[j], m.values[i]
}

func (m mapKV) Less(i, j int) bool {
	ki, kj := m.keys[i], m.keys[j]
	return compareKeys(ki, kj)
}

func compareKeys(ki, kj gnolang.TypedValue) bool {
	if ki.T == nil || kj.T == nil || ki.T.Kind() != kj.T.Kind() {
		return false
	}
	switch ki.T.Kind() {
	case gnolang.BoolKind:
		bi, bj := ki.GetBool(), kj.GetBool()
		// use == just to make it more explicit
		return bi == false && bj == true //nolint:staticcheck
	case gnolang.Float32Kind:
		return math.Float32frombits(ki.GetFloat32()) < math.Float32frombits(kj.GetFloat32())
	case gnolang.Float64Kind:
		return math.Float64frombits(ki.GetFloat64()) < math.Float64frombits(kj.GetFloat64())
	case gnolang.StringKind:
		return ki.GetString() < kj.GetString()
	case gnolang.IntKind,
		gnolang.Int8Kind,
		gnolang.Int16Kind,
		gnolang.Int32Kind,
		gnolang.Int64Kind:
		return ki.ConvertGetInt() < kj.ConvertGetInt()
	case gnolang.UintKind,
		gnolang.Uint8Kind,
		gnolang.Uint16Kind,
		gnolang.Uint32Kind,
		gnolang.Uint64Kind:
		return uint64(ki.ConvertGetInt()) < uint64(kj.ConvertGetInt())
	default:
		return false
	}
}

// get the n'th element of the given array or slice.
func X_arrayIndex(m *gnolang.Machine, v gnolang.TypedValue, n int) gnolang.TypedValue {
	switch v.T.Kind() {
	case gnolang.ArrayKind, gnolang.SliceKind:
		tv := gnolang.TypedValue{T: gnolang.IntType}
		tv.SetInt(int64(n))
		res := v.GetPointerAtIndex(m, m.Realm, m.Alloc, m.Store, &tv)
		return res.Deref()
	default:
		panic("invalid type to arrayIndex")
	}
}

// n'th field of the given struct.
func X_fieldByIndex(v gnolang.TypedValue, n int) (name string, value gnolang.TypedValue) {
	if v.T.Kind() != gnolang.StructKind {
		panic("invalid kind to fieldByIndex")
	}
	fldType := v.T.(*gnolang.StructType).Fields[n]
	name = string(fldType.Name)
	value.T = fldType.Type
	if v.V != nil {
		value = v.V.(*gnolang.StructValue).Fields[n]
	}
	return
}

func X_asByteSlice(v gnolang.TypedValue) (gnolang.TypedValue, bool) {
	switch {
	case v.T.Kind() == gnolang.SliceKind && v.T.Elem().Kind() == gnolang.Uint8Kind:
		return gnolang.TypedValue{
			T: &gnolang.SliceType{
				Elt: gnolang.Uint8Type,
			},
			V: v.V,
		}, true
	case v.T.Kind() == gnolang.ArrayKind && v.T.Elem().Kind() == gnolang.Uint8Kind:
		arrt := v.T.(*gnolang.ArrayType)
		return gnolang.TypedValue{
			T: &gnolang.SliceType{
				Elt: gnolang.Uint8Type,
			},
			V: &gnolang.SliceValue{
				Base:   v.V,
				Offset: 0,
				Length: arrt.Len,
				Maxcap: arrt.Len,
			},
		}, true
	default:
		return gnolang.TypedValue{}, false
	}
}
