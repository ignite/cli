package cache

// StorageOptions configures the cache storage.
type StorageOption func(*Storage)

// WithVersion sets a version for the storage.
// Version is used as prefix for any cached value.
func WithVersion(version string) StorageOption {
	return func(o *Storage) {
		o.version = version
	}
}
