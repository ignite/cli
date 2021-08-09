package data

import _ "embed" // embed is required for binary embedding.

//go:embed protoc-gen-dart_darwin_amd64
var binary []byte
