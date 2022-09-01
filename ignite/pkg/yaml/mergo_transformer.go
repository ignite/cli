package yaml

import "reflect"

// MapIndexTransformer defines a transformer for the "Mergo" project to be able to merge
// maps that were generated from a YAML file. These maps usually have an interface{}
// instead of a string type as keys even when the key values are actually strings.
// The transformer can be applied during a merge using mergo's WithTransformers option
// as `mergo.WithTransformers(MapIndexTransformer{})`.
type MapIndexTransformer struct{}

func (t MapIndexTransformer) Transformer(dstType reflect.Type) func(dst, src reflect.Value) error {
	// Make sure that transformation only applies to maps
	if dstType.Kind() != reflect.Map {
		return nil
	}

	return func(dst, src reflect.Value) error {
		if dst.CanSet() {
			iter := dst.MapRange()
			for iter.Next() {
				// Check if the key type is interface{}
				k := iter.Key()
				if k.Kind() != reflect.Interface {
					continue
				}

				// When the key is an interface{} cast its value to string and reasign it
				if s, ok := k.Interface().(string); ok {
					dst.SetMapIndex(reflect.ValueOf(s), iter.Value())
				}
			}
		}

		return nil
	}
}
