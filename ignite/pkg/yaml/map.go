package yaml

// Map defines a map type that uses strings as key value.
// The map implements the Unmarshaller interface to convert
// the unmarshalled map keys type from interface{} to string.
type Map map[string]interface{}

func (m *Map) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var raw map[interface{}]interface{}

	if err := unmarshal(&raw); err != nil {
		return err
	}

	*m = convertMapKeys(raw)

	return nil
}

func convertSlice(raw []interface{}) []interface{} {
	if len(raw) == 0 {
		return raw
	}

	if _, ok := raw[0].(map[interface{}]interface{}); !ok {
		return raw
	}

	values := make([]interface{}, len(raw))
	for i, v := range raw {
		values[i] = convertMapKeys(v.(map[interface{}]interface{}))
	}

	return values
}

func convertMapKeys(raw map[interface{}]interface{}) map[string]interface{} {
	m := make(map[string]interface{})

	for k, v := range raw {
		if value, _ := v.(map[interface{}]interface{}); value != nil {
			// Convert map keys to string
			v = convertMapKeys(value)
		} else if values, _ := v.([]interface{}); values != nil {
			// Make sure that maps inside slices also use strings as key
			v = convertSlice(values)
		}

		m[k.(string)] = v
	}

	return m
}
