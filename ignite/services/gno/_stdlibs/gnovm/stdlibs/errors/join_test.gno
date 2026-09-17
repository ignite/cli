// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors_test

import (
	"errors"
	"testing"
)

func TestJoinReturnsNil(t *testing.T) {
	if err := errors.Join(); err != nil {
		t.Errorf("errors.Join() = %v, want nil", err)
	}
	if err := errors.Join(nil); err != nil {
		t.Errorf("errors.Join(nil) = %v, want nil", err)
	}
	if err := errors.Join(nil, nil); err != nil {
		t.Errorf("errors.Join(nil, nil) = %v, want nil", err)
	}
}

func TestJoin(t *testing.T) {
	err1 := errors.New("err1")
	err2 := errors.New("err2")

	for _, test := range []struct {
		errs []error
		want []error
	}{{
		errs: []error{err1},
		want: []error{err1},
	}, {
		errs: []error{err1, err2},
		want: []error{err1, err2},
	}, {
		errs: []error{err1, nil, err2},
		want: []error{err1, err2},
	}} {
		got := errors.Join(test.errs...).(interface{ Unwrap() []error }).Unwrap()
		if len(got) != len(test.want) {
			t.Errorf("Join(%v) len = %d, want %d", test.errs, len(got), len(test.want))
			continue
		}
		for i, err := range got {
			if err != test.want[i] {
				t.Errorf("Join(%v)[%d] = %v, want %v", test.errs, i, err, test.want[i])
			}
		}
	}
}

func TestJoinErrorMethod(t *testing.T) {
	err1 := errors.New("err1")
	err2 := errors.New("err2")

	for _, test := range []struct {
		errs []error
		want string
	}{{
		errs: []error{err1},
		want: "err1",
	}, {
		errs: []error{err1, err2},
		want: "err1\nerr2",
	}, {
		errs: []error{err1, nil, err2},
		want: "err1\nerr2",
	}} {
		got := errors.Join(test.errs...).Error()
		if got != test.want {
			t.Errorf("Join(%v).Error() = %q, want %q", test.errs, got, test.want)
		}
	}
}

func TestJoinIs(t *testing.T) {
	err1 := errors.New("err1")
	err2 := errors.New("err2")
	err3 := errors.New("err3")

	joined := errors.Join(err1, err2)

	if !errors.Is(joined, err1) {
		t.Errorf("errors.Is(Join(err1, err2), err1) = false, want true")
	}
	if !errors.Is(joined, err2) {
		t.Errorf("errors.Is(Join(err1, err2), err2) = false, want true")
	}
	if errors.Is(joined, err3) {
		t.Errorf("errors.Is(Join(err1, err2), err3) = true, want false")
	}
}
