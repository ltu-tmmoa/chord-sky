package data

import (
	"crypto/sha1"
	"fmt"
)

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
func (storage *MemoryStorage) Get(key *ID) ([]byte, error) {
	return storage.data[key.String()], nil
}

// GetKeyRange gets all keys that lexically located within [fromKey, toKey).
func (storage *MemoryStorage) GetKeyRange(fromKey, toKey *ID) ([]*ID, error) {
	keys := []*ID{}
	for skey := range storage.data {
		key, ok := ParseID(skey, fromKey.Bits())
		if !ok {
			panic(fmt.Sprint("Illegal key in memory storage:", skey))
		}
		if IDIntervalContainsIE(fromKey, toKey, key) {
			keys = append(keys, key)
		}
	}

	return keys, nil
}

// Set stores provided key/value pair, potentially replacing an existing
// such.
func (storage *MemoryStorage) Set(key *ID, value []byte) error {
	storage.data[key.String()] = value
	return nil
}

// Remove attempts to remove one key/value pair from store with a key
// matching given.
func (storage *MemoryStorage) Remove(key *ID) error {
	delete(storage.data, key.String())
	return nil
}

// GetKeys gets all keys that are located on this storage node.
func (storage *MemoryStorage) GetAllKeys() ([]*ID, error) {
	keys := []*ID{}
	for skey := range storage.data {
		key, ok := parseID(skey)
		if !ok {
			panic(fmt.Sprint("Illegal key in memory storage:", skey))
		}
		keys = append(keys, key)
	}
	return keys, nil
}

// Haidar's fix ;)
const (
	idBits = sha1.Size * 8
)

func parseID(s string) (*ID, bool) {
	return ParseID(s, idBits)
}
