package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBufFiles(t *testing.T) {
	want := []string{"buf.work.yaml"}
	protoDir, err := os.ReadDir("files/{{protoDir}}")
	require.NoError(t, err)
	for _, e := range protoDir {
		want = append(want, filepath.Join("{{protoDir}}", e.Name()))
	}

	got, err := BufFiles()
	require.NoError(t, err)
	require.ElementsMatch(t, want, got)
}

func TestCutTemplatePrefix(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
		ok   bool
	}{
		{
			name: "with prefix",
			arg:  "{{protoDir}}/myvalue",
			want: "myvalue",
			ok:   true,
		},
		{
			name: "with 2 prefix",
			arg:  "{{protoDir}}/{{protoDir}}/myvalue",
			want: "{{protoDir}}/myvalue",
			ok:   true,
		},
		{
			name: "without prefix",
			arg:  "myvalue",
			want: "myvalue",
			ok:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := CutTemplatePrefix(tt.arg)
			require.Equal(t, tt.ok, ok)
			require.Equal(t, tt.want, got)
		})
	}
}
