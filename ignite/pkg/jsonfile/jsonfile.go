package jsonfile

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ignite-hq/cli/ignite/pkg/tarball"

	"github.com/buger/jsonparser"
	"github.com/pkg/errors"
)

const (
	keySeparator = "."
)

type (
	// JSONFile represents the file
	JSONFile struct {
		file        readWriteSeeker
		tarballPath string
		url         string
		updates     map[string][]byte
	}
	// UpdateFileOption configures file update.
	UpdateFileOption func(map[string][]byte)

	writeTruncate interface {
		Truncate(size int64) error
	}
	readWriteSeeker interface {
		io.ReadWriteSeeker
		Close() error
		Sync() error
	}
)

var (
	// ErrFieldNotFound parameter not found into json
	ErrFieldNotFound = errors.New("JSON field not found")
	// ErrInvalidValueType invalid value type
	ErrInvalidValueType = errors.New("invalid value type")
)

// New creates a new JSONFile
func New(file readWriteSeeker) *JSONFile {
	return &JSONFile{
		updates: make(map[string][]byte),
		file:    file,
	}
}

// FromPath parse JSONFile object from path
func FromPath(path string) (*JSONFile, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0600)

	if err != nil {
		return nil, errors.Wrap(err, "cannot open the file")
	}
	return New(file), nil
}

// FromURL fetches the file from the given URL and returns its content.
func FromURL(ctx context.Context, url, path, tarballFileName string) (*JSONFile, error) {
	// TODO create a cache system to avoid download genesis with the same hash again

	// Download the file from URL
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Remove the old file if exists and create a new one
	if err := os.RemoveAll(path); err != nil {
		return nil, err
	}
	file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create the file")
	}

	// Copy the downloaded file to buffer and the opened file
	var buf bytes.Buffer
	if _, err := io.Copy(file, io.TeeReader(resp.Body, &buf)); err != nil {
		return nil, err
	}

	// Check if the downloaded file is a tarball and extract only the necessary JSON file
	var ext bytes.Buffer
	tarballPath, err := tarball.ExtractFile(&buf, &ext, tarballFileName)
	if err != nil && err != tarball.ErrNotGzipType {
		return nil, err
	} else if err == nil {
		// Erase the tarball bite code from the file and copy the correct one
		if err := truncate(file, 0); err != nil {
			return nil, err
		}
		if _, err := io.Copy(file, &ext); err != nil {
			return nil, err
		}
	}

	return &JSONFile{
		updates:     make(map[string][]byte),
		file:        file,
		url:         url,
		tarballPath: tarballPath,
	}, nil
}

// Field return the param by key and the position into byte slice from the file reader.
// Key can be a path to a nested parameter eg: app_state.staking.accounts
func (f *JSONFile) Field(key string, param interface{}) error {
	if err := f.Reset(); err != nil {
		return err
	}
	dec := json.NewDecoder(f.file)
	// Split the keys by the separator to find nested JSON parameters eg: app_state.staking.accounts
	keys := strings.Split(key, keySeparator)
	// Instead of unmarshal the whole content of a file
	// this will decode one line/record at a time
	for {
		t, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		name, ok := t.(string)
		if !ok {
			continue
		}
		// If the nested path was found
		if name == keys[0] {
			if len(keys) > 1 {
				keys = keys[1:]
				continue
			}

			// Try to decode the all data first
			err := dec.Decode(&param)
			if err == nil {
				return nil
			}

			// If not decode, only get the JSON value from the key
			t, err := dec.Token()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			switch t := t.(type) {
			case int:
				param = strconv.Itoa(t)
			case string:
				param = t
			default:
				return ErrInvalidValueType
			}
			return nil
		}
	}
	return errors.Wrap(ErrFieldNotFound, key)
}

// WithKeyValue update a file value object by key
func WithKeyValue(key string, value string) UpdateFileOption {
	return func(update map[string][]byte) {
		update[key] = []byte(`"` + value + `"`)
	}
}

// WithTime update a time value
func WithTime(key string, t int64) UpdateFileOption {
	return func(update map[string][]byte) {
		formatted := time.Unix(t, 0).UTC().Format(time.RFC3339Nano)
		update[key] = []byte(`"` + formatted + `"`)
	}
}

// WithKeyIntValue update a file int value object by key
func WithKeyIntValue(key string, value int) UpdateFileOption {
	return func(update map[string][]byte) {
		update[key] = []byte(strconv.Itoa(value))
	}
}

// Update updates the file file with the new parameters by key
func (f *JSONFile) Update(opts ...UpdateFileOption) error {
	for _, opt := range opts {
		opt(f.updates)
	}
	if err := f.Reset(); err != nil {
		return err
	}
	_, err := io.Copy(f, f.file)
	return err
}

// Write implement the write method for io.Writer interface
func (f *JSONFile) Write(p []byte) (int, error) {
	var err error
	length := len(p)
	for key, value := range f.updates {
		p, err = jsonparser.Set(p, value, strings.Split(key, keySeparator)...)
		if err != nil {
			return 0, err
		}
		delete(f.updates, key)
	}

	err = truncate(f.file, 0)
	if err != nil {
		return 0, err
	}

	if err := f.Reset(); err != nil {
		return 0, err
	}
	n, err := f.file.Write(p)
	if err != nil {
		return n, err
	}

	if n != len(p) {
		return n, io.ErrShortWrite
	}

	// FIXME passing the new byte slice length throws an error
	// because the reader has less byte length than the writer
	// https://cs.opensource.google/go/go/+/refs/tags/go1.18:src/io/io.go;l=432
	return length, nil
}

// truncate remove the current file content
func truncate(rws io.WriteSeeker, size int) error {
	t, ok := rws.(writeTruncate)
	if !ok {
		return errors.New("truncate: unable to truncate")
	}
	return t.Truncate(int64(size))
}

// Close the file
func (f *JSONFile) Close() error {
	return f.file.Close()
}

// URL returns the genesis URL
func (f *JSONFile) URL() string {
	return f.url
}

// TarballPath returns the tarball path
func (f *JSONFile) TarballPath() string {
	return f.tarballPath
}

// Hash returns the hash of the file
func (f *JSONFile) Hash() (string, error) {
	if err := f.Reset(); err != nil {
		return "", err
	}
	h := sha256.New()
	if _, err := io.Copy(h, f.file); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// String returns the file string
func (f *JSONFile) String() (string, error) {
	if err := f.Reset(); err != nil {
		return "", err
	}
	data, err := io.ReadAll(f.file)
	return string(data), err
}

// Reset sets the offset for the next Read or Write to 0
func (f *JSONFile) Reset() error {
	// TODO find a better way to reset or create a
	// read of copy the writer with io.TeeReader
	_, err := f.file.Seek(0, 0)
	return err
}
