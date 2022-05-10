package cache

import (
	"bytes"
	"encoding/gob"
	"errors"
	"os"
	"path/filepath"
	"sync"

	bolt "go.etcd.io/bbolt"
)

const cacheFolderName = "cache"
const dbName = "ignite_cache.db"

var (
	ErrorNotFound = errors.New("no value was found with the provided key")

	// We need to make sure we don't open multiple connections to the same cache file, so we keep track of them here
	dbInstances = make(map[string]*bolt.DB)
	mu          = sync.Mutex{}
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

// NewNamespacedStorage creates a separate cache storage for a top-level namespace
// It is safe to call multiple times
func NewNamespacedStorage(rootCacheDir string, topLevelNamespace string) (Storage, error) {
	return NewStorage(filepath.Join(rootCacheDir, cacheFolderName, topLevelNamespace))
}

// NewStorage create a local db file (if it doesn't exist) and opens a connection to it (or reuses an existing one)
// Usually NewNamespacedStorage is more appropriate to use since it makes it possible to clear the cache for separate top-level namespaces
// It is safe to call multiple times
func NewStorage(dir string) (Storage, error) {
	mu.Lock()
	defer mu.Unlock()
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

// Close closes the underlying database and cleans up memory
// Attempting to use the same Storage instance or any pre-existing Cache instances
// will result in errors.
func (s Storage) Close() error {
	mu.Lock()
	defer mu.Unlock()

	dbKey := filepath.Dir(s.db.Path())
	delete(dbInstances, dbKey)

	return s.db.Close()
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
func (c Cache[T]) Get(key string) (val T, err error) {
	err = c.storage.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(c.namespace))
		if b == nil {
			return ErrorNotFound
		}
		c := b.Cursor()
		if k, v := c.Seek([]byte(key)); bytes.Equal(k, []byte(key)) {
			if v == nil {
				return ErrorNotFound
			}

			var decodedVal T
			d := gob.NewDecoder(bytes.NewReader(v))
			if err := d.Decode(&decodedVal); err != nil {
				return err
			}

			val = decodedVal
		} else {
			return ErrorNotFound
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
