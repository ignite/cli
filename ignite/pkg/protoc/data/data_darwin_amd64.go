package data

import _ "embed" // embed is required for binary embedding.

//go:embed protoc-darwin-amd64
var binary []byte
