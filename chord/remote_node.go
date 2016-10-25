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
	ipAddr net.IPAddr
	id     Hash
}

// NewRemoteNode creates a new remote node from given address.
func NewRemoteNode(ipAddr *net.IPAddr) *RemoteNode {
	return newRemoteNode(ipAddr, hash(ipAddr, HashBitsMax))
}

func newRemoteNode(ipAddr *net.IPAddr, id *Hash) *RemoteNode {
	node := new(RemoteNode)
	node.ipAddr = *ipAddr
	node.id = *id
	return node
}

// ID returns node ID.
func (node *RemoteNode) ID() ID {
	return &node.id
}

// IPAddr provides node network address.
func (node *RemoteNode) IPAddr() *net.IPAddr {
	return &node.ipAddr
}

// Finger resolves Chord node at given finger table offset i.
//
// The result is only defined for i in [1,M], where M is the amount of bits set
// at node ring creation.
func (node *RemoteNode) Finger(i int) *Finger {
	u, err := url.Parse("node/fingers")
	u.Host = fmt.Sprintf("%s:8080", node.IPAddr().String())
	u.Scheme = schemeHTTP
	if err != nil {
		log.Logger.Fatal(err)
	}

	q := u.Query()
	q.Set("id", strconv.Itoa(i))
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	defer resp.Body.Close()
	if err != nil {
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	addr, err := net.ResolveIPAddr("ip", string(body))
	if err != nil {
		return nil
	}
	finger := newFinger(node.ID(), i)
	finger.SetNodeFromIPAddress(addr)
	return finger
}

// Successor yields the next node in this node's ring.
func (node *RemoteNode) Successor() (Node, error) {
	u, err := url.Parse("node/successor")
	u.Host = fmt.Sprintf("%s:8080", node.IPAddr().String())
	u.Scheme = schemeHTTP
	if err != nil {
		log.Logger.Fatal(err)
	}

	resp, err := http.Get(u.String())
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	addr, err := net.ResolveIPAddr("ip", string(body))
	if err != nil {
		return nil, err
	}

	return NewRemoteNode(addr), nil
}

// Predecessor yields the previous node in this node's ring.
func (node *RemoteNode) Predecessor() (Node, error) {
	u, err := url.Parse("node/predecessor")
	u.Host = fmt.Sprintf("%s:8080", node.IPAddr().String())
	u.Scheme = schemeHTTP
	if err != nil {
		log.Logger.Fatal(err)
	}

	resp, err := http.Get(u.String())
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	addr, err := net.ResolveIPAddr("ip", string(body))
	if err != nil {
		return nil, err
	}

	return NewRemoteNode(addr), nil
}

// FindSuccessor asks this node to find successor of given ID.
func (node *RemoteNode) FindSuccessor(id ID) (Node, error) {
	u, err := url.Parse("node/successors")
	u.Host = fmt.Sprintf("%s:8080", node.IPAddr().String())
	u.Scheme = schemeHTTP
	if err != nil {
		log.Logger.Fatal(err)
	}

	q := u.Query()
	q.Set("id", id.String())
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	addr, err := net.ResolveIPAddr("ip", string(body))
	if err != nil {
		return nil, err
	}

	return NewRemoteNode(addr), nil
}

// FindPredecessor asks this node to find a predecessor of given ID.
func (node *RemoteNode) FindPredecessor(id ID) (Node, error) {
	u, err := url.Parse("node/predecessors")
	u.Host = fmt.Sprintf("%s:8080", node.IPAddr().String())
	u.Scheme = schemeHTTP
	if err != nil {
		log.Logger.Fatal(err)
	}

	q := u.Query()
	q.Set("id", id.String())
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	addr, err := net.ResolveIPAddr("ip", string(body))
	if err != nil {
		return nil, err
	}

	return NewRemoteNode(addr), nil
}

// SetSuccessor attempts to set this node's successor to given node.
func (node *RemoteNode) SetSuccessor(successor Node) error {
	u, err := url.Parse("node/successor")
	u.Host = fmt.Sprintf("%s:8080", node.IPAddr().String())
	u.Scheme = schemeHTTP
	if err != nil {
		log.Logger.Fatal(err)
	}

	body := strings.NewReader(successor.IPAddr().IP.String())
	req, err := http.NewRequest(http.MethodPut, u.String(), body)
	if err != nil {
		return err
	}

	res, err := (&http.Client{}).Do(req)
	defer res.Body.Close()
	return err
}

// SetPredecessor attempts to set this node's predecessor to given node.
func (node *RemoteNode) SetPredecessor(predecessor Node) error {
	u, err := url.Parse("node/predecessor")
	u.Host = fmt.Sprintf("%s:8080", node.IPAddr().String())
	u.Scheme = schemeHTTP
	if err != nil {
		log.Logger.Fatal(err)
	}

	body := strings.NewReader(predecessor.IPAddr().IP.String())
	req, err := http.NewRequest(http.MethodPut, u.String(), body)
	if err != nil {
		return err
	}

	res, err := (&http.Client{}).Do(req)
	defer res.Body.Close()
	return err
}

// String turns Node into its canonical string representation.
func (node *RemoteNode) String() string {
	return node.id.String()
}
