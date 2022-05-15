package cache

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	bolt "go.etcd.io/bbolt"
)

const dbName = "ignite_cache.db"

var ErrorNotFound = errors.New("no value was found with the provided key")

// Storage is meant to be passed around and used by the New function (which provides namespacing and type-safety)
type Storage struct {
	db *bolt.DB
}

// Cache is a namespaced and type-safe key-value store
type Cache[T any] struct {
	storage   Storage
	namespace string
}

// NewStorage create a local db file (if it doesn't exist) and opens a connection to it
// Storage needs to be closed before another is opened.
// If another process already has the same db file open, NewStorage will wait for it to be closed
func NewStorage(dir string) (Storage, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return Storage{}, err
	}

	storagePath := filepath.Join(dir, dbName)
	fmt.Println("Opening db")
	db, err := bolt.Open(storagePath, 0640, &bolt.Options{Timeout: 5 * time.Minute})
	if err != nil {
		return Storage{}, err
	}
	fmt.Println("db opened")

	return Storage{db}, nil
}

// Close closes the underlying database
// Attempting to use the same Storage instance or any pre-existing Cache instances
// will result in errors.
func (s Storage) Close() error {
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
