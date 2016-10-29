package chord

import (
	"fmt"
	"net"

	"github.com/ltu-tmmoa/chord-sky/data"
)

// Represents some Chord node available remotely.
type remoteNode struct {
	addr    net.TCPAddr
	id      data.ID
	pool    *nodePool
	storage data.Storage
}

func newRemoteNode(addr *net.TCPAddr, pool *nodePool) *remoteNode {
	return &remoteNode{
		addr: *addr,
		id:   *addrToID(addr),
		pool: pool,
	}
}

func (node *remoteNode) ID() *data.ID {
	return &node.id
}

func (node *remoteNode) TCPAddr() *net.TCPAddr {
	return &node.addr
}

func (node *remoteNode) FingerStart(i int) *data.ID {
	m := node.ID().Bits()
	verifyIndexOrPanic(m, i)
	return calcfingerStart(node.ID(), i-1)
}

func (node *remoteNode) FingerNode(i int) (Node, error) {
	return node.httpGetNodef("fingers/%d", i)
}

func (node *remoteNode) SetFingerNode(i int, fing Node) error {
	return node.httpPut(fmt.Sprintf("fingers/%d", i), fing.TCPAddr().String())
}

func (node *remoteNode) Heartbeat() {
	node.httpHeartbeat("heartbeat")
}

func (node *remoteNode) Successor() (Node, error) {
	return node.httpGetNodef("successor")
}

func (node *remoteNode) SuccessorList() ([]Node, error) {
	return node.httpGetNodesf("successors")
}

func (node *remoteNode) Predecessor() (Node, error) {
	return node.httpGetNodef("predecessor")
}

func (node *remoteNode) FindSuccessor(id *data.ID) (Node, error) {
	return node.httpGetNodef("successors?id=%s", id.String())
}

func (node *remoteNode) FindPredecessor(id *data.ID) (Node, error) {
	return node.httpGetNodef("predecessors?id=%s", id.String())
}

func (node *remoteNode) SetSuccessor(succ Node) error {
	return node.httpPut("successor", succ.TCPAddr().String())
}

func (node *remoteNode) SetPredecessor(pred Node) error {
	return node.httpPut("predecessor", pred.TCPAddr().String())
}

func (node *remoteNode) Storage() data.Storage {
	return node.storage
}

func (node *remoteNode) String() string {
	return fmt.Sprintf("%s@%s", node.ID().String(), node.TCPAddr().String())
}
