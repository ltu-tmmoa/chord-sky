package chord

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ltu-tmmoa/chord-sky/data"
)

type remoteStorage struct {
	node *remoteNode
}

func newRemoteStorage(node *remoteNode) *remoteStorage {
	return &remoteStorage{
		node: node,
	}
}

// Get attempts to get value associated with given key.
//
// Acquiring a value of `nil` is not considered an error.
func (storage *remoteStorage) Get(key *data.ID) ([]byte, error) {
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
func (storage *remoteStorage) GetKeyRange(fromKey, toKey *data.ID) ([]*data.ID, error) {
	// TODO

	return []*data.ID{}, nil
}

// Set stores provided key/value pair, potentially replacing an existing
// such.
func (storage *remoteStorage) Set(key *data.ID, value []byte) error {
	// TODO

	return nil
}

// Remove attempts to remove one key/value pair from store with a key
// matching given.
func (storage *remoteStorage) Remove(key *data.ID) error {
	// TODO

	return nil
}
