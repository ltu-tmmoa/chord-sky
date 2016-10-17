package chord

import (
	"fmt"
	"math/big"
	"math/rand"
	"net"
)

// Node represents a potential member of a Chord ring.
type LocalNode struct {
	addr        net.Addr
	id          *Hash
	fingers     []*Finger
	predecessor Node
}

func newLocalNode(addr net.Addr, id *Hash) *LocalNode {
	node := new(LocalNode)
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

// BigInt returns node identifier as a big.Int.
func (node *LocalNode) BigInt() *big.Int {
	return node.id.BigInt()
}

// Bits returns amount of significant bits in node identifier.
func (node *LocalNode) Bits() int {
	return node.id.Bits()
}

// Cmp compares this node's identifier to given ID.
func (node *LocalNode) Cmp(other ID) int {
	return node.id.Cmp(other)
}

// Diff calculates the difference between this node's identifier and given ID.
func (node *LocalNode) Diff(other ID) ID {
	return node.id.Diff(other)
}

// Eq determines if this node's identifier and given ID are equal.
func (node *LocalNode) Eq(other ID) bool {
	return node.id.Eq(other)
}

// Finger resolves Chord node at given finger table offset i.
//
// The result is only defined for i in [1,M], where M is the amount of bits set
// at node ring creation.
func (node *LocalNode) Finger(i int) *Finger {
	if 1 > i || i > node.id.bits {
		panic(fmt.Sprintf("%d not in [1,%d]", i, node.id.bits))
	}
	return node.finger(i)
}

func (node *LocalNode) finger(i int) *Finger {
	return node.fingers[i-1]
}

// Successor yields the next node in this node's ring.
func (node *LocalNode) Successor() Node {
	return node.finger(1).node
}

// Predecessor yields the previous node in this node's ring.
func (node *LocalNode) Predecessor() Node {
	return node.predecessor
}

// FindSuccessor asks this node to find successor of given ID.
//
// See Chord paper figure 4.
func (node *LocalNode) FindSuccessor(id ID) Node {
	node0 := node.FindPredecessor(id)
	return node0.Successor()
}

// FindPredecessor asks node to find id's predecessor.
//
// See Chord paper figure 4.
func (node *LocalNode) FindPredecessor(id ID) Node {
	return findPredecessor(node, id)
}

// Asks node to find id's predecessor.
//
// See Chord paper figure 4.
func findPredecessor(n Node, id ID) Node {
	n0 := n
	for !idIntervalContainsEI(n0, n0.Successor(), id) {
		n0 = closestPrecedingFinger(n0, id)
	}
	return n0
}

// Returns closest finger preceding ID.
//
// See Chord paper figure 4.
func closestPrecedingFinger(n Node, id ID) Node {
	for i := n.Bits(); i > 0; i-- {
		if f := n.Finger(i).node; idIntervalContainsEE(n, id, f) {
			return f
		}
	}
	return n
}

// SetSuccessor attempts to set this node's successor to given node.
func (node *LocalNode) SetSuccessor(successor Node) {
	node.finger(1).node = successor
}

// SetPredecessor attempts to set this node's predecessor to given node.
func (node *LocalNode) SetPredecessor(predecessor Node) {
	node.predecessor = predecessor
}

// Join makes this node join the ring of given other node.
//
// If given node is nil, this node will form its own ring.
func (node *LocalNode) Join(node0 Node) {
	if node0 != nil {
		if node.Bits() != node0.Bits() {
			node.id = hash(node.addr, node0.Bits())
		}
		node.initFingerTable(node0)
		updateOthers(node)
		// TODO: Move keys in (predecessor,n] from successor
	} else {
		m := node.Bits()
		for i := 1; i <= m; i++ {
			node.finger(i).node = node
		}
		node.predecessor = node
	}
}

// Update all nodes whose finger tables should refer to node.
//
// See Chord paper figure 6.
func updateOthers(n Node) {
	m := n.Bits()
	for i := 2; i <= m; i++ {
		var id ID
		{
			subtrahend := big.Int{}
			subtrahend.SetInt64(2)
			subtrahend.Exp(&subtrahend, big.NewInt(int64(i-1)), nil)
			id = n.Diff(newHash(subtrahend, m))
		}
		predecessor := n.FindPredecessor(id)
		updateFingerTable(predecessor, n, i)
	}
}

// If s is the i:th finger of node, update node's finger table with s.
//
// See Chord paper figure 6.
func updateFingerTable(n, s Node, i int) {
	finger := n.Finger(i)
	if idIntervalContainsIE(finger.Start(), finger.Node(), s) {
		finger.node = s
		predecessor := n.Predecessor()
		updateFingerTable(predecessor, s, i)
	}
}

// Initializes finger table of local node; node0 is an arbitrary node already in
// the network.
//
// See Chord paper figure 6.
func (node *LocalNode) initFingerTable(node0 Node) {
	// Add this node to node0 node's ring.
	{
		successor := node0.FindSuccessor(node.finger(1).Start())

		node.SetSuccessor(successor)
		node.SetPredecessor(successor.Predecessor())

		successor.Predecessor().SetSuccessor(node)
		successor.SetPredecessor(node)
	}
	// Update this node's finger table.
	{
		m := node.Bits()
		for i := 1; i < m; i++ {
			this := node.finger(i)
			next := node.finger(i + 1)
			if idIntervalContainsIE(node, this.Node(), next.Start()) {
				next.node = this.Node()
			} else {
				next.node = node0.FindSuccessor(next.Start())
			}
		}
	}
}

// Stabilize attempts to fix any ring issues arising from joining or leaving Chord ring nodes.
//
// Recommended to be called periodically in order to ensure node data integrity.
func (node *LocalNode) Stabilize() {
	x := node.Successor().Predecessor()
	if idIntervalContainsEE(node, node.Successor(), x) {
		node.SetSuccessor(x)
	}
	notify(node.Successor(), node)
}

func notify(n, n0 Node) {
	if n.Predecessor() == nil || idIntervalContainsEE(n.Predecessor(), n, n0) {
		n.SetPredecessor(n0)
	}
}

// FixRandomFinger refreshes this node's finger table entries in relation to Chord ring changes.
//
// Recommended to be called periodically in order to ensure finger table integrity.
func (node *LocalNode) FixRandomFinger() {
	i := rand.Int() % len(node.fingers)
	finger := node.fingers[i]
	finger.node = node.FindSuccessor(finger.Start())
}

// FixAllFingers refreshes all of this node's finger table entries.
func (node *LocalNode) FixAllFingers() {
	for _, finger := range node.fingers {
		finger.node = node.FindSuccessor(finger.Start())
	}
}

// PrintRing outputs this node's ring to console.
func (node *LocalNode) PrintRing() {
	fmt.Printf("Node %v ring: %v", node.String(), node.String())
	successor := node.Successor()
	for !node.Eq(successor) {
		fmt.Printf(" => %v", successor.String())
		successor = successor.Successor()
	}
	fmt.Println()
}

// String produces canonical string representation of this node.
func (node *LocalNode) String() string {
	return node.id.String()
}
