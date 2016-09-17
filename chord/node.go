package chord

import "net"

// Node represents a potential member of a Chord ring.
type Node struct {
	addr        net.Addr
	id          *ID
	predecessor *Node
}

func newNode(addr net.Addr, id *ID) *Node {
	node := new(Node)
	node.addr = addr
	node.id = id
	node.predecessor = nil
	return node
}

// Finger resolves Chord node at given finger table offset i.
//
// The result is only defined for offset in [1,M).
func (node *Node) Finger(i int) *Finger {
	return newFinger(node, i)
}

// Successor yields the next node in this node's ring.
func (node *Node) Successor() *Node {
	return node.Finger(1).Node()
}

// Predecessor yields the previous node in this node's ring.
func (node *Node) Predecessor() *Node {
	return node.predecessor
}

// String produces canonical string representation of this node.
func (node *Node) String() string {
	return node.id.String()
}
