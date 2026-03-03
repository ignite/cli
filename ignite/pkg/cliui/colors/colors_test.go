package colors

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSprintFunc(t *testing.T) {
	got := SprintFunc(Green)("ignite")
	require.Contains(t, got, "ignite")
}

func TestFormatters(t *testing.T) {
	require.Contains(t, Info("a"), "a")
	require.Contains(t, Infof("%s", "b"), "b")
	require.Contains(t, Error("c"), "c")
	require.Contains(t, Success("d"), "d")
	require.Contains(t, Modified("e"), "e")
	require.Contains(t, Name("f"), "f")
	require.Contains(t, Mnemonic("g"), "g")
	require.Contains(t, Spinner("h"), "h")
	require.Contains(t, Faint("i"), "i")
}
