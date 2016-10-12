package chord

import (
	"fmt"
	"testing"
)

func TestNodeJoin2(t *testing.T) {
	nodes := prepareNodes(0, 1)

	nodes[0].Join(nil)
	nodes[1].Join(nodes[0])

	{
		expectFingerNodeID := prepareNodeFingerTester(t, nodes[0])
		expectFingerNodeID(1, 1)
		expectFingerNodeID(2, 0)
		expectFingerNodeID(3, 0)
	}
	{
		expectFingerNodeID := prepareNodeFingerTester(t, nodes[1])
		expectFingerNodeID(1, 0)
		expectFingerNodeID(2, 0)
		expectFingerNodeID(3, 0)
	}
}

func TestNodeJoin3(t *testing.T) {
	nodes := prepareNodes(0, 1, 3)

	nodes[0].Join(nil)
	nodes[1].Join(nodes[0])
	nodes[2].Join(nodes[1])

	{
		expectFingerNodeID := prepareNodeFingerTester(t, nodes[0])
		expectFingerNodeID(1, 1)
		expectFingerNodeID(2, 3)
		expectFingerNodeID(3, 0)
	}
	{
		expectFingerNodeID := prepareNodeFingerTester(t, nodes[1])
		expectFingerNodeID(1, 3)
		expectFingerNodeID(2, 3)
		expectFingerNodeID(3, 0)
	}
	{
		expectFingerNodeID := prepareNodeFingerTester(t, nodes[2])
		expectFingerNodeID(1, 0)
		expectFingerNodeID(2, 0)
		expectFingerNodeID(3, 0)
	}
}

func TestNodeJoin4(t *testing.T) {
	nodes := prepareNodes(0, 1, 3, 6)

	fmt.Print("\nTEST 4 START\n\n")

	nodes[0].Join(nil)
	nodes[1].Join(nodes[0])
	nodes[2].Join(nodes[1])
	nodes[3].Join(nodes[2])

	fmt.Print("\nTEST 4 STOP\n\n")

	{
		expectFingerNodeID := prepareNodeFingerTester(t, nodes[0])
		expectFingerNodeID(1, 1)
		expectFingerNodeID(2, 3)
		expectFingerNodeID(3, 6)
	}
	{
		expectFingerNodeID := prepareNodeFingerTester(t, nodes[1])
		expectFingerNodeID(1, 3)
		expectFingerNodeID(2, 3)
		expectFingerNodeID(3, 6)
	}
	{
		expectFingerNodeID := prepareNodeFingerTester(t, nodes[2])
		expectFingerNodeID(1, 6)
		expectFingerNodeID(2, 6)
		expectFingerNodeID(3, 0)
	}
	{
		expectFingerNodeID := prepareNodeFingerTester(t, nodes[3])
		expectFingerNodeID(1, 0)
		expectFingerNodeID(2, 0)
		expectFingerNodeID(3, 3)
	}
	printNodes(nodes)
}

/*
func TestNodeJoin(t *testing.T) {
	nodes := prepareNodes(0, 1, 2, 3, 4, 5, 6, 7)

	nodes[0].Join(nil)
	nodes[1].Join(nodes[0])
	nodes[2].Join(nodes[1])
	nodes[3].Join(nodes[1])
	nodes[4].Join(nodes[1])
	nodes[5].Join(nodes[4])
	nodes[6].Join(nodes[3])
	nodes[7].Join(nodes[3])

	printNodes(nodes)

	//nodes[0].PrintRing()
}*/

func prepareNodeFingerTester(t *testing.T, node *Node) func(int, int64) {
	return func(finger int, nodeID int64) {
		if n := node.Finger(finger).Node(); !n.Eq(newHash64(nodeID, M3)) {
			t.Errorf("{%v}.finger(%v).node expected to be %v, was %v.", node, finger, nodeID, n)
		}
	}
}

func prepareNodes(ids ...int64) []*Node {
	nodes := make([]*Node, len(ids))
	for i, s := range ids {
		nodes[i] = newNode(stringAddr(fmt.Sprintf("%02d", s)), newHash64(s, M3))
	}
	return nodes
}

func printNodes(nodes []*Node) {
	for _, node := range nodes {
		fmt.Printf("Node: %v\n", node)
		for i := 1; i <= M3; i++ {
			finger := node.finger(i)
			fmt.Printf("  Finger: %v, %v, %v\n", i, finger.Interval(), finger.Node())
		}
		fmt.Printf("  Predecessor: %v\n", node.predecessor)
	}
}
