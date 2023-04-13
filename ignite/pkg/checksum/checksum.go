package checksum

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ignite/cli/ignite/pkg/xexec"
)

// Sum reads files from dirPath, calculates sha256 for each file and creates a new checksum
// file for them in outPath.
func Sum(dirPath, outPath string) error {
	var b bytes.Buffer

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, info := range files {
		path := filepath.Join(dirPath, info.Name())
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		h := sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			return err
		}

		// Note that checksum entry has two spaces as separator to follow
		// FIPS-180-2 regarding the character prefix for text file types.
		// This is required for tools like sha256sum with a strict verification.
		if _, err := b.WriteString(fmt.Sprintf("%x  %s\n", h.Sum(nil), info.Name())); err != nil {
			return err
		}
	}

	return os.WriteFile(outPath, b.Bytes(), 0o666)
}

// Binary returns SHA256 hash of executable file, file is searched by name in PATH.
func Binary(binaryName string) (string, error) {
	// get binary path
	binaryPath, err := xexec.ResolveAbsPath(binaryName)
	if err != nil {
		return "", err
	}
	f, err := os.Open(binaryPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// Strings concatenates all inputs and returns SHA256 hash of them.
func Strings(inputs ...string) string {
	h := sha256.New()
	for _, input := range inputs {
		h.Write([]byte(input))
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
