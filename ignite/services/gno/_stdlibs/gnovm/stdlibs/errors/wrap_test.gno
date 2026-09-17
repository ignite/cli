// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors_test

import (
	"errors"
	"testing"
)

func TestIs(t *testing.T) {
	err1 := errors.New("1")
	erra := wrapped{"wrap 2", err1}
	errb := wrapped{"wrap 3", erra}

	err3 := errors.New("3")

	poser := &poser{"either 1 or 3", func(target error) bool {
		return target == err1 || target == err3
	}}

	testCases := []struct {
		err    error
		target error
		match  bool
	}{
		{nil, nil, true},
		{err1, nil, false},
		{err1, err1, true},
		{erra, err1, true},
		{errb, err1, true},
		{err1, err3, false},
		{erra, err3, false},
		{errb, err3, false},
		{poser, err1, true},
		{poser, err3, true},
		{poser, erra, false},
		{poser, errb, false},
		{errorUncomparable{}, errorUncomparable{}, true},
		{errorUncomparable{}, &errorUncomparable{}, false},
		{&errorUncomparable{}, errorUncomparable{}, true},
		{&errorUncomparable{}, &errorUncomparable{}, false},
		{errorUncomparable{}, err1, false},
		{errorUncomparable{}, nil, false},
		{nil, errorUncomparable{}, false},
		{multiErr{}, err1, false},
		{multiErr{err1, err3}, err1, true},
		{multiErr{err3, err1}, err1, true},
		{multiErr{err1, err3}, errors.New("x"), false},
		{multiErr{err3, errb}, errb, true},
		{multiErr{err3, errb}, erra, true},
		{multiErr{err3, errb}, err1, true},
		{multiErr{errb, err3}, err1, true},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			if got := errors.Is(tc.err, tc.target); got != tc.match {
				t.Errorf("Is(%v, %v) = %v, want %v", tc.err, tc.target, got, tc.match)
			}
		})
	}
}

// TestIsUncomparableNoMatch locks in the #5713 dependency: an uncomparable
// target whose Is returns false must not match. This fails until #5713 makes
// the VM panic on comparing uncomparable types, since comparable() currently
// returns true and the == branch wrongly matches on the nil slice fields.
func TestIsUncomparableNoMatch(t *testing.T) {
	a := errorUncomparableNoMatch{}
	b := errorUncomparableNoMatch{}
	if got := errors.Is(a, b); got != false {
		t.Errorf("Is(%v, %v) = %v, want false", a, b, got)
	}
}

func TestUnwrap(t *testing.T) {
	err1 := errors.New("1")
	erra := wrapped{"wrap 2", err1}

	testCases := []struct {
		err  error
		want error
	}{
		{nil, nil},
		{wrapped{"wrapped", nil}, nil},
		{err1, nil},
		{erra, err1},
		{wrapped{"wrap 3", erra}, erra},
	}
	for _, tc := range testCases {
		if got := errors.Unwrap(tc.err); got != tc.want {
			t.Errorf("Unwrap(%v) = %v, want %v", tc.err, got, tc.want)
		}
	}
}

type poser struct {
	msg string
	f   func(error) bool
}

func (p *poser) Error() string { return p.msg }
func (p *poser) Is(err error) bool {
	return p.f(err)
}

type wrapped struct {
	msg string
	err error
}

func (e wrapped) Error() string { return e.msg }
func (e wrapped) Unwrap() error { return e.err }

type errorUncomparable struct {
	f []string
}

func (errorUncomparable) Error() string {
	return "uncomparable error"
}

func (errorUncomparable) Is(target error) bool {
	_, ok := target.(errorUncomparable)
	if !ok {
		return false
	}
	return true
}

type errorUncomparableNoMatch struct {
	f []string
}

func (errorUncomparableNoMatch) Error() string { return "uncomparable no match" }

func (errorUncomparableNoMatch) Is(target error) bool { return false }

type multiErr []error

func (m multiErr) Error() string   { return "multiErr" }
func (m multiErr) Unwrap() []error { return []error(m) }
