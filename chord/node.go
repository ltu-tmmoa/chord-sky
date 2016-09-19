package chord

import (
	"fmt"
	"math/big"
	"net"
)

// Node represents a potential member of a Chord ring.
type Node struct {
	addr        net.Addr
	id          *Hash
	fingers     []*Finger
	predecessor *Node
}

func newNode(addr net.Addr, id *Hash) *Node {
	node := new(Node)
	node.addr = addr
	node.id = id

	fingers := make([]*Finger, id.bits)
	for i := range fingers {
		fingers[i] = newFinger(id, i+1)
	}
	node.fingers = fingers
	node.predecessor = nil
	return node
}

// NewNodeRing creates a new node in its own ring.
//
// The created ring has a capacity of 2^bits node members.
func NewNodeRing(addr net.Addr, bits int) *Node {
	node := newNode(addr, hash(addr, bits))
	node.Join(nil)
	return node
}

// AsInt returns node identifier as a big.Int.
func (node *Node) AsInt() *big.Int {
	return node.id.AsInt()
}

// Bits returns amount of significant bits in node identifier.
func (node *Node) Bits() int {
	return node.id.Bits()
}

// Cmp compares this node's identifier to given ID.
func (node *Node) Cmp(other ID) int {
	return node.id.Cmp(other)
}

// Diff calculates the difference between this node's identifier and given ID.
func (node *Node) Diff(other ID) ID {
	return node.id.Diff(other)
}

// Eq determines if this node's identifier and given ID are equal.
func (node *Node) Eq(other ID) bool {
	return node.id.Eq(other)
}

// Finger resolves Chord node at given finger table offset i.
//
// The result is only defined for i in [1,M], where M is the amount of bits set
// at node ring creation.
func (node *Node) Finger(i int) *Finger {
	if 1 > i || i > node.id.bits {
		panic(fmt.Sprintf("%d not in [1,%d]", i, node.id.bits))
	}
	return node.finger(i)
}

func (node *Node) finger(i int) *Finger {
	return node.fingers[i-1]
}

// Successor yields the next node in this node's ring.
func (node *Node) Successor() *Node {
	return node.finger(1).node
}

// Predecessor yields the previous node in this node's ring.
func (node *Node) Predecessor() *Node {
	return node.predecessor
}

// FindSuccessor asks this node to find successor of given ID.
func (node *Node) FindSuccessor(id ID) *Node {
	node0 := node.findPredecessor(id)
	return node0.Successor()
}

// Asks node to find id' predecessor.
//
// See Chord paper figure 4.
func (node *Node) findPredecessor(id ID) *Node {
	node0 := node
	for id.Cmp(node0) <= 0 && id.Cmp(node0.Successor()) > 0 {
		node0 = node0.closestPrecedingFinger(id)
	}
	return node0
}

// Returns closest finger preceding ID.
//
// See Chord paper figure 4.
func (node *Node) closestPrecedingFinger(id ID) *Node {
	for i := node.Bits(); i > 0; i-- {
		if f := node.finger(i).node; f.Cmp(node) > 0 && f.Cmp(id) < 0 {
			return f
		}
	}
	return node
}

// Join makes this node join the ring of given other node.
//
// If given node is nil, this node will form its own ring.
func (node *Node) Join(other *Node) {
	if other != nil {
		if node.Bits() != other.Bits() {
			node.id = hash(node.addr, other.Bits())
		}
		node.initFingerTable(other)
		node.updateOthers()
		// TODO: Move keys in (predecessor,n] from successor
	} else {
		m := node.Bits()
		for i := 1; i <= m; i++ {
			node.finger(i).node = node
		}
		node.predecessor = node
	}
}

// Initialize finger table of local node; node0 is an arbitary node already in
// the network.
//
// See Chord paper figure 6.
func (node *Node) initFingerTable(node0 *Node) {
	// Add this node to node0 node's ring.
	{
		successor := node0.FindSuccessor(node.finger(1).start)
		node.finger(1).node = successor
		node.predecessor = successor.predecessor
		successor.predecessor = node
	}
	// Update this node's finger table.
	{
		m := node.Bits()
		for i := 1; i < m; i++ {
			this := node.finger(i)
			next := node.finger(i + 1)
			if next.start.Cmp(node) >= 0 && next.start.Cmp(this.node) < 0 {
				next.node = this.node
			} else {
				next.node = node0.FindSuccessor(next.start)
			}
		}
	}
}

// Update all nodes whose finger tables should refer to node.
//
// See Chord paper figure 6.
func (node *Node) updateOthers() {
	m := node.Bits()
	for i := 1; i <= m; i++ {
		var id ID
		{
			subtrahend := big.Int{}
			subtrahend.SetInt64(2)
			subtrahend.Exp(&subtrahend, big.NewInt(int64(i-1)), nil)
			id = node.Diff(newHash(subtrahend, m))
		}
		predecessor := node.findPredecessor(id)
		predecessor.updateFingerTable(node, i)
	}
}

// If s is the i:th finger of node, update node's finger table with s.
//
// See Chord paper figure 6.
func (node *Node) updateFingerTable(s *Node, i int) {
	finger := node.finger(i)
	if s.Cmp(node) >= 0 && s.Cmp(finger.node) < 0 {
		finger.node = s
		predecessor := node.predecessor
		predecessor.updateFingerTable(s, i)
	}
}

// PrintRing outputs this node's ring to console.
func (node *Node) PrintRing() {
	node0 := node
	for ; node0 != nil && !node.Eq(node0.id); node0 = node0.Successor() {
		fmt.Println(node0.String())
	}
}

// String produces canonical string representation of this node.
func (node *Node) String() string {
	return node.id.String()
}
