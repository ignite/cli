package jsondoc

import (
	"encoding/json"

	"github.com/goccy/go-yaml"
)

// Doc represents a JSON encoded data.
type Doc []byte

// MarshalYAML converts Doc to a YAML encoded data during YAML marshaling.
func (d Doc) MarshalYAML() ([]byte, error) {
	var out interface{}
	if err := json.Unmarshal(d, &out); err != nil {
		return nil, err
	}
	return yaml.Marshal(out)
}

// Pretty converts a Doc to a human readible string.
func (d Doc) Pretty() (string, error) {
	proposalyaml, err := yaml.Marshal(d)
	return string(proposalyaml), err
}
