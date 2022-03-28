package xgzip

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"os"
	"path/filepath"
)

var (
	// ErrGzipFileNotFound the file not found in the gzip
	ErrGzipFileNotFound = errors.New("file not found in the gzip")
	// ErrFileNotFound the file gzip not found into the folder
	ErrFileNotFound = errors.New("gzip file not found")
	// ErrInvalidGzipFile the gzip file is invalid
	ErrInvalidGzipFile = errors.New("invalid gzip file")
)

// ReadFile founds and reads a specific file into a gzip file and folders recursively
func ReadFile(source, file string) ([]byte, error) {
	reader, err := os.Open(source)
	if os.IsNotExist(err) {
		return nil, ErrFileNotFound
	} else if err != nil {
		return nil, err
	}
	defer reader.Close()

	archive, err := gzip.NewReader(reader)
	if err == io.EOF {
		return nil, ErrInvalidGzipFile
	} else if err != nil {
		return nil, err
	}
	defer archive.Close()

	tarReader := tar.NewReader(archive)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			return nil, ErrGzipFileNotFound
		} else if err != nil {
			return nil, err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			continue
		case tar.TypeReg:
			name := filepath.Base(header.Name)
			if file == name {
				var bout bytes.Buffer
				if _, err := io.Copy(&bout, tarReader); err != nil {
					return nil, err
				}
				return bout.Bytes(), nil
			}
		default:
			continue
		}
	}
}
