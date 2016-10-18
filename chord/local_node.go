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
		fingers[i] = newFinger(id, i + 1)
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
	return node.fingers[i - 1]
}

// Successor yields the next node in this node's ring.
func (node *LocalNode) Successor() (Node, error) {
	return node.finger(1).Node()
}

// Predecessor yields the previous node in this node's ring.
func (node *LocalNode) Predecessor() (Node, error) {
	return node.predecessor, nil
}

// FindSuccessor asks this node to find successor of given ID.
//
// See Chord paper figure 4.
func (node *LocalNode) FindSuccessor(id ID) (Node, error) {
	if node0, err := node.FindPredecessor(id); err != nil {
		return nil, err
	} else {
		return node0.Successor()
	}
}

// FindPredecessor asks node to find id's predecessor.
//
// See Chord paper figure 4.
func (node *LocalNode) FindPredecessor(id ID) (Node, error) {
	return findPredecessor(node, id)
}

// Asks node to find id's predecessor.
//
// See Chord paper figure 4.
func findPredecessor(n Node, id ID) (Node, error) {
	n0 := n
	for {
		successor, err := n0.Successor()
		if err != nil {
			return nil, err
		}
		if idIntervalContainsEI(n0, successor, id) {
			return n0, nil
		}
		n0, err = closestPrecedingFinger(n0, id)
		if err != nil {
			return nil, err
		}
	}
	return n0, nil
}

// Returns closest finger preceding ID.
//
// See Chord paper figure 4.
func closestPrecedingFinger(n Node, id ID) (Node, error) {
	for i := n.Bits(); i > 0; i-- {
		f, err := n.Finger(i).Node()
		if err != nil {
			return nil, err
		}
		if idIntervalContainsEE(n, id, f) {
			return f, nil
		}
	}
	return n, nil
}

// SetSuccessor attempts to set this node's successor to given node.
func (node *LocalNode) SetSuccessor(successor Node) error {
	node.finger(1).SetNode(successor)
	return nil
}

// SetPredecessor attempts to set this node's predecessor to given node.
func (node *LocalNode) SetPredecessor(predecessor Node) error {
	node.predecessor = predecessor
	return nil
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
			node.finger(i).SetNode(node)
		}
		node.predecessor = node
	}
}

// Update all nodes whose finger tables should refer to node.
//
// See Chord paper figure 6.
func updateOthers(n Node) error {
	m := n.Bits()
	for i := 2; i <= m; i++ {
		var id ID
		{
			subtrahend := big.Int{}
			subtrahend.SetInt64(2)
			subtrahend.Exp(&subtrahend, big.NewInt(int64(i - 1)), nil)
			id = n.Diff(newHash(subtrahend, m))
		}
		predecessor, err := n.FindPredecessor(id)
		if err != nil {
			return err
		}
		err = updateFingerTable(predecessor, n, i)
		if err != nil {
			return err
		}
	}
	return nil
}

// If s is the i:th finger of node, update node's finger table with s.
//
// See Chord paper figure 6.
func updateFingerTable(n, s Node, i int) error {
	finger := n.Finger(i)
	fingerNode, err := finger.Node()
	if err != nil {
		return err
	}
	if idIntervalContainsIE(finger.Start(), fingerNode, s) {
		finger.SetNode(s)
		predecessor, err := n.Predecessor()
		if err != nil {
			return err
		}
		updateFingerTable(predecessor, s, i)
	}
	return nil
}

// Initializes finger table of local node; node0 is an arbitrary node already in
// the network.
//
// See Chord paper figure 6.
func (node *LocalNode) initFingerTable(node0 Node) error {
	var err error
	// Add this node to node0 node's ring.
	{
		var successor Node
		var predecessor Node

		successor, err = node0.FindSuccessor(node.finger(1).Start())
		if err != nil {
			return err
		}
		predecessor, err = successor.Predecessor()
		if err != nil {
			return err
		}

		node.SetSuccessor(successor)
		node.SetPredecessor(predecessor)

		if err = predecessor.SetSuccessor(node); err != nil {
			return err
		}
		if err = successor.SetPredecessor(node); err != nil {
			return err
		}
	}
	// Update this node's finger table.
	{
		m := node.Bits()
		for i := 1; i < m; i++ {
			this := node.finger(i)
			next := node.finger(i + 1)
			thisNode, err := this.Node()
			if err != nil {
				return err
			}
			var nextNode Node
			if idIntervalContainsIE(node, thisNode, next.Start()) {
				nextNode, err = this.Node()
			} else {
				nextNode, err = node0.FindSuccessor(next.Start())
			}
			if err != nil {
				return err
			}
			next.SetNode(nextNode)
		}
	}
	return nil
}

// Stabilize attempts to fix any ring issues arising from joining or leaving Chord ring nodes.
//
// Recommended to be called periodically in order to ensure node data integrity.
func (node *LocalNode) Stabilize() error {
	successor, err := node.Successor()
	if err != nil {
		return err
	}
	x, err := successor.Predecessor()
	if err != nil {
		return err
	}
	if idIntervalContainsEE(node, successor, x) {
		if err = node.SetSuccessor(x); err != nil {
			return err
		}
	}
	successor, err = node.Successor()
	if err != nil {
		return err
	}
	return notify(successor, node)
}

func notify(n, n0 Node) error {
	predecessor, err := n.Predecessor()
	if err != nil {
		return nil
	}
	if predecessor == nil || idIntervalContainsEE(predecessor, n, n0) {
		return n.SetPredecessor(n0)
	}
	return nil
}

// FixRandomFinger refreshes this node's finger table entries in relation to Chord ring changes.
//
// Recommended to be called periodically in order to ensure finger table integrity.
func (node *LocalNode) FixRandomFinger() error {
	i := rand.Int() % len(node.fingers)
	finger := node.fingers[i]
	fingerNode, err := node.FindSuccessor(finger.Start())
	finger.SetNode(fingerNode)
	return err
}

// FixAllFingers refreshes all of this node's finger table entries.
func (node *LocalNode) FixAllFingers() error {
	for _, finger := range node.fingers {
		fingerNode, err := node.FindSuccessor(finger.Start())
		if err != nil {
			return err
		}
		finger.SetNode(fingerNode)
	}
	return nil
}

// PrintRing outputs this node's ring to console.
func (node *LocalNode) PrintRing() {
	fmt.Printf("Node %v ring: %v", node.String(), node.String())
	successor, err := node.Successor()
	if err != nil {
		fmt.Println(err)
		return
	}
	for !node.Eq(successor) {
		fmt.Printf(" => %v", successor.String())
		successor, err = successor.Successor()
		if err != nil {
			fmt.Printf(" (%v)", err)
			return
		}
	}
	fmt.Println()
}

// String produces canonical string representation of this node.
func (node *LocalNode) String() string {
	return node.id.String()
}
