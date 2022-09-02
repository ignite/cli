package yaml

// Map defines a map type that uses strings as key value.
// The map implements the Unmarshaller interface to convert
// the unmershalled map keys type from interface{} to string.
type Map map[string]interface{}

func (m *Map) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var raw map[interface{}]interface{}

	if err := unmarshal(&raw); err != nil {
		return err
	}

	*m = convertMapKeys(raw)

	return nil
}

func convertMapKeys(raw map[interface{}]interface{}) map[string]interface{} {
	m := make(map[string]interface{})

	for k, v := range raw {
		if value, ok := v.(map[interface{}]interface{}); ok {
			v = convertMapKeys(value)
		}

		m[k.(string)] = v
	}

	return m
}
