package chord

import (
	"net"
	"sync"
)

// Holds a single local node and a set of remote nodes, allowing management of
// remote node lifetimes.
type nodePool struct {
	lnode *localNode
	nodes map[string]Node
	mutex sync.Mutex
}

func newNodePool(laddr *net.TCPAddr) *nodePool {
	lnode := newLocalNode(laddr)
	return &nodePool{
		lnode: lnode,
		nodes: map[string]Node{
			laddr.String(): lnode,
		},
	}
}

func (pool *nodePool) getOrCreateNode(addr *net.TCPAddr) Node {
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

func (pool *nodePool) removeNode(node Node) {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	key := node.TCPAddr().String()
	if node, ok := pool.nodes[key]; ok && node != pool.lnode {
		pool.lnode.disassociateNode(node)
		delete(pool.nodes, key)
	}
}

func (pool *nodePool) refresh() error {
	defer func() {
		for _, node := range pool.nodes {
			rnode, ok := node.(*remoteNode)
			if ok {
				rnode.Heartbeat()
			}
		}
	}()

	if err := pool.lnode.stabilize(); err != nil {
		return err
	}
	if err := pool.lnode.fixRandomFinger(); err != nil {
		return err
	}
	return nil
}
