package chord

import (
	"fmt"
	"math/big"
	"math/rand"
	"net"
	"sync"
)

// LocalNode represents a potential member of a Chord ring.
type LocalNode struct {
	tcpAddr     net.TCPAddr
	id          ID
	fingerTable FingerTable
	predecessor Node
	mutex       sync.RWMutex
}

// NewLocalNode creates a new local node from given address, which ought to be
// the application's public-facing IP address.
func NewLocalNode(tcpAddr *net.TCPAddr) *LocalNode {
	return newLocalNode(tcpAddr, identity(tcpAddr, HashBitsMax))
}

func newLocalNode(tcpAddr *net.TCPAddr, id *ID) *LocalNode {
	return &LocalNode{
		tcpAddr:     *tcpAddr,
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
	return &node.tcpAddr
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

// Join makes this node join the ring of given other node.
//
// Returns error only if failing to resolve successor or predecessor. All other
// operations are carried out on a best-effort basis.
//
// If given node is nil, this node will form its own ring.
func (node *LocalNode) Join(node0 Node) {
	if node0 != nil {
		node.initFingerTable(node0)
		node.updateOthers()
		// TODO: Move keys in (predecessor,n] from successor
	} else {
		node.SetSuccessor(node)
		node.SetPredecessor(node)
	}
}

// Initializes finger table of local node; node0 is an arbitrary node already
// in the network.
//
// Panics if failing to resolve successor or predecessor. All other operations
// are carried out on a best-effort basis.
//
// See Chord paper figure 6.
func (node *LocalNode) initFingerTable(node0 Node) {
	// Add this node to node0 node's ring.
	{
		succ := <-node0.FindSuccessor(node.FingerStart(1))
		if succ == nil {
			panic("Failed to resolve successor node.")
		}
		pred := <-succ.Predecessor()
		if pred == nil {
			panic("Failed to resolve predecessor node.")
		}

		node.SetSuccessor(succ)
		node.SetPredecessor(pred)

		pred.SetSuccessor(node)
		succ.SetPredecessor(node)
	}
	// Update this node's finger table, on best-effort basis.
	{
		m := node.id.Bits()
		for i := 1; i < m; i++ {
			this := node.fingerNode(i)
			nextStart := node.FingerStart(i + 1)

			var n Node
			if idIntervalContainsIE(node.ID(), this.ID(), nextStart) {
				n = this
			} else {
				n = <-node0.FindSuccessor(nextStart)
				if n == nil {
					continue
				}
			}
			node.setFingerNodeUnlocked(i+1, n)
		}
	}
}

// Update all nodes whose finger tables should refer to this node.
//
// Operates on a best-effort basis.
//
// See Chord paper figure 6.
func (node *LocalNode) updateOthers() {
	m := node.ID().Bits()
	for i := 2; i <= m; i++ {
		var id *ID
		{
			subt := big.Int{}
			subt.SetInt64(2)
			subt.Exp(&subt, big.NewInt(int64(i-1)), nil)

			id = node.ID().Diff(NewID(&subt, m))
		}
		pred := <-node.FindPredecessor(id)
		if pred == nil {
			continue
		}
		node.updateFingerTable(pred, node, i)
	}
}

// If s is the i:th finger of node, update node's finger table with s.
//
// Operates on a best-effort basis.
//
// See Chord paper figure 6.
func (node *LocalNode) updateFingerTable(n, s Node, i int) {
	fingNode := <-n.FingerNode(i)
	if fingNode == nil {
		return
	}
	if idIntervalContainsIE(n.FingerStart(i), fingNode.ID(), s.ID()) {
		n.SetFingerNode(i, s)
		pred := <-n.Predecessor()
		if pred == nil {
			return
		}
		node.updateFingerTable(pred, s, i)
	}
}

// Stabilize attempts to fix any ring issues arising from joining or leaving
// Chord ring nodes.
//
// Panics if fails to resolve successor's predecessor.
//
// Recommended to be called periodically in order to ensure node data
// integrity.
func (node *LocalNode) Stabilize() error {
	succ := node.successor()

	x := <-succ.Predecessor()
	if x == nil {
		return fmt.Errorf("Node stabilization failed. Unable to resolve %s predecessor.", succ.String())
	}
	if idIntervalContainsEE(node.ID(), succ.ID(), x.ID()) {
		node.SetSuccessor(x)
	}
	succ = node.successor()
	node.notify(succ)
	return nil
}

func (node *LocalNode) notify(node0 Node) {
	pred := <-node0.Predecessor()
	if pred == nil || idIntervalContainsEE(pred.ID(), node0.ID(), node.ID()) {
		node0.SetPredecessor(node)
	}
}

// FixRandomFinger refreshes this node's finger table entries in relation to Chord ring changes.
//
// Recommended to be called periodically in order to ensure finger table integrity.
func (node *LocalNode) FixRandomFinger() error {
	return node.FixFinger((rand.Int() % node.ID().Bits()) + 1)
}

// FixFinger refreshes finger indicated by given index i.
func (node *LocalNode) FixFinger(i int) error {
	succ := <-node.FindSuccessor(node.FingerStart(i))
	if succ != nil {
		node.setFingerNodeUnlocked(i, succ)
		return nil
	}
	return fmt.Errorf("Finger %d fix failed. Unable to resolve its successor node.", i)
}

// FixAllFingers refreshes all of this node's finger table entries.
func (node *LocalNode) FixAllFingers() error {
	for i := range node.fingerTable {
		if err := node.FixFinger(i + 1); err != nil {
			return err
		}
	}
	return nil
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
	return fmt.Sprintf("%s@%s", node.id.String(), node.tcpAddr.String())
}
