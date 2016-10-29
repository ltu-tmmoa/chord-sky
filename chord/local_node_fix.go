package chord

import (
	"math/rand"

	"github.com/ltu-tmmoa/chord-sky/data"
)

// Attempts to fix any ring issues arising from joining or leaving chord ring
// nodes.
//
// Should be called periodically in order to ensure node data integrity.
func (node *localNode) stabilize() error {
	succ := node.successor()

	x, err := succ.Predecessor()
	if err != nil {
		return err
	}
	if data.IDIntervalContainsEE(node.ID(), succ.ID(), x.ID()) {
		node.SetSuccessor(succ)
	}
	succ = node.successor()
	return node.notify(succ)
}

func (node *localNode) notify(node0 Node) error {
	pred, err := node0.Predecessor()
	if err != nil {
		return err
	}
	if pred == nil || data.IDIntervalContainsEE(pred.ID(), node0.ID(), node.ID()) {
		if err = node0.SetPredecessor(node); err != nil {
			return err
		}
	}
	return nil
}

func (node *localNode) fixSuccessorList() error {
	succ, _ := node.Successor()
	succs := []Node{succ}

	var err error
	for i := 1; i < 3; i++ {
		succ, err = succ.Successor()
		if err != nil {
			return err
		}
		succs = append(succs, succ)
	}
	node.setSuccessorList(succs)
	return nil
}

func (node *localNode) fixRandomFinger() error {
	return node.fixFinger((rand.Int() % node.ID().Bits()) + 1)
}

func (node *localNode) fixFinger(i int) error {
	succ, err := node.FindSuccessor(node.FingerStart(i))
	if err != nil {
		return err
	}
	return node.SetFingerNode(i, succ)
}

func (node *localNode) fixAllFingers() error {
	for i := range node.ftable.fingers {
		if err := node.fixFinger(i + 1); err != nil {
			return err
		}
	}
	return nil
}
