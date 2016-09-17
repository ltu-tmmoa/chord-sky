package chord

import (
	"fmt"
	"net"
)

// Ring holds a circular list of Chord nodes.
type Ring struct {
	root *Node
	bits int
}

// NewRing creates a new empty Chord node ring with a node capacity of 2^bits.
func NewRing(bits int) *Ring {
	if bits < 1 {
		panic("bits < 1")
	}
	ring := new(Ring)
	ring.root = nil
	ring.bits = bits
	return ring
}

// AddNode adds Chord node with given address to ring.
func (ring *Ring) AddNode(addr net.Addr) {
	node := newNode(addr, hash(addr, &ring.bits))
	// TODO: Make sure new node is added at correct location in ring.
	if ring.root == nil {
		ring.root = node
	}
	node.predecessor = ring.root
	ring.root = node
}

// Print outputs the ring members to console.
func (ring *Ring) Print() {
	node0 := ring.root
	for ; node0 != nil && !ring.root.id.Eq(node0.id); node0 = node0.Successor() {
		fmt.Println(node0.String())
	}
}

// Root provides the first node added to the ring.
func (ring *Ring) Root() *Node {
	return ring.root
}
