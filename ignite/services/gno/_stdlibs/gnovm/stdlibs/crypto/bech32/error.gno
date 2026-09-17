// Copyright (c) 2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package bech32

// ErrMixedCase is returned when the bech32 string has both lower and uppercase
// characters.
type ErrMixedCase struct{}

func (e ErrMixedCase) Error() string {
	return "string not all lowercase or all uppercase"
}

// ErrInvalidBitGroups is returned when conversion is attempted between byte
// slices using bit-per-element of unsupported value.
type ErrInvalidBitGroups struct{}

func (e ErrInvalidBitGroups) Error() string {
	return "only bit groups between 1 and 8 allowed"
}

// ErrInvalidIncompleteGroup is returned when then byte slice used as input has
// data of wrong length.
type ErrInvalidIncompleteGroup struct{}

func (e ErrInvalidIncompleteGroup) Error() string {
	return "invalid incomplete group"
}

// ErrInvalidLength is returned when the bech32 string has an invalid length
// given the BIP-173 defined restrictions.
type ErrInvalidLength int

func (e ErrInvalidLength) Error() string {
	return "invalid bech32 string length"
}

// ErrInvalidCharacter is returned when the bech32 string has a character
// outside the range of the supported charset.
type ErrInvalidCharacter rune

func (e ErrInvalidCharacter) Error() string {
	return "invalid character in string: " + string(e)
}

// ErrInvalidSeparatorIndex is returned when the separator character '1' is
// in an invalid position in the bech32 string.
type ErrInvalidSeparatorIndex int

func (e ErrInvalidSeparatorIndex) Error() string {
	return "invalid separator index" + string(e)
}

// ErrNonCharsetChar is returned when a character outside of the specific
// bech32 charset is used in the string.
type ErrNonCharsetChar rune

func (e ErrNonCharsetChar) Error() string {
	return "invalid character not part of charset"
}

// ErrInvalidChecksum is returned when the extracted checksum of the string
// is different than what was expected. Both the original version, as well as
// the new bech32m checksum may be specified.
type ErrInvalidChecksum struct {
	Expected  string
	ExpectedM string
	Actual    string
}

func (e ErrInvalidChecksum) Error() string {
	return "invalid checksum (expected (bech32=" + e.Expected +
		" bech32m=)" + e.ExpectedM + ", got " + e.Actual + ")"
}

// ErrInvalidDataByte is returned when a byte outside the range required for
// conversion into a string was found.
type ErrInvalidDataByte byte

func (e ErrInvalidDataByte) Error() string {
	return "invalid data byte: " + string(e)
}
