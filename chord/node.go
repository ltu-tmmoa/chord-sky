package chord

import "net"

// Node represents some Chord node, available either locally or remotely.
type Node interface {
	// ID returns node ID.
	ID() ID

	// IPAddr provides node network address.
	IPAddr() *net.IPAddr

	// Finger resolves Chord node at given finger table offset i.
	//
	// The result is only defined for i in [1,M], where M is the amount of bits set
	// at node ring creation.
	Finger(i int) *Finger

	// Successor yields the next node in this node's ring.
	Successor() (Node, error)

	// Predecessor yields the previous node in this node's ring.
	Predecessor() (Node, error)

	// FindSuccessor asks this node to find successor of given ID.
	FindSuccessor(id ID) (Node, error)

	// FindPredecessor asks this node to find a predecessor of given ID.
	FindPredecessor(id ID) (Node, error)

	// SetSuccessor attempts to set this node's successor to given node.
	SetSuccessor(successor Node) error

	// SetPredecessor attempts to set this node's predecessor to given node.
	SetPredecessor(predecessor Node) error

	// String turns Node into its canonical string representation.
	String() string
}
