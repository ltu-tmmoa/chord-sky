package data

// Storage of string keys and byte array values.
type Storage interface {
	// Get attempts to get value associated with given key.
	//
	// Acquiring a value of `nil` is not considered an error.
	Get(key string) ([]byte, error)

	// GetKeyRange gets all keys that lexically located within [fromKey, toKey).
	GetKeyRange(fromKey, toKey string) ([]string, error)

	// Set stores provided key/value pair, potentially replacing an existing
	// such.
	Set(key string, value []byte) error

	// Remove attempts to remove one key/value pair from store with a key
	// matching given.
	Remove(key string) error
}
