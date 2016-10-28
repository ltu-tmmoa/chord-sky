package chord

import "math/big"

// Join makes this node join the ring of given other node.
//
// Returns error only if failing to resolve successor or predecessor. All other
// operations are carried out on a best-effort basis.
//
// If given node is nil, this node will form its own ring.
func (node *localNode) Join(node0 Node) {
	if node0 != nil {
		node.initfingerTable(node0)
		node.updateOthers()
		// TODO: Move keys in (predecessor,n] from successor
	} else {
		<-node.SetSuccessor(node)
		<-node.SetPredecessor(node)
	}
}

// Initializes finger table of local node; node0 is an arbitrary node already
// in the network.
//
// Panics if failing to resolve successor or predecessor. All other operations
// are carried out on a best-effort basis.
//
// See Chord paper figure 6.
func (node *localNode) initfingerTable(node0 Node) {
	// Add this node to node0 node's ring.
	{
		succ := <-node0.FindSuccessor(node.FingerStart(1))
		if succ == nil {
			panic("Failed to resolve successor node.")
		}
		pred := <-succ.Predecessor()
		if pred == nil {
			panic("Failed to resolve predecessor node.")
		}

		<-node.SetSuccessor(succ)
		<-node.SetPredecessor(pred)

		<-pred.SetSuccessor(node)
		<-succ.SetPredecessor(node)
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
				n = <-node0.FindSuccessor(nextStart)
				if n == nil {
					continue
				}
			}
			<-node.SetfingerNode(i+1, n)
		}
	}
}

// Update all nodes whose finger tables should refer to this node.
//
// Operates on a best-effort basis.
//
// See Chord paper figure 6.
func (node *localNode) updateOthers() {
	m := node.ID().Bits()
	for i := 2; i <= m; i++ {
		var id *ID
		{
			subt := big.Int{}
			subt.SetInt64(2)
			subt.Exp(&subt, big.NewInt(int64(i-1)), nil)

			id = node.ID().Diff(NewID(&subt, m))
		}
		pred := <-node.FindPredecessor(id)
		if pred == nil {
			continue
		}
		node.updatefingerTable(pred, node, i)
	}
}

// If s is the i:th finger of node, update node's finger table with s.
//
// Operates on a best-effort basis.
//
// See Chord paper figure 6.
func (node *localNode) updatefingerTable(n, s Node, i int) {
	fingNode := <-n.FingerNode(i)
	if fingNode == nil {
		return
	}
	if idIntervalContainsIE(n.FingerStart(i), fingNode.ID(), s.ID()) {
		<-n.SetfingerNode(i, s)
		pred := <-n.Predecessor()
		if pred == nil {
			return
		}
		node.updatefingerTable(pred, s, i)
	}
}
