package xyaml

// Map defines a map type that uses strings as key value.
// The map implements the Unmarshaller interface to convert
// the unmarshalled map keys type from interface{} to string.
type Map map[string]any

func (m *Map) UnmarshalYAML(unmarshal func(any) error) error {
	var raw map[any]any

	if err := unmarshal(&raw); err != nil {
		return err
	}

	*m = convertMapKeys(raw)

	return nil
}

func convertSlice(raw []any) []any {
	if len(raw) == 0 {
		return raw
	}

	if _, ok := raw[0].(map[any]any); !ok {
		return raw
	}

	values := make([]any, len(raw))
	for i, v := range raw {
		values[i] = convertMapKeys(v.(map[any]any))
	}

	return values
}

func convertMapKeys(raw map[any]any) map[string]any {
	m := make(map[string]any)

	for k, v := range raw {
		if value, _ := v.(map[any]any); value != nil {
			// Convert map keys to string
			v = convertMapKeys(value)
		} else if values, _ := v.([]any); values != nil {
			// Make sure that maps inside slices also use strings as key
			v = convertSlice(values)
		}

		m[k.(string)] = v
	}

	return m
}
