// +build darwin,arm64

package protoc

import _ "embed" // embed is required for binary embedding.

//go:embed data/protoc-darwin-arm64
var binary []byte
