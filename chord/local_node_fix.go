package chord

import "math/rand"

// Attempts to fix any ring issues arising from joining or leaving chord ring
// nodes.
//
// Should be called periodically in order to ensure node data integrity.
func (node *localNode) stabilize() error {
	succ := node.successor()

	x, err := (<-succ.Predecessor()).Unwrap()
	if err != nil {
		return err
	}
	if idIntervalContainsEE(node.ID(), succ.ID(), x.ID()) {
		<-node.SetSuccessor(x)
	}
	succ = node.successor()
	return node.notify(succ)
}

func (node *localNode) notify(node0 Node) error {
	pred, err := (<-node0.Predecessor()).Unwrap()
	if err != nil {
		return err
	}
	if pred == nil || idIntervalContainsEE(pred.ID(), node0.ID(), node.ID()) {
		<-node0.SetPredecessor(node)
	}
	return nil
}

func (node *localNode) fixRandomFinger() error {
	return node.fixFinger((rand.Int() % node.ID().Bits()) + 1)
}

func (node *localNode) fixFinger(i int) error {
	succ, err := (<-node.FindSuccessor(node.FingerStart(i))).Unwrap()
	if err != nil {
		return err
	}
	<-node.SetFingerNode(i, succ)
	return nil
}

func (node *localNode) fixAllFingers() error {
	for i := range node.ftable.fingers {
		if err := node.fixFinger(i + 1); err != nil {
			return err
		}
	}
	return nil
}
