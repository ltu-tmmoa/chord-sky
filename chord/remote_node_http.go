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
	onError := func(err error) {
		node.pool.removeNode(node.TCPAddr())
		log.Logger.Printf("Node %s hearbeat failure: %s", node, err.Error())
	}
	go func() {
		url := fmt.Sprintf("http://%s/node/%s", node.TCPAddr(), path)
		res, err := http.Get(url)
		if err != nil {
			onError(err)
			return
		}
		if res.Body == nil {
			onError(errors.New("No body in response."))
			return
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			onError(err)
			return
		}
		log.Logger.Println("Node", node, "heartbeat (", string(body), ").")
	}()
}

func (node *remoteNode) httpGetNodef(pathFormat string, pathArgs ...interface{}) <-chan Node {
	ch := make(chan Node, 1)
	onError := func(err error) {
		node.pool.removeNode(node.TCPAddr())
		log.Logger.Printf("Node %s disconnected: %s", node.String(), err.Error())
		ch <- nil
	}
	go func() {
		path := fmt.Sprintf(pathFormat, pathArgs...)
		url := fmt.Sprintf("http://%s/node/%s", node.TCPAddr(), path)
		res, err := http.Get(url)
		if err != nil {
			onError(err)
			return
		}
		if res.Body == nil {
			onError(errors.New("No body in response."))
			return
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			onError(err)
			return
		}
		addr, err := net.ResolveTCPAddr("tcp", string(body))
		if err != nil {
			onError(err)
			return
		}
		ch <- node.pool.getOrCreateNode(addr)
	}()
	return ch
}

func (node *remoteNode) httpPut(path, body string) <-chan *struct{} {
	ch := make(chan *struct{}, 1)
	onError := func(err error) {
		node.pool.removeNode(node.TCPAddr())
		log.Logger.Printf("Node %s disconnected: %s", node.String(), err.Error())
		ch <- nil
	}
	go func() {
		url := fmt.Sprintf("http://%s/node/%s", node.TCPAddr().String(), path)
		req, err := http.NewRequest(http.MethodPut, url, bytes.NewBufferString(body))
		if err != nil {
			onError(err)
			return
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			onError(err)
			return
		}
		if res.Body != nil {
			defer res.Body.Close()
		}
		if res.StatusCode < 200 || res.StatusCode > 299 {
			onError(fmt.Errorf("HTTP PUT %s -> %d %s", url, res.StatusCode, res.Status))
		}
		ch <- nil
	}()
	return ch
}
