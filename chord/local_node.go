package chord

import (
	"fmt"
	"io"
	"net"
)

// LocalNode represents a potential member of a Chord ring.
type LocalNode struct {
	addr        net.TCPAddr
	id          ID
	fingerTable *FingerTable
	predecessor Node
}

// NewLocalNode creates a new local node from given address, which ought to be
// the application's public-facing IP address.
func NewLocalNode(addr *net.TCPAddr) *LocalNode {
	return newLocalNode(addr, identity(addr, HashBitsMax))
}

func newLocalNode(addr *net.TCPAddr, id *ID) *LocalNode {
	node := &LocalNode{
		addr: *addr,
		id:   *id,
	}
	node.fingerTable = newFingerTable(node)
	return node
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
func (node *LocalNode) SetFingerNode(i int, fing Node) <-chan *struct{} {
	return node.getVoid(func() {
		node.fingerTable.SetFingerNode(i, fing)
	})
}

func (node *LocalNode) getVoid(f func()) <-chan *struct{} {
	ch := make(chan *struct{}, 1)
	f()
	ch <- nil
	return ch
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
		if node.predecessor == nil {
			node.predecessor = <-node.FindPredecessor(node.ID())
			if node.predecessor == nil {
				node.predecessor = node
			}
		}
		return node.predecessor
	})
}

// FindSuccessor asks this node to find successor of given ID.
//
// See Chord paper figure 4.
func (node *LocalNode) FindSuccessor(id *ID) <-chan Node {
	return node.getNode(func() Node {
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
func (node *LocalNode) SetSuccessor(succ Node) <-chan *struct{} {
	return node.getVoid(func() {
		node.SetFingerNode(1, succ)
	})
}

// SetPredecessor attempts to set this node's predecessor to given node.
func (node *LocalNode) SetPredecessor(pred Node) <-chan *struct{} {
	return node.getVoid(func() {
		node.predecessor = pred
	})
}

// DisassociateNodeByID removes any references held to node with an ID
// equivalent to given.
func (node *LocalNode) DisassociateNodeByID(id *ID) {
	node.fingerTable.RemoveFingerNodesByID(id)
	// TODO: Remove from successor list?

	if node.predecessor != nil && node.predecessor.ID().Eq(id) {
		node.predecessor = nil
	}
}

// WriteRingTextTo writes a list of the members of this node's ring to `w`.
//
// It might take a while before this returns, as it might need to call a lot of
// remote hosts to gather all required data.
func (node *LocalNode) WriteRingTextTo(w io.Writer) {
	succ := node.successor()
	for succ != nil {
		fmt.Fprintf(w, "%v\r\n", succ)
		if node.ID().Eq(succ.ID()) {
			break
		}
		succ = <-succ.Successor()
	}
}

// String produces canonical string representation of this node.
func (node *LocalNode) String() string {
	return fmt.Sprintf("%s@%s", node.id.String(), node.addr.String())
}
