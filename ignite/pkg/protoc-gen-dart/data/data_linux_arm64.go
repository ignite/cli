package data

import _ "embed" // embed is required for binary embedding.

//go:embed protoc-gen-dart_linux_arm64
var binary []byte
