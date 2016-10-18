package chord

import (
	"math/big"
)

type Node interface  {
	// BigInt turns Node ID into big.Int representation.
	BigInt() *big.Int

	// Bits returns amount of significant bits in Node ID.
	Bits() int

	// Cmp compares ID of this Node with given ID.
	//
	// Returns -1, 0 or 1 depending on if given other ID is lesser than, equal
	// to, or greater than this Node's ID.
	Cmp(other ID) int

	// Diff calculates the difference between this Node's ID and given other ID.
	Diff(other ID) ID

	// Eq determines if this Node's ID and given other ID are equal.
	Eq(other ID) bool

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
