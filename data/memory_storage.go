package data

import "strings"

// MemoryStorage provides in-memory storage.
type MemoryStorage struct {
	data map[string][]byte
}

// NewMemoryStorage creates a new MemoryStorage instance.
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: map[string][]byte{},
	}
}

// Get attempts to get value associated with given key, if any.
//
// Acquiring a value of `nil` is not considered an error.
func (storage *MemoryStorage) Get(key string) ([]byte, error) {
	return storage.data[key], nil
}

// GetKeyRange gets all keys that lexically located within [fromKey, toKey).
func (storage *MemoryStorage) GetKeyRange(fromKey, toKey string) ([]string, error) {
	keys := []string{}
	for key := range storage.data {
		cmp0 := strings.Compare(fromKey, key)
		cmp1 := strings.Compare(toKey, key)
		if cmp0 <= 0 && cmp1 > 0 {
			keys = append(keys, key)
		}
	}
	return keys, nil
}

// Set stores provided key/value pair, potentially replacing an existing
// such.
func (storage *MemoryStorage) Set(key string, value []byte) error {
	storage.data[key] = value
	return nil
}

// Remove attempts to remove one key/value pair from store with a key
// matching given.
func (storage *MemoryStorage) Remove(key string) error {
	delete(storage.data, key)
	return nil
}
