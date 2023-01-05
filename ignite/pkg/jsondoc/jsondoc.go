package jsondoc

import (
	"encoding/json"

	"github.com/goccy/go-yaml"
)

// Doc represents a JSON encoded data.
type Doc []byte

// ToDocs converts a list of JSON encoded data to docs.
func ToDocs(data [][]byte) []Doc {
	var docs []Doc
	for _, d := range data {
		docs = append(docs, d)
	}
	return docs
}

// MarshalYAML converts Doc to a YAML encoded data during YAML marshaling.
func (d Doc) MarshalYAML() ([]byte, error) {
	var out interface{}
	if err := json.Unmarshal(d, &out); err != nil {
		return nil, err
	}
	return yaml.Marshal(out)
}

// Pretty converts a Doc to a human-readable string.
func (d Doc) Pretty() (string, error) {
	proposalyaml, err := yaml.Marshal(d)
	return string(proposalyaml), err
}
