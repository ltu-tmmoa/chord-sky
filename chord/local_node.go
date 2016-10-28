package chord

import (
	"fmt"
	"io"
	"net"
)

// localNode represents a potential member of a Chord ring.
type localNode struct {
	addr        net.TCPAddr
	id          ID
	ftable      *fingerTable
	predecessor Node
}

// NewLocalNode creates a new local node from given address, which ought to be
// the application's public-facing IP address.
func newLocalNode(addr *net.TCPAddr) *localNode {
	return newLocalNodeID(addr, identity(addr, HashBitsMax))
}

func newLocalNodeID(addr *net.TCPAddr, id *ID) *localNode {
	node := &localNode{
		addr: *addr,
		id:   *id,
	}
	node.ftable = newfingerTable(node)
	return node
}

func (node *localNode) ID() *ID {
	return &node.id
}

func (node *localNode) TCPAddr() *net.TCPAddr {
	return &node.addr
}

func (node *localNode) FingerStart(i int) *ID {
	return node.ftable.fingerStart(i)
}

func (node *localNode) FingerNode(i int) <-chan Node {
	return node.getNode(func() Node {
		return node.fingerNode(i)
	})
}

func (node *localNode) fingerNode(i int) Node {
	return node.ftable.fingerNode(i)
}

func (node *localNode) SetfingerNode(i int, fing Node) <-chan *struct{} {
	return node.getVoid(func() {
		node.ftable.setfingerNode(i, fing)
	})
}

// Successor yields the next node in this node's ring.
func (node *localNode) Successor() <-chan Node {
	return node.FingerNode(1)
}

// Successor yields the next node in this node's ring.
func (node *localNode) successor() Node {
	return node.fingerNode(1)
}

func (node *localNode) Predecessor() <-chan Node {
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

func (node *localNode) FindSuccessor(id *ID) <-chan Node {
	return node.getNode(func() Node {
		node0 := <-node.FindPredecessor(id)
		if node0 == nil {
			return nil
		}
		return <-node0.Successor()
	})
}

func (node *localNode) FindPredecessor(id *ID) <-chan Node {
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

func (node *localNode) SetSuccessor(succ Node) <-chan *struct{} {
	return node.getVoid(func() {
		node.SetfingerNode(1, succ)
	})
}

func (node *localNode) SetPredecessor(pred Node) <-chan *struct{} {
	return node.getVoid(func() {
		node.predecessor = pred
	})
}

func (node *localNode) disassociateNode(n Node) {
	id := n.ID()
	node.ftable.removefingerNodesByID(id)
	// TODO: Remove from successor list?

	if node.predecessor != nil && node.predecessor.ID().Eq(id) {
		node.predecessor = nil
	}
}

// Writes a list of the members of this node's ring to `w`.
//
// It might take a while before this returns, as it might need to call a lot of
// remote hosts to gather all required data.
func (node *localNode) writeRingTextTo(w io.Writer) {
	succ := node.successor()
	for succ != nil {
		fmt.Fprintf(w, "%v\r\n", succ)
		if node.ID().Eq(succ.ID()) {
			break
		}
		succ = <-succ.Successor()
	}
}

func (node *localNode) getNode(f func() Node) <-chan Node {
	ch := make(chan Node, 1)
	ch <- f()
	return ch
}

func (node *localNode) getVoid(f func()) <-chan *struct{} {
	ch := make(chan *struct{}, 1)
	f()
	ch <- nil
	return ch
}

// String produces canonical string representation of this node.
func (node *localNode) String() string {
	return fmt.Sprintf("%s@%s", node.id.String(), node.addr.String())
}
