package chord

import (
	"net"
	"sync"
)

// PublicNode is used to expose Chord Node operations via RPC.
type PublicNode struct {
	node  Node
	mutex *sync.RWMutex
}

// NewPublicNode creates a new public node, wrapping given regular node, synchronizing it using provided mutex.
func NewPublicNode(node Node, mutex *sync.RWMutex) *PublicNode {
	publicNode := new(PublicNode)
	publicNode.node = node
	publicNode.mutex = mutex
	return publicNode
}

// GetFingerNode fetches Chord node at given finger table offset i.
//
// The result is only defined for i in [1,M], where M is the amount of bits set
// at node ring creation.
func (node *PublicNode) GetFingerNode(i *int, reply *net.IPAddr) error {
	node.mutex.RLock()
	defer node.mutex.RUnlock()

	*reply = *node.node.Finger(*i).Node().IPAddr()
	return nil
}

// GetSuccessor yields public node successor.
func (node *PublicNode) GetSuccessor(void *int, reply *net.IPAddr) error {
	node.mutex.RLock()
	defer node.mutex.RUnlock()

	successor, err := node.node.Successor()
	if err != nil {
		return err
	}
	*reply = *successor.IPAddr()
	return nil
}

// GetPredecessor yields public node predecessor.
func (node *PublicNode) GetPredecessor(void *int, reply *net.IPAddr) error {
	node.mutex.RLock()
	defer node.mutex.RUnlock()

	predecessor, err := node.node.Predecessor()
	if err != nil {
		return err
	}
	*reply = *predecessor.IPAddr()
	return nil
}

// FindSuccessor asks this node to find successor of given ID.
func (node *PublicNode) FindSuccessor(id *Hash, reply *net.IPAddr) error {
	node.mutex.RLock()
	defer node.mutex.RUnlock()

	successor, err := node.node.FindSuccessor(id)
	if err != nil {
		return err
	}
	*reply = *successor.IPAddr()
	return nil
}

// FindPredecessor asks this node to find a predecessor of given ID.
func (node *PublicNode) FindPredecessor(id *Hash, reply *net.IPAddr) error {
	node.mutex.RLock()
	defer node.mutex.RUnlock()

	predecessor, err := node.node.FindPredecessor(id)
	if err != nil {
		return err
	}
	*reply = *predecessor.IPAddr()
	return nil
}

// SetSuccessor attempts to set this node's successor to given node.
func (node *PublicNode) SetSuccessor(successor *net.IPAddr, void *int) error {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	return node.node.SetSuccessor(NewRemoteNode(successor))
}

// SetPredecessor attempts to set this node's predecessor to given node.
func (node *PublicNode) SetPredecessor(predecessor *net.IPAddr, void *int) error {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	return node.node.SetPredecessor(NewRemoteNode(predecessor))
}
