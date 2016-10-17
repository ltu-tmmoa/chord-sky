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
	Successor() Node

	// Predecessor yields the previous node in this node's ring.
	Predecessor() Node

	// FindSuccessor asks this node to find successor of given ID.
	FindSuccessor(id ID) Node

	// Join makes this node join the ring of given other node.
	//
	// If given node is nil, this node will form its own ring.
	Join(node0 Node)

	// Stabilize attempts to fix any ring issues arising from joining or leaving Chord ring nodes.
	//
	// Recommended to be called periodically in order to ensure node data integrity.
	Stabilize()

	// FixFingers refreshes this node's finger table entries in relation to Chord ring changes.
	//
	// Recommended to be called periodically in order to ensure finger table integrity.
	FixFingers()

	// FixAllFingers refreshes all of this node's finger table entries.
	FixAllFingers()

	// String turns Node into its canonical string representation.
	String() string
}
