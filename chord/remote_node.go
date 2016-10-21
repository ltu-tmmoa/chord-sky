package chord

import (
	"math/big"
	"net"
	  "net/http"
	  "io/ioutil"
	  "net/url"
	  "github.com/ltu-tmmoa/chord-sky/log"
	  "strconv"
)

// RemoteNode represents some Chord node available remotely.
type RemoteNode struct {
	ipAddr 		net.IPAddr
	id			Hash
	fingers     	[]*Finger
	predecessor 	Node
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

// BigInt turns Node ID into big.Int representation.
func (node *RemoteNode) BigInt() *big.Int {
	return node.id.BigInt()
}

// Bits returns amount of significant bits in Node ID.
func (node *RemoteNode) Bits() int {
	return node.id.Bits()
}

// Cmp compares ID of this Node with given ID.
//
// Returns -1, 0 or 1 depending on if given other ID is lesser than, equal
// to, or greater than this Node's ID.
func (node *RemoteNode) Cmp(other ID) int {
	return node.id.Cmp(other)
}

// Diff calculates the difference between this Node's ID and given other ID.
func (node *RemoteNode) Diff(other ID) ID {
	return node.id.Diff(other)
}

// Eq determines if this Node's ID and given other ID are equal.
func (node *RemoteNode) Eq(other ID) bool {
	return node.id.Eq(other)
}

// Hash turns this Node's ID into Hash representation.
func (node *RemoteNode) Hash() Hash {
	return node.id.Hash()
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
	  u, err := url.Parse(node.ipAddr.String()+"/node/finger")
	  if err != nil { log.Logger.Fatal(err) }

	  q := u.Query()
	  q.Set("finger", strconv.Itoa(i))
	  u.RawQuery = q.Encode()

	  resp, err := http.Get(u.String())
	  defer resp.Body.Close()
	  if err != nil { return nil }

	  body, err := ioutil.ReadAll(resp.Body)
	  addr, err := net.ResolveIPAddr("ip", string(body))
	  if err != nil { return nil }
	  node.fingers[i].SetNodeFromIPAddress(addr)
	  return node.fingers[i]
}

// Successor yields the next node in this node's ring.
func (node *RemoteNode) Successor() (Node, error) {

	  u, err := url.Parse(node.ipAddr.String()+"/node/successor")
	  if err != nil { log.Logger.Fatal(err) }

	  resp, err := http.Get(u.String())
	  defer resp.Body.Close()
	  if err != nil { return nil, err }

	  body, err := ioutil.ReadAll(resp.Body)
	  addr, err := net.ResolveIPAddr("ip", string(body))
	  if err != nil { return nil, err }

	  return NewRemoteNode(addr), nil
}


// Predecessor yields the previous node in this node's ring.
func (node *RemoteNode) Predecessor() (Node, error) {
	  u, err := url.Parse(node.ipAddr.String()+"/node/predecessor")
	  if err != nil { log.Logger.Fatal(err) }

	  resp, err := http.Get(u.String())
	  defer resp.Body.Close()
	  if err != nil { return nil, err }

	  body, err := ioutil.ReadAll(resp.Body)
	  addr, err := net.ResolveIPAddr("ip", string(body))
	  if err != nil { return nil, err }

	  return NewRemoteNode(addr), nil
}

// FindSuccessor asks this node to find successor of given ID.
func (node *RemoteNode) FindSuccessor(id ID) (Node, error) {
	  u, err := url.Parse(node.ipAddr.String()+"/node/successors")
	  if err != nil { log.Logger.Fatal(err) }
	  q := u.Query()
	  q.Set("id", id.String())
	  u.RawQuery = q.Encode()

	  resp, err := http.Get(u.String())
	  defer resp.Body.Close()
	  if err != nil { return nil, err }

	  body, err := ioutil.ReadAll(resp.Body)
	  addr, err := net.ResolveIPAddr("ip", string(body))
	  if err != nil { return nil, err }

	  return NewRemoteNode(addr), nil
}

// FindPredecessor asks this node to find a predecessor of given ID.
func (node *RemoteNode) FindPredecessor(id ID) (Node, error) {
	  u, err := url.Parse(node.ipAddr.String()+"/node/predecessors")
	  if err != nil { log.Logger.Fatal(err) }

	  q := u.Query()
	  q.Set("id", id.String())
	  u.RawQuery = q.Encode()

	  resp, err := http.Get(u.String())
	  defer resp.Body.Close()
	  if err != nil { return nil, err }

	  body, err := ioutil.ReadAll(resp.Body)
	  addr, err := net.ResolveIPAddr("ip", string(body))
	  if err != nil { return nil, err }

	  return NewRemoteNode(addr), nil
}

// SetSuccessor attempts to set this node's successor to given node.
func (node *RemoteNode) SetSuccessor(successor Node) error {
	  u, err := url.Parse(node.ipAddr.String()+"/node/successor")
	  if err != nil { log.Logger.Fatal(err) }

	  q := u.Query()
	  q.Set("ip", successor.IPAddr().String())
	  u.RawQuery = q.Encode()

	  resp, err := http.NewRequest("PUT", u.String(), nil)
	  defer resp.Body.Close()
	  if err != nil { return err }
	  return nil
}

// SetPredecessor attempts to set this node's predecessor to given node.
func (node *RemoteNode) SetPredecessor(predecessor Node) error {
	  u, err := url.Parse(node.ipAddr.String()+"/node/predecessor")
	  if err != nil { log.Logger.Fatal(err) }

	  q := u.Query()
	  q.Set("ip", predecessor.IPAddr().String())
	  u.RawQuery = q.Encode()

	  resp, err := http.NewRequest("PUT", u.String(), nil)
	  defer resp.Body.Close()
	  if err != nil { return err }
	  return nil
}


// String turns Node into its canonical string representation.
func (node *RemoteNode) String() string {
	  return node.id.String()
}