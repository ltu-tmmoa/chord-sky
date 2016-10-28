package chord

import (
	"fmt"
	"net"
	"sync"
)

// LocalNode represents a potential member of a Chord ring.
type LocalNode struct {
	addr        net.TCPAddr
	id          ID
	fingerTable FingerTable
	predecessor Node
	mutex       sync.RWMutex
}

// NewLocalNode creates a new local node from given address, which ought to be
// the application's public-facing IP address.
func NewLocalNode(addr *net.TCPAddr) *LocalNode {
	return newLocalNode(addr, identity(addr, HashBitsMax))
}

func newLocalNode(addr *net.TCPAddr, id *ID) *LocalNode {
	return &LocalNode{
		addr:        *addr,
		id:          *id,
		fingerTable: newFingerTable(id),
	}
}

// ID returns node ID.
func (node *LocalNode) ID() *ID {
	return &node.id
}

// TCPAddr provides node network IP address.
func (node *LocalNode) TCPAddr() *net.TCPAddr {
	return &node.addr
}

// FingerStart resolves start ID of finger table entry i.
//
// The result is only defined for i in [1,M], where M is the amount of bits set
// at node ring creation.
func (node *LocalNode) FingerStart(i int) *ID {
	node.mutex.RLock()
	defer node.mutex.RUnlock()

	return node.fingerTable.FingerStart(i)
}

// FingerNode resolves Chord node at given finger table offset i.
//
// The result is only defined for i in [1,M], where M is the amount of bits set
// at node ring creation.
func (node *LocalNode) FingerNode(i int) <-chan Node {
	return node.getNode(func() Node {
		return node.fingerNode(i)
	})
}

func (node *LocalNode) fingerNode(i int) Node {
	node.mutex.RLock()
	defer node.mutex.RUnlock()

	return node.fingerTable.FingerNode(i)
}

func (node *LocalNode) getNode(f func() Node) <-chan Node {
	ch := make(chan Node, 1)
	ch <- f()
	return ch
}

// SetFingerNode attempts to set this node's ith finger to given node.
//
// The operation is only valid for i in [1,M], where M is the amount of
// bits set at node ring creation.
func (node *LocalNode) SetFingerNode(i int, fing Node) {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	node.fingerTable.SetFingerNode(i, fing)
}

func (node *LocalNode) setFingerNodeUnlocked(i int, fing Node) {
	node.fingerTable.SetFingerNode(i, fing)
}

// Successor yields the next node in this node's ring.
func (node *LocalNode) Successor() <-chan Node {
	return node.FingerNode(1)
}

// Successor yields the next node in this node's ring.
func (node *LocalNode) successor() Node {
	return node.fingerNode(1)
}

// Predecessor yields the previous node in this node's ring.
func (node *LocalNode) Predecessor() <-chan Node {
	return node.getNode(func() Node {
		node.mutex.RLock()
		defer node.mutex.RUnlock()

		if node.predecessor == nil {
			node.predecessor = <-node.FindPredecessor(node.ID())
		}
		return node.predecessor
	})
}

// FindSuccessor asks this node to find successor of given ID.
//
// See Chord paper figure 4.
func (node *LocalNode) FindSuccessor(id *ID) <-chan Node {
	return node.getNode(func() Node {
		node.mutex.RLock()
		defer node.mutex.RUnlock()

		node0 := <-node.FindPredecessor(id)
		if node0 == nil {
			return nil
		}
		return <-node0.Successor()
	})
}

// FindPredecessor asks node to find id's predecessor.
//
// See Chord paper figure 4.
func (node *LocalNode) FindPredecessor(id *ID) <-chan Node {
	return node.getNode(func() Node {
		node.mutex.RLock()
		defer node.mutex.RUnlock()

		var n0 Node
		n0 = node
		for {
			succ := <-n0.Successor()
			if succ == nil {
				return nil
			}
			if idIntervalContainsEI(n0.ID(), succ.ID(), id) {
				return n0
			}
			n0 = closestPrecedingFinger(n0, id)
		}
	})
}

// Returns closest finger preceding ID.
//
// See Chord paper figure 4.
func closestPrecedingFinger(n Node, id *ID) Node {
	for i := n.ID().Bits(); i > 0; i-- {
		f := <-n.FingerNode(i)
		if f == nil {
			return nil
		}
		if idIntervalContainsEE(n.ID(), id, f.ID()) {
			return f
		}
	}
	return n
}

// SetSuccessor attempts to set this node's successor to given node.
func (node *LocalNode) SetSuccessor(succ Node) {
	node.SetFingerNode(1, succ)
}

// SetPredecessor attempts to set this node's predecessor to given node.
func (node *LocalNode) SetPredecessor(pred Node) {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	node.predecessor = pred
}

// DisassociateNodesByID removes any references held to node with an ID
// equivalent to given.
func (node *LocalNode) DisassociateNodesByID(id *ID) {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	node.fingerTable.RemoveFingerNodesByID(id)
	// TODO: Remove from successor list?

	if node.predecessor.ID().Eq(id) {
		node.predecessor = nil
	}
}

// PrintRing outputs this node's ring to console.
func (node *LocalNode) PrintRing() {
	fmt.Printf("Node %v ring:\n", node.String())
	succ := node.successor()
	for succ != nil && !node.ID().Eq(succ.ID()) {
		fmt.Printf(" => %v\n", succ)
		succ = node.successor()
	}
	fmt.Println()
}

// String produces canonical string representation of this node.
func (node *LocalNode) String() string {
	return fmt.Sprintf("%s@%s", node.id.String(), node.addr.String())
}
