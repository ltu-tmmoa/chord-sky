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

func (node *RemoteNode) httpGetNodef(pathFormat string, pathArgs ...interface{}) <-chan Node {
	ch := make(chan Node, 1)
	onError := func(err error) {
		log.Logger.Print(err.Error())
		ch <- nil
	}
	go func() {
		path := fmt.Sprintf(pathFormat, pathArgs...)
		url := fmt.Sprintf("http://%s/node/%s", node.TCPAddr().String(), path)
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
		ch <- NewRemoteNode(addr)
	}()
	return ch
}

func (node *RemoteNode) httpPut(path, body string) {
	onError := func(err error) {
		log.Logger.Print(err.Error())
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
	}()
}
