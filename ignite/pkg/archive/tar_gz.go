package archive

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

// CreateArchive creates a tar.gz archive from a list of files.
func CreateArchive(dir string, buf io.Writer) error {
	// Create new Writers for gzip and tar
	// These writers are chained. Writing to the tar writer will
	// write to the gzip writer which in turn will write to
	// the "buf" writer
	gw := gzip.NewWriter(buf)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	return filepath.WalkDir(dir, func(path string, _ os.DirEntry, _ error) error {
		return addToArchive(tw, path)
	})
}

func addToArchive(tw *tar.Writer, filename string) error {
	// Open the file which will be written into the archive
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get FileInfo about our file providing file size, mode, etc.
	info, err := file.Stat()
	if err != nil {
		return err
	}

	// Create a tar Header from the FileInfo data
	if info.IsDir() {
		hdr, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}
		hdr.Name = filename
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		return nil
	}

	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}

	// Use full path as name (FileInfoHeader only takes the basename)
	// If we don't do this the directory structure would
	// not be preserved
	// https://golang.org/src/archive/tar/common.go?#L626
	header.Name = filename

	// Write file header to the tar archive
	err = tw.WriteHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(tw, file)
	if err != nil {
		return err
	}

	return nil
}

// ExtractArchive extracts a tar.gz archive to the specified directory.
func ExtractArchive(outDir string, gzipStream io.Reader) error {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		targetPath := filepath.Join(outDir, header.Name) //nolint:gosec // We trust the tar file

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, 0o755); err != nil {
				return err
			}
		case tar.TypeReg:
			outFile, err := os.Create(targetPath)
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil { //nolint:gosec // We trust the tar file
				return err
			}
			outFile.Close()

		default:
			return errors.Errorf("unknown type: %s in %s", string(header.Typeflag), header.Name)
		}
	}

	return nil
}
