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
	succlist    []Node
	predecessor Node
}

// NewLocalNode creates a new local node from given address, which ought to be
// the application's public-facing IP address.
func newLocalNode(addr *net.TCPAddr) *localNode {
	return newLocalNodeID(addr, hashAddr(addr))
}

func newLocalNodeID(addr *net.TCPAddr, id *ID) *localNode {
	node := &localNode{
		addr: *addr,
		id:   *id,
	}
	node.ftable = newFingerTable(node)
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

func (node *localNode) FingerNode(i int) (Node, error) {
	return node.fingerNode(i), nil
}

func (node *localNode) fingerNode(i int) Node {
	return node.ftable.fingerNode(i)
}

func (node *localNode) SetFingerNode(i int, fing Node) error {
	node.ftable.setFingerNode(i, fing)
	return nil
}

func (node *localNode) Successor() (Node, error) {
	return node.FingerNode(1)
}

func (node *localNode) SuccessorList() ([]Node, error) {
	return node.succlist, nil
}

func (node *localNode) successor() Node {
	return node.fingerNode(1)
}

func (node *localNode) Predecessor() (Node, error) {
	if node.predecessor == nil {
		pred, err := node.FindPredecessor(node.ID())
		if err != nil {
			return nil, err
		}
		if node.predecessor == nil {
			node.predecessor = pred
		}
	}
	return node.predecessor, nil
}

func (node *localNode) FindSuccessor(id *ID) (Node, error) {
	pred, err := node.FindPredecessor(id)
	if err != nil {
		return nil, err
	}
	succ, err := pred.Successor()
	if err != nil {
		return nil, err
	}
	return succ, nil
}

func (node *localNode) FindPredecessor(id *ID) (Node, error) {
	var n0 Node
	n0 = node
	for {
		succ, err := n0.Successor()
		if err != nil {
			return nil, err
		}
		if idIntervalContainsEI(n0.ID(), succ.ID(), id) {
			return n0, nil
		}
		n0, err = closestPrecedingFinger(n0, id)
		if err != nil {
			return nil, err
		}
	}
}

// Returns closest finger preceding ID.
//
// See Chord paper figure 4.
func closestPrecedingFinger(n Node, id *ID) (Node, error) {
	for i := n.ID().Bits(); i > 0; i-- {
		f, err := n.FingerNode(i)
		if err != nil {
			return nil, err
		}
		if idIntervalContainsEE(n.ID(), id, f.ID()) {
			return f, nil
		}
	}
	return n, nil
}

func (node *localNode) SetSuccessor(succ Node) error {
	node.ftable.setFingerNode(1, succ)
	node.succlist = []Node{succ}
	return nil
}

func (node *localNode) setSuccessorList(succs []Node) error {
	nsuccs := []Node{}
	for _, succ := range succs {
		id := succ.ID()
		hits := 0
		for _, osucc := range node.succlist {
			oid := osucc.ID()
			if id.Eq(oid) {
				hits++
			}
		}
		if hits == 0 {
			nsuccs = append(nsuccs, succ)
		}
	}
	for _, nsucc := range nsuccs {
		fmt.Println("Sending data to", nsucc)
		go func() {
			// TODO: Backup data with nsucc.
		}()
	}
	node.succlist = succs
	return nil
}

func (node *localNode) SetPredecessor(pred Node) error {
	node.predecessor = pred
	return nil
}

func (node *localNode) disassociateNode(n Node) {
	id := n.ID()
	node.ftable.removeFingerNodesByID(id)
	succlist := node.succlist
	for i, succ := range succlist {
		if succ.ID().Eq(id) {
			succlist = append(node.succlist[:i], node.succlist[i+1:]...)
		}
	}
	if len(succlist) > 0 {
		node.ftable.setFingerNode(1, succlist[0])
	}
	node.succlist = succlist
	if node.predecessor != nil && node.predecessor.ID().Eq(id) {
		node.predecessor = nil
	}
}

// Writes a list of the members of this node's ring to `w`.
//
// It might take a while before this returns, as it might need to call a lot of
// remote hosts to gather all required data.
func (node *localNode) writeRingTextTo(w io.Writer) {
	var err error

	succ := node.successor()
	for succ != nil {
		fmt.Fprintf(w, "%v\r\n", succ)
		if node.ID().Eq(succ.ID()) {
			break
		}
		succ, err = succ.Successor()
		if err != nil {
			fmt.Fprintf(w, "%v\r\n", err.Error())
		}
	}
}

// String produces canonical string representation of this node.
func (node *localNode) String() string {
	return fmt.Sprintf("%s@%s", node.id.String(), node.addr.String())
}
