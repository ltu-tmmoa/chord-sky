package chord

import "net"

// Node represents some Chord node, available either locally or remotely.
type Node interface {
	// ID returns node ID.
	ID() ID

	// IPAddr provides node network address.
	IPAddr() *net.IPAddr

	// Finger resolves provides finger interval at provided offset i.
	//
	// The result is only defined for i in [1,M], where M is the amount of bits set
	// at node ring creation.
	Finger(i int) *Finger

	// FingerNode resolves Chord node at given finger table offset i.
	//
	// The result is only defined for i in [1,M], where M is the amount of bits set
	// at node ring creation.
	FingerNode(i int) (Node, error)

	// SetFingerNode attempts to set this node's ith finger to given node.
	//
	// The operation is only valid for i in [1,M], where M is the amount of
	// bits set at node ring creation.
	SetFingerNode(i int, fing Node) error

	// Successor yields the next node in this node's ring.
	Successor() (Node, error)

	// Predecessor yields the previous node in this node's ring.
	Predecessor() (Node, error)

	// FindSuccessor asks this node to find successor of given ID.
	FindSuccessor(id ID) (Node, error)

	// FindPredecessor asks this node to find a predecessor of given ID.
	FindPredecessor(id ID) (Node, error)

	// SetSuccessor attempts to set this node's successor to given node.
	SetSuccessor(succ Node) error

	// SetPredecessor attempts to set this node's predecessor to given node.
	SetPredecessor(pred Node) error

	// String turns Node into its canonical string representation.
	String() string
}
