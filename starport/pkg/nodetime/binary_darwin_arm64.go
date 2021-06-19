// +build darwin arm64

package nodetime

import _ "embed" // embed is required for binary embedding.

//go:embed nodetime-darwin-arm64.tar.gz
var binaryCompressed []byte
