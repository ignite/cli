package errorview

import (
	stdErrors "errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorString(t *testing.T) {
	view := NewError(stdErrors.New("  hello world  "))
	out := view.String()

	require.Contains(t, out, "hello world")
	require.NotContains(t, out, "  hello world  ")
}
