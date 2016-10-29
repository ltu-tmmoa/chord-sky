package chord

import (
	"math/big"

	"github.com/ltu-tmmoa/chord-sky/data"
	"github.com/ltu-tmmoa/chord-sky/log"
)

// Join makes this node join the ring of given other node.
//
// If given node is nil, this node will form its own ring.
//
// See Chord paper figure 6.
func (node *localNode) join(node0 Node) {
	if node0 != nil {
		node.initfingerTable(node0)
		node.updateOthers()
		if err := node.downloadStorageOf(node.successor()); err != nil {
			log.Logger.Println("Failed to download storage.", err.Error())
		}
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
func (node *localNode) initfingerTable(node0 Node) error {
	// Add this node to node0 node's ring.
	{
		succ, err := node0.FindSuccessor(node.FingerStart(1))
		if err != nil {
			return err
		}
		pred, err := succ.Predecessor()
		if err != nil {
			return err
		}

		node.SetSuccessor(succ)
		node.SetPredecessor(pred)

		if err = pred.SetSuccessor(node); err != nil {
			return err
		}
		if err = succ.SetPredecessor(node); err != nil {
			return err
		}
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
				n, _ = node0.FindSuccessor(nextStart)
				if n == nil {
					continue
				}
			}
			node.SetFingerNode(i+1, n)
		}
	}
	return nil
}

// Update all nodes whose finger tables should refer to this node.
//
// Operates on a best-effort basis.
//
// See Chord paper figure 6.
func (node *localNode) updateOthers() {
	m := node.ID().Bits()
	for i := 2; i <= m; i++ {
		var id *data.ID
		{
			subt := big.Int{}
			subt.SetInt64(2)
			subt.Exp(&subt, big.NewInt(int64(i-1)), nil)

			id = node.ID().Diff(data.NewID(&subt, m))
		}
		pred, err := node.FindPredecessor(id)
		if err != nil {
			log.Logger.Println(err.Error())
		}
		if pred != nil {
			node.updatefingerTable(pred, node, i)
		}
	}
}

// If s is the i:th finger of node, update node's finger table with s.
//
// Operates on a best-effort basis.
//
// See Chord paper figure 6.
func (node *localNode) updatefingerTable(n, s Node, i int) {
	fingNode, err := n.FingerNode(i)
	if err != nil {
		log.Logger.Println(err.Error())
	}
	if fingNode == nil {
		return
	}
	if idIntervalContainsIE(n.FingerStart(i), fingNode.ID(), s.ID()) {
		if err = n.SetFingerNode(i, s); err != nil {
			return
		}
		pred, err := n.Predecessor()
		if err != nil {
			log.Logger.Println(err.Error())
		}
		if pred == nil {
			return
		}
		node.updatefingerTable(pred, s, i)
	}
}
