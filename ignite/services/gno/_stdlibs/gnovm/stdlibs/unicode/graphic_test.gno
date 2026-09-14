// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unicode_test

import (
	"testing"
	uu "unicode"
)

// Independently check that the special "Is" functions work
// in the Latin-1 range through the property table.

func TestIsControlLatin1(t *testing.T) {
	for i := rune(0); i <= uu.MaxLatin1; i++ {
		got := uu.IsControl(i)
		want := false
		switch {
		case 0x00 <= i && i <= 0x1F:
			want = true
		case 0x7F <= i && i <= 0x9F:
			want = true
		}
		if got != want {
			t.Errorf("%U incorrect: got %t; want %t", i, got, want)
		}
	}
}

func TestIsLetterLatin1(t *testing.T) {
	for i := rune(0); i <= uu.MaxLatin1; i++ {
		got := uu.IsLetter(i)
		want := uu.Is(uu.Letter, i)
		if got != want {
			t.Errorf("%U incorrect: got %t; want %t", i, got, want)
		}
	}
}

func TestIsUpperLatin1(t *testing.T) {
	for i := rune(0); i <= uu.MaxLatin1; i++ {
		got := uu.IsUpper(i)
		want := uu.Is(uu.Upper, i)
		if got != want {
			t.Errorf("%U incorrect: got %t; want %t", i, got, want)
		}
	}
}

func TestIsLowerLatin1(t *testing.T) {
	for i := rune(0); i <= uu.MaxLatin1; i++ {
		got := uu.IsLower(i)
		want := uu.Is(uu.Lower, i)
		if got != want {
			t.Errorf("%U incorrect: got %t; want %t", i, got, want)
		}
	}
}

func TestNumberLatin1(t *testing.T) {
	for i := rune(0); i <= uu.MaxLatin1; i++ {
		got := uu.IsNumber(i)
		want := uu.Is(uu.Number, i)
		if got != want {
			t.Errorf("%U incorrect: got %t; want %t", i, got, want)
		}
	}
}

func TestIsPrintLatin1(t *testing.T) {
	for i := rune(0); i <= uu.MaxLatin1; i++ {
		got := uu.IsPrint(i)
		want := uu.In(i, uu.PrintRanges...)
		if i == ' ' {
			want = true
		}
		if got != want {
			t.Errorf("%U incorrect: got %t; want %t", i, got, want)
		}
	}
}

func TestIsGraphicLatin1(t *testing.T) {
	for i := rune(0); i <= uu.MaxLatin1; i++ {
		got := uu.IsGraphic(i)
		want := uu.In(i, uu.GraphicRanges...)
		if got != want {
			t.Errorf("%U incorrect: got %t; want %t", i, got, want)
		}
	}
}

func TestIsPunctLatin1(t *testing.T) {
	for i := rune(0); i <= uu.MaxLatin1; i++ {
		got := uu.IsPunct(i)
		want := uu.Is(uu.Punct, i)
		if got != want {
			t.Errorf("%U incorrect: got %t; want %t", i, got, want)
		}
	}
}

func TestIsSpaceLatin1(t *testing.T) {
	for i := rune(0); i <= uu.MaxLatin1; i++ {
		got := uu.IsSpace(i)
		want := uu.Is(uu.White_Space, i)
		if got != want {
			t.Errorf("%U incorrect: got %t; want %t", i, got, want)
		}
	}
}

func TestIsSymbolLatin1(t *testing.T) {
	for i := rune(0); i <= uu.MaxLatin1; i++ {
		got := uu.IsSymbol(i)
		want := uu.Is(uu.Symbol, i)
		if got != want {
			t.Errorf("%U incorrect: got %t; want %t", i, got, want)
		}
	}
}
