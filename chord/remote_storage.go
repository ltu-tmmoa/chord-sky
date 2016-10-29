package chord

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type remoteStorage struct {
	node *remoteNode
}

// Get attempts to get value associated with given key.
//
// Acquiring a value of `nil` is not considered an error.
func (storage *remoteStorage) Get(key string) ([]byte, error) {
	node := storage.node

	url := fmt.Sprintf("http://%s/storage/%s", node.TCPAddr(), key)
	res, err := http.Get(url)
	if err != nil {
		node.disconnect(err)
		return nil, err
	}
	if res.Body == nil {
		return nil, nil
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		node.disconnect(err)
		return nil, err
	}
	return body, nil
}

// GetKeyRange gets all keys that lexically located within [fromKey, toKey).
func (storage *remoteStorage) GetKeyRange(fromKey, toKey string) ([]string, error) {
	// TODO

	return []string{}, nil
}

// Set stores provided key/value pair, potentially replacing an existing
// such.
func (storage *remoteStorage) Set(key string, value []byte) error {
	// TODO

	return nil
}

// Remove attempts to remove one key/value pair from store with a key
// matching given.
func (storage *remoteStorage) Remove(key string) error {
	// TODO

	return nil
}
