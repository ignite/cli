package cache

import (
	"bytes"
	"encoding/gob"
	"os"
	"path/filepath"

	"github.com/ignite-hq/cli/ignite/chainconfig"
	"github.com/ignite-hq/cli/ignite/pkg/xfilepath"
	bolt "go.etcd.io/bbolt"
)

const cacheFolderName = "cache"
const dbName = "ignite_cache.db"

var (
	// This is the top level folder where each chain's cache will be stored
	// When using NewChainStorage, each chain gets a separate folder
	cacheFolder = xfilepath.Join(chainconfig.ConfigDirPath, xfilepath.Path(cacheFolderName))

	// We need to make sure we don't open multiple connections to the same cache file, so we keep track of them here
	dbInstances = make(map[string]*bolt.DB)
)

// Storage is meant to be passed around and used by the New function (which provides namespacing and type-safety)
type Storage struct {
	db *bolt.DB
}

// Cache is a namespaced and type-safe key-value store
type Cache[T any] struct {
	storage   Storage
	namespace string
}

// NewChainStorage creates a separate cache storage for a chain with chainName
// It is safe to call multiple times
func NewChainStorage(chainName string) (Storage, error) {
	dir, err := xfilepath.Join(cacheFolder, xfilepath.Path(chainName))()
	if err != nil {
		return Storage{}, err
	}

	return NewStorage(dir)
}

// NewStorage create a local db file (if it doesn't exist) and opens a connection to it (or reuses an existing one)
// Usually NewChainStorage is more appropriate to use since it makes it possible to clear the cache for a particular project
// It is safe to call multiple times
func NewStorage(dir string) (Storage, error) {
	db, ok := dbInstances[dir]
	if ok {
		return Storage{db}, nil
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return Storage{}, err
	}

	storagePath := filepath.Join(dir, dbName)
	db, err := bolt.Open(storagePath, 0640, nil)
	if err != nil {
		return Storage{}, err
	}

	dbInstances[dir] = db
	return Storage{db}, nil
}

// New creates a namespaced and typesafe key-value Cache
func New[T any](storage Storage, namespace string) Cache[T] {
	return Cache[T]{
		storage:   storage,
		namespace: namespace,
	}
}

// Clear deletes all namespaces and cached values
func (s Storage) Clear() error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			return tx.DeleteBucket(name)
		})
	})
}

// Put sets key to value within the namespace
// If the key already exists, it will be overwritten
func (c Cache[T]) Put(key string, value T) error {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(value); err != nil {
		return err
	}
	result := buf.Bytes()

	return c.storage.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(c.namespace))
		if err != nil {
			return err
		}
		return b.Put([]byte(key), result)
	})
}

// Get fetches the value of key within the namespace.
// If no value exists, it will return found == false
func (c Cache[T]) Get(key string) (val T, found bool, err error) {
	err = c.storage.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(c.namespace))
		if b == nil {
			found = false
			return nil
		}
		c := b.Cursor()
		if k, v := c.Seek([]byte(key)); bytes.Equal(k, []byte(key)) {
			if v == nil {
				found = false
				return nil
			}

			var decodedVal T
			d := gob.NewDecoder(bytes.NewReader(v))
			if err := d.Decode(&decodedVal); err != nil {
				return err
			}

			val = decodedVal
			found = true
		} else {
			found = false
		}

		return nil
	})

	return
}

// Delete removes a value for key within the namespace
func (c Cache[T]) Delete(key string) error {
	return c.storage.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(c.namespace))
		if b == nil {
			return nil
		}

		return b.Delete([]byte(key))
	})
}
