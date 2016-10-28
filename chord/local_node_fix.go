package chord

import (
	"fmt"
	"math/rand"
)

// Attempts to fix any ring issues arising from joining or leaving chord ring
// nodes.
//
// Should be called periodically in order to ensure node data integrity.
func (node *localNode) stabilize() error {
	succ := node.successor()

	x := <-succ.Predecessor()
	if x == nil {
		return fmt.Errorf("Node stabilization failed. Unable to resolve %s predecessor.", succ)
	}
	if idIntervalContainsEE(node.ID(), succ.ID(), x.ID()) {
		<-node.SetSuccessor(x)
	}
	succ = node.successor()
	node.notify(succ)
	return nil
}

func (node *localNode) notify(node0 Node) {
	pred := <-node0.Predecessor()
	if pred == nil || idIntervalContainsEE(pred.ID(), node0.ID(), node.ID()) {
		<-node0.SetPredecessor(node)
	}
}

func (node *localNode) fixRandomFinger() error {
	return node.fixFinger((rand.Int() % node.ID().Bits()) + 1)
}

func (node *localNode) fixFinger(i int) error {
	succ := <-node.FindSuccessor(node.FingerStart(i))
	if succ != nil {
		<-node.SetfingerNode(i, succ)
		return nil
	}
	return fmt.Errorf("Finger %d fix failed. Unable to resolve its successor node.", i)
}

func (node *localNode) fixAllFingers() error {
	for i := range node.ftable.fingers {
		if err := node.fixFinger(i + 1); err != nil {
			return err
		}
	}
	return nil
}
