package randstr

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunesReturnsRequestedLength(t *testing.T) {
	for _, n := range []int{0, 1, 16, 64} {
		got := Runes(n)
		require.Len(t, got, n)
	}
}

func TestRunesUsesOnlyLowercaseLetters(t *testing.T) {
	allowed := make(map[rune]struct{}, len(letterRunes))
	for _, r := range letterRunes {
		allowed[r] = struct{}{}
	}

	got := Runes(128)
	for _, r := range got {
		_, ok := allowed[r]
		require.Truef(t, ok, "unexpected rune %q in %q", r, got)
	}
}
