package chord

import (
	"fmt"
	"math/rand"
)

// Stabilize attempts to fix any ring issues arising from joining or leaving
// Chord ring nodes.
//
// Recommended to be called periodically in order to ensure node data
// integrity.
func (node *localNode) Stabilize() error {
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

// FixRandomFinger refreshes this node's finger table entries in relation to Chord ring changes.
//
// Recommended to be called periodically in order to ensure finger table integrity.
func (node *localNode) FixRandomFinger() error {
	return node.FixFinger((rand.Int() % node.ID().Bits()) + 1)
}

// FixFinger refreshes finger indicated by given index i.
func (node *localNode) FixFinger(i int) error {
	succ := <-node.FindSuccessor(node.FingerStart(i))
	if succ != nil {
		<-node.SetfingerNode(i, succ)
		return nil
	}
	return fmt.Errorf("Finger %d fix failed. Unable to resolve its successor node.", i)
}

// FixAllFingers refreshes all of this node's finger table entries.
func (node *localNode) FixAllFingers() error {
	for i := range node.ftable.fingers {
		if err := node.FixFinger(i + 1); err != nil {
			return err
		}
	}
	return nil
}
