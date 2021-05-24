// +build darwin,amd64

package protoc

import _ "embed" // embed is required for binary embedding.

//go:embed data/protoc-darwin-amd64
var binary []byte
