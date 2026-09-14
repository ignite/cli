// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

// Export internal functions for testing.
var (
	ExpGo      = exp
	Exp2Go     = exp2
	HypotGo    = hypot
	SqrtGo     = sqrt
	TrigReduce = trigReduce
)

const ReduceThreshold = reduceThreshold
