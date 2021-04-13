// +build linux

package protoc

import _ "embed" // embed is required for binary embedding.

//go:embed data/protoc-linux-amd64
var binary []byte
