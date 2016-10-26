package chord

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/ltu-tmmoa/chord-sky/log"
)

const (
	schemeHTTP = "http"
)

// RemoteNode represents some Chord node available remotely.
type RemoteNode struct {
	tcpAddr net.TCPAddr
	id      ID
}

// NewRemoteNode creates a new remote node from given address.
func NewRemoteNode(tcpAddr *net.TCPAddr) *RemoteNode {
	return newRemoteNode(tcpAddr, identity(tcpAddr, HashBitsMax))
}

func newRemoteNode(tcpAddr *net.TCPAddr, id *ID) *RemoteNode {
	node := new(RemoteNode)
	node.tcpAddr = *tcpAddr
	node.id = *id
	return node
}

// ID returns node ID.
func (node *RemoteNode) ID() *ID {
	return &node.id
}

// TCPAddr provides node network address.
func (node *RemoteNode) TCPAddr() *net.TCPAddr {
	return &node.tcpAddr
}

// FingerStart resolves start ID of finger table entry i.
//
// The result is only defined for i in [1,M], where M is the amount of bits set
// at node ring creation.
func (node *RemoteNode) FingerStart(i int) *ID {
	m := node.ID().Bits()
	verifyIndexOrPanic(m, i)
	return calcFingerStart(node.ID(), i)
}

// FingerNode resolves Chord node at given finger table offset i.
//
// The result is only defined for i in [1,M], where M is the amount of bits set
// at node ring creation.
func (node *RemoteNode) FingerNode(i int) (Node, error) {
	u, err := url.Parse("node/fingers")
	u.Host = node.TCPAddr().String()
	u.Scheme = schemeHTTP
	if err != nil {
		log.Logger.Fatal(err)
	}

	q := u.Query()
	q.Set("id", strconv.Itoa(i))
	u.RawQuery = q.Encode()

	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	tcpAddr, err := net.ResolveTCPAddr("ip", string(body))
	if err != nil {
		return nil, err
	}
	return NewRemoteNode(tcpAddr), nil
}

// SetFingerNode attempts to set this node's ith finger to given node.
//
// The operation is only valid for i in [1,M], where M is the amount of
// bits set at node ring creation.
func (node *RemoteNode) SetFingerNode(i int, fing Node) error {
	u, err := url.Parse(fmt.Sprintf("node/fingers?id=%d", i))
	u.Host = node.TCPAddr().String()
	u.Scheme = schemeHTTP
	if err != nil {
		log.Logger.Fatal(err)
	}

	body := strings.NewReader(fing.TCPAddr().IP.String())
	req, err := http.NewRequest(http.MethodPut, u.String(), body)
	if err != nil {
		return err
	}

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	return nil
}

// RemoveFingerNodesByID attempts to remove all nodes from this node's
// finger table that match given ID.
func (node *RemoteNode) RemoveFingerNodesByID(id *ID) error {
	url := fmt.Sprintf("http://%s:8080/node/fingers?id=%s", node.TCPAddr().IP.String(), id.String())
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	return nil
}

// Successor yields the next node in this node's ring.
func (node *RemoteNode) Successor() (Node, error) {
	u, err := url.Parse("node/successor")
	u.Host = node.TCPAddr().String()
	u.Scheme = schemeHTTP
	if err != nil {
		log.Logger.Fatal(err)
	}

	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	tcpAddr, err := net.ResolveTCPAddr("ip", string(body))
	if err != nil {
		return nil, err
	}
	return NewRemoteNode(tcpAddr), nil
}

// Predecessor yields the previous node in this node's ring.
func (node *RemoteNode) Predecessor() (Node, error) {
	u, err := url.Parse("node/predecessor")
	u.Host = node.TCPAddr().String()
	u.Scheme = schemeHTTP
	if err != nil {
		log.Logger.Fatal(err)
	}

	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	tcpAddr, err := net.ResolveTCPAddr("ip", string(body))
	if err != nil {
		return nil, err
	}
	return NewRemoteNode(tcpAddr), nil
}

// FindSuccessor asks this node to find successor of given ID.
func (node *RemoteNode) FindSuccessor(id *ID) (Node, error) {
	u, err := url.Parse("node/successors")
	u.Host = node.TCPAddr().String()
	u.Scheme = schemeHTTP
	if err != nil {
		log.Logger.Fatal(err)
	}

	q := u.Query()
	q.Set("id", id.String())
	u.RawQuery = q.Encode()

	res, err := http.Get(u.String())
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	tcpAddr, err := net.ResolveTCPAddr("ip", string(body))
	if err != nil {
		return nil, err
	}

	return NewRemoteNode(tcpAddr), nil
}

// FindPredecessor asks this node to find a predecessor of given ID.
func (node *RemoteNode) FindPredecessor(id *ID) (Node, error) {
	u, err := url.Parse("node/predecessors")
	u.Host = node.TCPAddr().String()
	u.Scheme = schemeHTTP
	if err != nil {
		log.Logger.Fatal(err)
	}

	q := u.Query()
	q.Set("id", id.String())
	u.RawQuery = q.Encode()

	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	tcpAddr, err := net.ResolveTCPAddr("ip", string(body))
	if err != nil {
		return nil, err
	}

	return NewRemoteNode(tcpAddr), nil
}

// SetSuccessor attempts to set this node's successor to given node.
func (node *RemoteNode) SetSuccessor(successor Node) error {
	u, err := url.Parse("node/successor")
	u.Host = node.TCPAddr().String()
	u.Scheme = schemeHTTP
	if err != nil {
		log.Logger.Fatal(err)
	}

	body := strings.NewReader(successor.TCPAddr().IP.String())
	req, err := http.NewRequest(http.MethodPut, u.String(), body)
	if err != nil {
		return err
	}

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	return nil
}

// SetPredecessor attempts to set this node's predecessor to given node.
func (node *RemoteNode) SetPredecessor(predecessor Node) error {
	u, err := url.Parse("node/predecessor")
	u.Host = node.TCPAddr().String()
	u.Scheme = schemeHTTP
	if err != nil {
		log.Logger.Fatal(err)
	}

	body := strings.NewReader(predecessor.TCPAddr().IP.String())
	req, err := http.NewRequest(http.MethodPut, u.String(), body)
	if err != nil {
		return err
	}

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	return nil
}

// String turns Node into its canonical string representation.
func (node *RemoteNode) String() string {
	return fmt.Sprintf("%s %s", node.id.String(), node.tcpAddr.String())
}
