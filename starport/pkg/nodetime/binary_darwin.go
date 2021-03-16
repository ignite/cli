// +build darwin

package nodetime

import _ "embed" // embed is required for binary embedding.

//go:embed nodetime-darwin.tar.gz
var binaryCompressed []byte
