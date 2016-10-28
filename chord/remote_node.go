package chord

import (
	"fmt"
	"net"
)

// Represents some Chord node available remotely.
type remoteNode struct {
	addr net.TCPAddr
	id   ID
	pool *NodePool
}

func newRemoteNode(addr *net.TCPAddr, pool *NodePool) *remoteNode {
	return &remoteNode{
		addr: *addr,
		id:   *hashAddr(addr),
		pool: pool,
	}
}

func (node *remoteNode) ID() *ID {
	return &node.id
}

func (node *remoteNode) TCPAddr() *net.TCPAddr {
	return &node.addr
}

func (node *remoteNode) FingerStart(i int) *ID {
	m := node.ID().Bits()
	verifyIndexOrPanic(m, i)
	return calcfingerStart(node.ID(), i-1)
}

func (node *remoteNode) FingerNode(i int) <-chan Node {
	return node.httpGetNodef("fingers/%d", i)
}

func (node *remoteNode) SetfingerNode(i int, fing Node) <-chan *struct{} {
	return node.httpPut(fmt.Sprintf("fingers/%d", i), fing.TCPAddr().String())
}

func (node *remoteNode) Heartbeat() {
	node.httpHeartbeat("heartbeat")
}

func (node *remoteNode) Successor() <-chan Node {
	return node.httpGetNodef("successor")
}

func (node *remoteNode) Predecessor() <-chan Node {
	return node.httpGetNodef("predecessor")
}

func (node *remoteNode) FindSuccessor(id *ID) <-chan Node {
	return node.httpGetNodef("successors?id=%s", id.String())
}

func (node *remoteNode) FindPredecessor(id *ID) <-chan Node {
	return node.httpGetNodef("predecessors?id=%s", id.String())
}

func (node *remoteNode) SetSuccessor(succ Node) <-chan *struct{} {
	return node.httpPut("successor", succ.TCPAddr().String())
}

func (node *remoteNode) SetPredecessor(pred Node) <-chan *struct{} {
	return node.httpPut("predecessor", pred.TCPAddr().String())
}

func (node *remoteNode) String() string {
	return fmt.Sprintf("%s@%s", node.ID().String(), node.TCPAddr().String())
}
