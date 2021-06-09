package checksum

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Sum reads files from dirPath, calculates sha256 for each file and creates a new checksum
// file for them at outPath.
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
