package chord

import (
	"net"

	"github.com/ltu-tmmoa/chord-sky/data"
)

// Node represents some Chord node, available either locally or remotely.
type Node interface {
	// ID returns node ID.
	ID() *data.ID

	// TCPAddr provides node network address.
	TCPAddr() *net.TCPAddr

	// fingerStart resolves start ID of finger table entry i.
	//
	// The result is only defined for i in [1,M], where M is the amount of bits
	// set at node ring creation.
	FingerStart(i int) *data.ID

	// fingerNode resolves Chord node at given finger table offset i.
	//
	// The result is only defined for i in [1,M], where M is the amount of bits
	// set at node ring creation.
	FingerNode(i int) (Node, error)

	// SetFingerNode attempts to set this node's ith finger to given node.
	//
	// The operation is only valid for i in [1,M], where M is the amount of
	// bits set at node ring creation.
	SetFingerNode(i int, fing Node) error

	// Successor yields the next node in this node's ring.
	Successor() (Node, error)

	// SuccessorList yields a list of nodes succeeding the current one.
	SuccessorList() ([]Node, error)

	// Predecessor yields the previous node in this node's ring.
	Predecessor() (Node, error)

	// FindSuccessor asks this node to find successor of given ID.
	FindSuccessor(id *data.ID) (Node, error)

	// FindPredecessor asks this node to find a predecessor of given ID.
	FindPredecessor(id *data.ID) (Node, error)

	// SetSuccessor attempts to set this node's successors to given node.
	SetSuccessor(succ Node) error

	// SetPredecessor attempts to set this node's predecessor to given node.
	SetPredecessor(pred Node) error

	// Storage exposes the data held by the node.
	Storage() data.Storage

	// String turns Node into its canonical string representation.
	String() string
}
