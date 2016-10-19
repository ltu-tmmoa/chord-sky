package chord

import (
	"math/big"
	"net"
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

// BigInt turns Node ID into big.Int representation.
func (node *RemoteNode) BigInt() *big.Int {
	return nil
}

// Bits returns amount of significant bits in Node ID.
func (node *RemoteNode) Bits() int {
	return -1
}

// Cmp compares ID of this Node with given ID.
//
// Returns -1, 0 or 1 depending on if given other ID is lesser than, equal
// to, or greater than this Node's ID.
func (node *RemoteNode) Cmp(other ID) int {
	return 0
}

// Diff calculates the difference between this Node's ID and given other ID.
func (node *RemoteNode) Diff(other ID) ID {
	return nil
}

// Eq determines if this Node's ID and given other ID are equal.
func (node *RemoteNode) Eq(other ID) bool {
	return false
}

// Hash turns this Node's ID into Hash representation.
func (node *RemoteNode) Hash() Hash {
	return node.id
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
	return nil
}

// Successor yields the next node in this node's ring.
func (node *RemoteNode) Successor() (Node, error) {
	return nil, nil
}

// Predecessor yields the previous node in this node's ring.
func (node *RemoteNode) Predecessor() (Node, error) {
	return nil, nil
}

// FindSuccessor asks this node to find successor of given ID.
func (node *RemoteNode) FindSuccessor(id ID) (Node, error) {
	return nil, nil
}

// FindPredecessor asks this node to find a predecessor of given ID.
func (node *RemoteNode) FindPredecessor(id ID) (Node, error) {
	return nil, nil
}

// SetSuccessor attempts to set this node's successor to given node.
func (node *RemoteNode) SetSuccessor(successor Node) error {
	return nil
}

// SetPredecessor attempts to set this node's predecessor to given node.
func (node *RemoteNode) SetPredecessor(predecessor Node) error {
	return nil
}

// String turns Node into its canonical string representation.
func (node *RemoteNode) String() string {
	return ""
}
