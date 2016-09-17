package chord

import (
	"fmt"
	"math/big"
	"net"
)

// Node represents a potential member of a Chord ring.
type Node struct {
	addr        net.Addr
	id          *ID
	fingers     []*Finger
	predecessor *Node
}

func newNode(addr net.Addr, id *ID) *Node {
	node := new(Node)
	node.addr = addr
	node.id = id

	fingers := make([]*Finger, *id.bits)
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
	node := newNode(addr, hash(addr, &bits))
	node.Join(nil)
	return node
}

// Finger resolves Chord node at given finger table offset i.
//
// The result is only defined for i in [1,M], where M is the amount of bits set
// at node ring creation.
func (node *Node) Finger(i int) *Finger {
	if 1 > i || i > *node.id.bits {
		panic(fmt.Sprintf("%d not in [1,%d]", i, *node.id.bits))
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

// FindSuccessor finds successor node of given ID.
func (node *Node) FindSuccessor(id *ID) *Node {
	node0 := node.findPredecessor(id)
	return node0.Successor()
}

func (node *Node) findPredecessor(id *ID) *Node {
	node0 := node
	for id.cmp(node0.id) <= 0 && id.cmp(node0.Successor().id) > 0 {
		node0 = node0.closestPrecedingFinger(id)
	}
	return node0
}

func (node *Node) closestPrecedingFinger(id *ID) *Node {
	for i := *node.id.bits; i > 0; i-- {
		if f := node.finger(i).node; f.id.cmp(node.id) > 0 && f.id.cmp(id) < 0 {
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
		if *node.id.bits != *other.id.bits {
			node.id = hash(node.addr, other.id.bits)
		}
		node.initFingerTable(other)
		node.updateOthers()
		// TODO: Move keys in (predecessor,n] from successor
	} else {
		m := *node.id.bits
		for i := 1; i <= m; i++ {
			node.finger(i).node = node
		}
		node.predecessor = node
	}
}

func (node *Node) initFingerTable(other *Node) {
	// Add this node to other node's ring.
	{
		successor := other.FindSuccessor(node.finger(1).start)
		node.finger(1).node = successor
		node.predecessor = successor.predecessor
		successor.predecessor = node
	}
	// Update this node's finger table.
	{
		m := *node.id.bits
		for i := 1; i < m; i++ {
			this := node.finger(i)
			next := node.finger(i + 1)
			if next.start.cmp(node.id) >= 0 && next.start.cmp(this.node.id) < 0 {
				next.node = this.node
			} else {
				next.node = other.FindSuccessor(next.start)
			}
		}
	}
}

func (node *Node) updateOthers() {
	m := *node.id.bits
	for i := 1; i <= m; i++ {
		subtrahend := big.NewInt(2)
		subtrahend.Exp(subtrahend, big.NewInt(int64(i-1)), nil)
		predecessor := node.findPredecessor(node.id.diff(subtrahend))
		predecessor.updateFingerTable(node, i)
	}
}

func (node *Node) updateFingerTable(s *Node, i int) {
	finger := node.finger(i)
	if s.id.cmp(node.id) >= 0 && s.id.cmp(finger.node.id) < 0 {
		finger.node = s
		predecessor := node.predecessor
		predecessor.updateFingerTable(s, i)
	}
}

// PrintRing outputs this node's ring to console.
func (node *Node) PrintRing() {
	node0 := node
	for ; node0 != nil && !node.id.eq(node0.id); node0 = node0.Successor() {
		fmt.Println(node0.String())
	}
}

// String produces canonical string representation of this node.
func (node *Node) String() string {
	return node.id.String()
}
