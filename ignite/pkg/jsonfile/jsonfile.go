package jsonfile

import (
	"bufio"
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

	"github.com/buger/jsonparser"
	"github.com/pkg/errors"

	"github.com/ignite/cli/ignite/pkg/tarball"
)

const (
	keySeparator = "."
)

var (
	// ErrFieldNotFound parameter not found into json.
	ErrFieldNotFound = errors.New("JSON field not found")

	// ErrInvalidValueType invalid value type.
	ErrInvalidValueType = errors.New("invalid value type")

	// ErrInvalidURL invalid file URL.
	ErrInvalidURL = errors.New("invalid file URL")
)

type (
	// JSONFile represents the JSON file and also implements the io.write interface,
	// saving directly to the file.
	JSONFile struct {
		file        ReadWriteSeeker
		tarballPath string
		url         string
		updates     map[string][]byte
		cache       []byte
	}

	// UpdateFileOption configures file update function with key and value.
	UpdateFileOption func(map[string][]byte)
)

type (
	// writeTruncate represents the truncate method from io.WriteSeeker interface.
	writeTruncate interface {
		Truncate(size int64) error
	}

	// ReadWriteSeeker represents the owns ReadWriteSeeker interface inherit from io.ReadWriteSeeker.
	ReadWriteSeeker interface {
		io.ReadWriteSeeker
		Close() error
		Sync() error
	}
)

// New creates a new JSONFile.
func New(file ReadWriteSeeker) *JSONFile {
	return &JSONFile{
		updates: make(map[string][]byte),
		file:    file,
	}
}

// FromPath parses a JSONFile object from path.
func FromPath(path string) (*JSONFile, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0o600)
	if err != nil {
		return nil, errors.Wrap(err, "cannot open the file")
	}
	return New(file), nil
}

// FromURL fetches the file from the given URL and returns its content.
// If tarballFileName is not empty, the URL is interpreted as a tarball file,
// tarballFileName is extracted from it and is returned instead of the URL
// content.
func FromURL(ctx context.Context, url, destPath, tarballFileName string) (*JSONFile, error) {
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

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrInvalidURL
	}

	// Remove the old file if exists and create a new one
	if err := os.RemoveAll(destPath); err != nil {
		return nil, err
	}
	file, err := os.OpenFile(destPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0o600)
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
	if err != nil && !errors.Is(err, tarball.ErrNotGzipType) && !errors.Is(err, tarball.ErrInvalidFileName) {
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

// Bytes returns the jsonfile byte array.
func (f *JSONFile) Bytes() ([]byte, error) {
	file := f.cache
	if file != nil {
		return file, nil
	}
	if err := f.Reset(); err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(f.file)
	for scanner.Scan() {
		file = append(file, scanner.Bytes()...)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	f.cache = file
	return file, nil
}

// Field returns the param by key and the position into byte slice from the file reader.
// Key can be a path to a nested parameters eg: app_state.staking.accounts.
func (f *JSONFile) Field(key string, param interface{}) error {
	file, err := f.Bytes()
	if err != nil {
		return err
	}

	value, dataType, _, err := jsonparser.Get(file, strings.Split(key, keySeparator)...)
	if errors.Is(err, jsonparser.KeyPathNotFoundError) {
		return ErrFieldNotFound
	} else if err != nil {
		return err
	}

	switch dataType {
	case jsonparser.Boolean, jsonparser.Array, jsonparser.Number, jsonparser.Object:
		err := json.Unmarshal(value, param)
		if _, ok := err.(*json.UnmarshalTypeError); ok { //nolint:errorlint
			return ErrInvalidValueType
		} else if err != nil {
			return err
		}
	case jsonparser.String:
		paramStr, ok := param.(*string)
		if !ok {
			return ErrInvalidValueType
		}
		*paramStr, err = jsonparser.ParseString(value)
		if err != nil {
			return err
		}
	case jsonparser.NotExist:
	case jsonparser.Null:
	case jsonparser.Unknown:
	default:
		return ErrInvalidValueType
	}
	return nil
}

// WithKeyValue updates a file value object by key.
func WithKeyValue(key string, value string) UpdateFileOption {
	return func(update map[string][]byte) {
		update[key] = []byte(`"` + value + `"`)
	}
}

// WithKeyValueByte updates a file byte value object by key.
func WithKeyValueByte(key string, value []byte) UpdateFileOption {
	return func(update map[string][]byte) {
		update[key] = value
	}
}

// WithKeyValueTimestamp updates a time value.
func WithKeyValueTimestamp(key string, t int64) UpdateFileOption {
	return func(update map[string][]byte) {
		formatted := time.Unix(t, 0).UTC().Format(time.RFC3339Nano)
		update[key] = []byte(`"` + formatted + `"`)
	}
}

// WithKeyValueInt updates a file int value object by key.
func WithKeyValueInt(key string, value int64) UpdateFileOption {
	return func(update map[string][]byte) {
		update[key] = []byte(strconv.FormatInt(value, 10))
	}
}

// WithKeyValueUint updates a file uint value object by key.
func WithKeyValueUint(key string, value uint64) UpdateFileOption {
	return WithKeyValueInt(key, int64(value))
}

// Update updates the file with the new parameters by key.
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

// Write implement the write method for io.Writer interface.
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
	f.cache = p

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

	return length, nil
}

// truncate removes the current file content.
func truncate(rws io.WriteSeeker, size int) error {
	t, ok := rws.(writeTruncate)
	if !ok {
		return errors.New("truncate: unable to truncate")
	}

	return t.Truncate(int64(size))
}

// Close the file.
func (f *JSONFile) Close() error {
	return f.file.Close()
}

// URL returns the genesis URL.
func (f *JSONFile) URL() string {
	return f.url
}

// TarballPath returns the tarball path.
func (f *JSONFile) TarballPath() string {
	return f.tarballPath
}

// Hash returns the hash of the file.
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

// String returns the file string.
func (f *JSONFile) String() (string, error) {
	if err := f.Reset(); err != nil {
		return "", err
	}

	data, err := io.ReadAll(f.file)
	return string(data), err
}

// Reset sets the offset for the next Read or Write to 0.
func (f *JSONFile) Reset() error {
	// TODO find a better way to reset or create a
	// read of copy the writer with io.TeeReader
	_, err := f.file.Seek(0, 0)
	return err
}
