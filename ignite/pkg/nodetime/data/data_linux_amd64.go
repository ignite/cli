package data

import _ "embed" // embed is required for binary embedding.

//go:embed nodetime-linux-amd64.tar.gz
var binaryCompressed []byte
