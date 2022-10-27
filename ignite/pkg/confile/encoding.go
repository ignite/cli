package confile

import (
	"encoding/json"
	"io"

	"github.com/goccy/go-yaml"
	"github.com/pelletier/go-toml"
)

// EncodingCreator defines a constructor to create an EncodingCreator from
// an io.ReadWriter.
type EncodingCreator interface {
	Create(io.ReadWriter) EncodeDecoder
}

// EncodeDecoder combines Encoder and Decoder.
type EncodeDecoder interface {
	Encoder
	Decoder
}

// Encoder should encode a v into io.Writer given to EncodingCreator.
type Encoder interface {
	Encode(v interface{}) error
}

// Decoder should decode a v from io.Reader given to EncodingCreator.
type Decoder interface {
	Decode(v interface{}) error
}

// Encoding implements EncodeDecoder
type Encoding struct {
	Encoder
	Decoder
}

// NewEncoding returns a new EncodeDecoder implementation from e end d.
func NewEncoding(e Encoder, d Decoder) EncodeDecoder {
	return &Encoding{
		Encoder: e,
		Decoder: d,
	}
}

// DefaultJSONEncodingCreator implements EncodingCreator for JSON encoding.
var DefaultJSONEncodingCreator = &JSONEncodingCreator{}

// DefaultYAMLEncodingCreator implements EncodingCreator for YAML encoding.
var DefaultYAMLEncodingCreator = &YAMLEncodingCreator{}

// DefaultTOMLEncodingCreator implements EncodingCreator for TOML encoding.
var DefaultTOMLEncodingCreator = &TOMLEncodingCreator{}

// JSONEncodingCreator implements EncodingCreator for JSON encoding.
type JSONEncodingCreator struct{}

func (e *JSONEncodingCreator) Create(rw io.ReadWriter) EncodeDecoder {
	return NewEncoding(json.NewEncoder(rw), json.NewDecoder(rw))
}

// YAMLEncodingCreator implements EncodingCreator for JSON encoding.
type YAMLEncodingCreator struct{}

func (e *YAMLEncodingCreator) Create(rw io.ReadWriter) EncodeDecoder {
	return NewEncoding(yaml.NewEncoder(rw), yaml.NewDecoder(rw))
}

// TOMLEncodingCreator implements EncodingCreator for JSON encoding.
type TOMLEncodingCreator struct{}

func (e *TOMLEncodingCreator) Create(rw io.ReadWriter) EncodeDecoder {
	return NewEncoding(toml.NewEncoder(rw), toml.NewDecoder(rw))
}
