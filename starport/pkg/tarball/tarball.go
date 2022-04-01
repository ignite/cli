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
)

// ExtractFile founds and reads a specific file into a gzip file and folders recursively
func ExtractFile(reader io.Reader, fileName string) (io.Reader, string, error) {
	archive, err := gzip.NewReader(reader)
	if err == io.EOF || err == gzip.ErrHeader {
		return reader, "", nil
	} else if err != nil {
		return nil, "", err
	}
	defer archive.Close()

	tarReader := tar.NewReader(archive)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			return nil, "", ErrGzipFileNotFound
		} else if err != nil {
			return nil, header.Name, err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			continue
		case tar.TypeReg:
			name := filepath.Base(header.Name)
			if fileName == name {
				genesis, err := io.ReadAll(tarReader)
				if err != nil {
					return nil, "", err
				}
				return bytes.NewReader(genesis), header.Name, err
			}
		default:
			continue
		}
	}
}
