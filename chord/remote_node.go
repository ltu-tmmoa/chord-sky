package chord

import (
	"fmt"
	"net"
)

// RemoteNode represents some Chord node available remotely.
type RemoteNode struct {
	addr net.TCPAddr
	id   ID
	pool *NodePool
}

// NewRemoteNode creates a new remote node from given address and
// disconnection handler.
//
// The node will have its TCP connection initialized when it is first used. The
// disconnection handler is called when the `Disconnect()` method is called, or
// whenever any network error occurs. In either case, the remote node is to be
// treated as stale and should be cleaned up.
//
// The `onDisconnect` function may be passed on to other remote nodes by this
// node.
func NewRemoteNode(addr *net.TCPAddr, pool *NodePool) *RemoteNode {
	return &RemoteNode{
		addr: *addr,
		id:   *Identity(addr),
		pool: pool,
	}
}

// ID returns node ID.
func (node *RemoteNode) ID() *ID {
	return &node.id
}

// TCPAddr provides node network address.
func (node *RemoteNode) TCPAddr() *net.TCPAddr {
	return &node.addr
}

// FingerStart resolves start ID of finger table entry i.
//
// The result is only defined for i in [1,M], where M is the amount of bits set
// at node ring creation.
func (node *RemoteNode) FingerStart(i int) *ID {
	m := node.ID().Bits()
	verifyIndexOrPanic(m, i)
	return calcFingerStart(node.ID(), i-1)
}

// FingerNode resolves Chord node at given finger table offset i.
//
// The result is only defined for i in [1,M], where M is the amount of bits set
// at node ring creation.
func (node *RemoteNode) FingerNode(i int) <-chan Node {
	return node.httpGetNodef("fingers/%d", i)
}

// SetFingerNode attempts to set this node's ith finger to given node.
//
// The operation is only valid for i in [1,M], where M is the amount of
// bits set at node ring creation.
func (node *RemoteNode) SetFingerNode(i int, fing Node) {
	node.httpPut(fmt.Sprintf("fingers/%d", i), fing.TCPAddr().String())
}

// Successor yields the next node in this node's ring.
func (node *RemoteNode) Successor() <-chan Node {
	return node.httpGetNodef("successor")
}

// Predecessor yields the previous node in this node's ring.
func (node *RemoteNode) Predecessor() <-chan Node {
	return node.httpGetNodef("predecessor")
}

// FindSuccessor asks this node to find successor of given ID.
func (node *RemoteNode) FindSuccessor(id *ID) <-chan Node {
	return node.httpGetNodef("successors?id=%s", id.String())
}

// FindPredecessor asks this node to find a predecessor of given ID.
func (node *RemoteNode) FindPredecessor(id *ID) <-chan Node {
	return node.httpGetNodef("predecessors?id=%s", id.String())
}

// SetSuccessor attempts to set this node's successor to given node.
func (node *RemoteNode) SetSuccessor(succ Node) {
	node.httpPut("successor", succ.TCPAddr().String())
}

// SetPredecessor attempts to set this node's predecessor to given node.
func (node *RemoteNode) SetPredecessor(pred Node) {
	node.httpPut("predecessor", pred.TCPAddr().String())
}

// String turns Node into its canonical string representation.
func (node *RemoteNode) String() string {
	return fmt.Sprintf("%s@%s", node.ID().String(), node.TCPAddr().String())
}
