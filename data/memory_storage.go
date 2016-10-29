package data

// MemoryStorage provides in-memory storage.
type MemoryStorage struct {
	data map[string][]byte
}

// Get attempts to get value associated with given key, if any.
//
// Acquiring a value of `nil` is not considered an error.
func (storage *MemoryStorage) Get(key string) ([]byte, error) {
	return storage.data[key], nil
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
