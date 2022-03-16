package checksum

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
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

		if _, err := b.WriteString(fmt.Sprintf("%x %s\n", h.Sum(nil), info.Name())); err != nil {
			return err
		}
	}

	return os.WriteFile(outPath, b.Bytes(), 0666)
}

func BinaryChecksum(binaryName string) (string, error) {
	// get binary path
	binaryPath, err := exec.LookPath(binaryName)
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

func SHA256Checksum(inputs ...[]byte) string {
	h := sha256.New()
	for _, input := range inputs {
		h.Write(input)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
