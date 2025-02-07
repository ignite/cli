package xos

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStdOutRedirect(t *testing.T) {
	tests := []struct {
		name         string
		output       io.Writer
		isRedirected bool
		err          error
	}{
		{
			name:         "test stdOut",
			output:       os.Stdout,
			isRedirected: true,
		},
		{
			name:         "test stdErr",
			output:       os.Stderr,
			isRedirected: true,
		},
		{
			name:         "test stdErr",
			output:       &bytes.Buffer{},
			isRedirected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, cancel, err := StdOutRedirect(tt.output)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
			defer cancel()

			require.Equal(t, tt.isRedirected, isatty.IsTerminal(got.Fd()))
		})
	}
}
