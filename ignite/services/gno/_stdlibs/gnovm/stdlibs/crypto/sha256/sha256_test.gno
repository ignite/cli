package sha256_test

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestSha256Sum(t *testing.T) {
	result := sha256.Sum256([]byte("sha256 this string"))
	got := result[:]
	expected := "1af1dfa857bf1d8814fe1af8983c18080019922e557f15a8a0d3db739d77aacb"

	if hex.EncodeToString(got) != expected {
		t.Errorf("got %v(%T), expected %v(%T)", hex.EncodeToString(got), got, expected, expected)
	}
}
