// +build darwin

package protoc

import _ "embed" // embed is required for binary embedding.

//go:embed data/protoc-darwin-amd64
var binary []byte
