package tarball

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"path/filepath"
	"strings"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

var (
	// ErrGzipFileNotFound the file not found in the gzip.
	ErrGzipFileNotFound = errors.New("file not found in the gzip")
	// ErrNotGzipType the file is not a gzip.
	ErrNotGzipType = errors.New("file is not a gzip type")
	// ErrInvalidFileName the file name is invalid.
	ErrInvalidFileName = errors.New("invalid file name")
	// ErrInvalidFilePath the file path is invalid.
	ErrInvalidFilePath = errors.New("invalid file path")
	// ErrFileTooLarge the file is too large to extract.
	ErrFileTooLarge = errors.New("file too large to extract")
)

// ExtractFile founds and reads a specific file into a gzip file and folders recursively.
func ExtractFile(reader io.Reader, out io.Writer, fileName string) (string, error) {
	if fileName == "" {
		return "", ErrInvalidFileName
	}
	archive, err := gzip.NewReader(reader)
	// Verify if is a GZIP file
	if errors.Is(err, io.EOF) || errors.Is(err, gzip.ErrHeader) {
		return "", ErrNotGzipType
	} else if err != nil {
		return "", err
	}
	defer archive.Close()

	tarReader := tar.NewReader(archive)
	// Read the tarball files and find only the necessary file
	for {
		header, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			return "", ErrGzipFileNotFound
		} else if err != nil {
			return header.Name, err
		}

		// Validate the file path
		if !isValidPath(header.Name) {
			return "", ErrInvalidFilePath
		}

		switch header.Typeflag {
		case tar.TypeDir:
			continue
		case tar.TypeReg:
			name := filepath.Base(header.Name)
			if fileName == name {
				// Limit the size of the file to extract
				if header.Size > 100<<20 { // 100 MB limit
					return "", ErrFileTooLarge
				}
				limitedReader := io.LimitReader(tarReader, 1000<<20) // 1000 MB limit
				_, err := io.Copy(out, limitedReader)
				return header.Name, err
			}
		default:
			continue
		}
	}
}

// isValidPath checks for directory traversal attacks.
func isValidPath(filePath string) bool {
	cleanPath := filepath.Clean(filePath)
	return !strings.Contains(cleanPath, "..")
}
