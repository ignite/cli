package tarball

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"path/filepath"
)

var (
	// ErrGzipFileNotFound the file not found in the gzip
	ErrGzipFileNotFound = errors.New("file not found in the gzip")
	// ErrInvalidGzipFile the gzip file is invalid
	ErrInvalidGzipFile = errors.New("invalid gzip file")
)

// IsTarball checks if the file is a tarball
func IsTarball(tarball []byte) error {
	r, err := gzip.NewReader(bytes.NewReader(tarball))
	if err == io.EOF || err == gzip.ErrHeader {
		return ErrInvalidGzipFile
	} else if err != nil {
		return err
	}
	return r.Close()
}

// ReadFile founds and reads a specific file into a gzip file and folders recursively
func ReadFile(tarball []byte, file string) ([]byte, error) {
	archive, err := gzip.NewReader(bytes.NewReader(tarball))
	if err == io.EOF || err == gzip.ErrHeader {
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
