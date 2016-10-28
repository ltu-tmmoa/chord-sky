package chord

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/ltu-tmmoa/chord-sky/log"
)

func (node *remoteNode) httpHeartbeat(path string) {
	go func() {
		url := fmt.Sprintf("http://%s/node/%s", node.TCPAddr(), path)
		res, err := http.Get(url)
		if err != nil {
			node.disconnect(err)
			return
		}
		if res.Body == nil {
			node.disconnect(errors.New("No body in response."))
			return
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			node.disconnect(err)
			return
		}
		log.Logger.Println("Node", node, "heartbeat (", string(body), ").")
	}()
}

func (node *remoteNode) httpGetNodef(pathFormat string, pathArgs ...interface{}) <-chan NodeErr {
	return newChanNodeErr(func() (Node, error) {
		path := fmt.Sprintf(pathFormat, pathArgs...)
		url := fmt.Sprintf("http://%s/node/%s", node.TCPAddr(), path)
		res, err := http.Get(url)
		if err != nil {
			node.disconnect(err)
			return nil, err
		}
		if res.Body == nil {
			node.disconnect(errors.New("No body in response."))
			return nil, err
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			node.disconnect(err)
			return nil, err
		}
		addr, err := net.ResolveTCPAddr("tcp", string(body))
		if err != nil {
			node.disconnect(err)
			return nil, err
		}
		return node.pool.getOrCreateNode(addr), nil
	})
}

func (node *remoteNode) httpPut(path, body string) <-chan error {
	return newChanErr(func() error {
		url := fmt.Sprintf("http://%s/node/%s", node.TCPAddr().String(), path)
		req, err := http.NewRequest(http.MethodPut, url, bytes.NewBufferString(body))
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
			err := fmt.Errorf("HTTP PUT %s -> %d %s", url, res.StatusCode, res.Status)
			node.disconnect(err)
			return err
		}
		return nil
	})
}

func (node *remoteNode) disconnect(err error) {
	node.pool.removeNode(node)
	log.Logger.Printf("Node %s disconnected: %s", node.String(), err.Error())
}
