// +build darwin,amd64

package nodetime

import _ "embed" // embed is required for binary embedding.

//go:embed nodetime-darwin-x64.tar.gz
var binaryCompressed []byte
