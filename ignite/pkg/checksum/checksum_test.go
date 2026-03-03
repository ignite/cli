package checksum

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStrings(t *testing.T) {
	h := sha256.Sum256([]byte("abc"))
	require.Equal(t, fmt.Sprintf("%x", h[:]), Strings("a", "b", "c"))
}

func TestSum(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "a.txt"), []byte("alpha"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "b.txt"), []byte("beta"), 0o600))

	out := filepath.Join(t.TempDir(), "checksums.txt")
	require.NoError(t, Sum(dir, out))

	content, err := os.ReadFile(out)
	require.NoError(t, err)
	text := string(content)
	require.Contains(t, text, "  a.txt\n")
	require.Contains(t, text, "  b.txt\n")
}

func TestBinary(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "fake-bin")
	data := []byte("#!/bin/sh\necho test\n")
	require.NoError(t, os.WriteFile(bin, data, 0o700))

	want := sha256.Sum256(data)
	got, err := Binary(bin)
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("%x", want[:]), got)
}

func TestBinaryReturnsErrorForMissingFile(t *testing.T) {
	_, err := Binary(strings.TrimSpace(filepath.Join(t.TempDir(), "missing-bin")))
	require.Error(t, err)
}
