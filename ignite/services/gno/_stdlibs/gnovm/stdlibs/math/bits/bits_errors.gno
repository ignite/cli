// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bits

import (
	"errors"
)

var overflowError = errors.New("math/bits: integer overflow")

var divideError = errors.New("math/bits: integer divide by zero")
