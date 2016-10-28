package chord

import (
	"net"
	"sync"
)

// NodePool holds a single local node and a set of remote nodes, facilitating
// management of remote node lifetimes.
type NodePool struct {
	lnode *localNode
	nodes map[string]Node
	mutex sync.Mutex
}

// NewNodePool creates a new pool of a single local node with the given local
// interface address.
func NewNodePool(laddr *net.TCPAddr) *NodePool {
	lnode := newLocalNode(laddr)
	return &NodePool{
		lnode: lnode,
		nodes: map[string]Node{
			laddr.String(): lnode,
		},
	}
}

// GetOrCreateNode gets existing node with same address, or creates a new one
// if no such is already known.
func (pool *NodePool) GetOrCreateNode(addr *net.TCPAddr) Node {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	key := addr.String()
	if node, ok := pool.nodes[key]; ok {
		return node
	}
	node := newRemoteNode(addr, pool)
	pool.nodes[key] = node
	return node
}

// RemoveNode removes node from pool with an `addr` matching given, if present.
//
// Removed nodes are always disassociated with the local node.
//
// Attempts at removing the local node from the pool will be ignored.
func (pool *NodePool) RemoveNode(addr *net.TCPAddr) {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	key := addr.String()
	if node, ok := pool.nodes[key]; ok && node != pool.lnode {
		pool.lnode.disassociateNode(node)
		delete(pool.nodes, key)
	}
}

// Refresh causes a hearbeat message to be sent to each remote node in pool.
func (pool *NodePool) Refresh() error {
	defer func() {
		for _, node := range pool.nodes {
			rnode, ok := node.(*remoteNode)
			if ok {
				rnode.Heartbeat()
			}
		}
	}()

	if err := pool.lnode.Stabilize(); err != nil {
		return err
	}
	if err := pool.lnode.FixRandomFinger(); err != nil {
		return err
	}
	return nil
}
