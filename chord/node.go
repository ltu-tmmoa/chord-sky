package chord

import "net"

// Node represents some Chord node, available either locally or remotely.
type Node interface {
	// ID returns node ID.
	ID() *ID

	// TCPAddr provides node network address.
	TCPAddr() *net.TCPAddr

	// fingerStart resolves start ID of finger table entry i.
	//
	// The result is only defined for i in [1,M], where M is the amount of bits
	// set at node ring creation.
	FingerStart(i int) *ID

	// fingerNode resolves Chord node at given finger table offset i.
	//
	// The result is only defined for i in [1,M], where M is the amount of bits
	// set at node ring creation.
	FingerNode(i int) <-chan Node

	// SetfingerNode attempts to set this node's ith finger to given node.
	//
	// The operation is only valid for i in [1,M], where M is the amount of
	// bits set at node ring creation.
	SetfingerNode(i int, fing Node) <-chan *struct{}

	// Successor yields the next node in this node's ring.
	Successor() <-chan Node

	// Predecessor yields the previous node in this node's ring.
	Predecessor() <-chan Node

	// FindSuccessor asks this node to find successor of given ID.
	FindSuccessor(id *ID) <-chan Node

	// FindPredecessor asks this node to find a predecessor of given ID.
	FindPredecessor(id *ID) <-chan Node

	// SetSuccessor attempts to set this node's successor to given node.
	SetSuccessor(succ Node) <-chan *struct{}

	// SetPredecessor attempts to set this node's predecessor to given node.
	SetPredecessor(pred Node) <-chan *struct{}

	// String turns Node into its canonical string representation.
	String() string
}
