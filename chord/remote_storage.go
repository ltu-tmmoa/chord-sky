package chord

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"github.com/ltu-tmmoa/chord-sky/data"
	  "encoding/base64"
	  "net/url"
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

	url := fmt.Sprintf("http://%s/storage/%s", node.TCPAddr().String(), key.String())
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
	  node := storage.node

	  // http://<IP:PORT>/storage?fromKey=x00&toKey=x11
	  url, err := url.Parse(fmt.Sprintf("http://%s/storage", node.TCPAddr().String()))
	  if err != nil {
		    node.disconnect(err)
		    return nil, err
	  }
	  q := url.Query()
	  q.Set("fromKey", fromKey.String())
	  q.Set("toKey",toKey.String())

	  res, err := http.Get(url.String())
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

// Set stores provided key/value pair, potentially replacing an existing
// such.
func (storage *remoteStorage) Set(key *data.ID, value []byte) error {
	  node := storage.node

	  url := fmt.Sprintf("http://%s/storage/%s", node.TCPAddr().String(), key.String())

	  // Base64 encoding, RFC 4648.
	  str := base64.StdEncoding.EncodeToString(value)
	  req, err := http.NewRequest(http.MethodPut, url, str)
	  if req.Body != nil {
		    defer req.Body.Close()
	  }
	  if err != nil {
		    node.disconnect(err)
		    return err
	  }
	  res, err := http.DefaultClient.Do(req)
	  if err != nil {
		    node.disconnect(err)
		    return err
	  }
	  if res.Body != nil {
		    defer res.Body.Close()
	  }
	  if res.StatusCode < 200 || res.StatusCode > 299 {
		    err := fmt.Errorf("HTTP storage Put %s -> %d %s", url, res.StatusCode, res.Status)
		    node.disconnect(err)
		    return err
	  }
	  return nil
}

// Remove attempts to remove one key/value pair from store with a key
// matching given.
func (storage *remoteStorage) Remove(key *data.ID) error {
	  node := storage.node

	  url := fmt.Sprintf("http://%s/storage/%s", node.TCPAddr().String(), key.String())
	  req, err := http.NewRequest(http.MethodDelete, url, nil)
	  if req.Body != nil {
		    defer req.Body.Close()
	  }
	  if err != nil {
		    node.disconnect(err)
		    return err
	  }
	  res, err := http.DefaultClient.Do(req)
	  if err != nil {
		    node.disconnect(err)
		    return err
	  }
	  if res.Body != nil {
		    defer res.Body.Close()
	  }
	  if res.StatusCode < 200 || res.StatusCode > 299 {
		    err := fmt.Errorf("HTTP storage Delete %s -> %d %s", url, res.StatusCode, res.Status)
		    node.disconnect(err)
		    return err
	  }
	  return nil
}
