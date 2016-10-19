package chord

import "testing"

func TestNodeJoin2(t *testing.T) {
	nodes := prepareNodes(0, 1)

	nodes[0].Join(nil)
	nodes[1].Join(nodes[0])

	for _, node := range nodes {
		node.FixAllFingers()
	}
	{
		expectPredecessorID, expectFingerNodeID := prepareNodeFingerTests(t, nodes[0])
		expectPredecessorID(1)
		expectFingerNodeID(1, 1)
		expectFingerNodeID(2, 0)
		expectFingerNodeID(3, 0)
	}
	{
		expectPredecessorID, expectFingerNodeID := prepareNodeFingerTests(t, nodes[1])
		expectPredecessorID(0)
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

	for _, node := range nodes {
		node.FixAllFingers()
	}
	{
		expectPredecessorID, expectFingerNodeID := prepareNodeFingerTests(t, nodes[0])
		expectPredecessorID(3)
		expectFingerNodeID(1, 1)
		expectFingerNodeID(2, 3)
		expectFingerNodeID(3, 0)
	}
	{
		expectPredecessorID, expectFingerNodeID := prepareNodeFingerTests(t, nodes[1])
		expectPredecessorID(0)
		expectFingerNodeID(1, 3)
		expectFingerNodeID(2, 3)
		expectFingerNodeID(3, 0)
	}
	{
		expectPredecessorID, expectFingerNodeID := prepareNodeFingerTests(t, nodes[2])
		expectPredecessorID(1)
		expectFingerNodeID(1, 0)
		expectFingerNodeID(2, 0)
		expectFingerNodeID(3, 0)
	}
}

func TestNodeJoin4(t *testing.T) {
	nodes := prepareNodes(0, 1, 3, 6)

	nodes[0].Join(nil)
	nodes[1].Join(nodes[0])
	nodes[2].Join(nodes[1])
	nodes[3].Join(nodes[2])

	for _, node := range nodes {
		node.FixAllFingers()
	}
	{
		expectPredecessorID, expectFingerNodeID := prepareNodeFingerTests(t, nodes[0])
		expectPredecessorID(6)
		expectFingerNodeID(1, 1)
		expectFingerNodeID(2, 3)
		expectFingerNodeID(3, 6)
	}
	{
		expectPredecessorID, expectFingerNodeID := prepareNodeFingerTests(t, nodes[1])
		expectPredecessorID(0)
		expectFingerNodeID(1, 3)
		expectFingerNodeID(2, 3)
		expectFingerNodeID(3, 6)
	}
	{
		expectPredecessorID, expectFingerNodeID := prepareNodeFingerTests(t, nodes[2])
		expectPredecessorID(1)
		expectFingerNodeID(1, 6)
		expectFingerNodeID(2, 6)
		expectFingerNodeID(3, 0)
	}
	{
		expectPredecessorID, expectFingerNodeID := prepareNodeFingerTests(t, nodes[3])
		expectPredecessorID(3)
		expectFingerNodeID(1, 0)
		expectFingerNodeID(2, 0)
		expectFingerNodeID(3, 3)
	}
}

func TestNodeJoin8(t *testing.T) {
	nodes := prepareNodes(0, 1, 2, 3, 4, 5, 6, 7)

	nodes[0].Join(nil)
	nodes[1].Join(nodes[0])
	nodes[2].Join(nodes[1])
	nodes[3].Join(nodes[1])
	nodes[4].Join(nodes[1])
	nodes[5].Join(nodes[4])
	nodes[6].Join(nodes[3])
	nodes[7].Join(nodes[3])

	for _, node := range nodes {
		node.FixAllFingers()
	}
	{
		expectPredecessorID, expectFingerNodeID := prepareNodeFingerTests(t, nodes[0])
		expectPredecessorID(7)
		expectFingerNodeID(1, 1)
		expectFingerNodeID(2, 2)
		expectFingerNodeID(3, 4)
	}
	{
		expectPredecessorID, expectFingerNodeID := prepareNodeFingerTests(t, nodes[1])
		expectPredecessorID(0)
		expectFingerNodeID(1, 2)
		expectFingerNodeID(2, 3)
		expectFingerNodeID(3, 5)
	}
	{
		expectPredecessorID, expectFingerNodeID := prepareNodeFingerTests(t, nodes[2])
		expectPredecessorID(1)
		expectFingerNodeID(1, 3)
		expectFingerNodeID(2, 4)
		expectFingerNodeID(3, 6)
	}
	{
		expectPredecessorID, expectFingerNodeID := prepareNodeFingerTests(t, nodes[3])
		expectPredecessorID(2)
		expectFingerNodeID(1, 4)
		expectFingerNodeID(2, 5)
		expectFingerNodeID(3, 7)
	}
	{
		expectPredecessorID, expectFingerNodeID := prepareNodeFingerTests(t, nodes[4])
		expectPredecessorID(3)
		expectFingerNodeID(1, 5)
		expectFingerNodeID(2, 6)
		expectFingerNodeID(3, 0)
	}
	{
		expectPredecessorID, expectFingerNodeID := prepareNodeFingerTests(t, nodes[5])
		expectPredecessorID(4)
		expectFingerNodeID(1, 6)
		expectFingerNodeID(2, 7)
		expectFingerNodeID(3, 1)
	}
	{
		expectPredecessorID, expectFingerNodeID := prepareNodeFingerTests(t, nodes[6])
		expectPredecessorID(5)
		expectFingerNodeID(1, 7)
		expectFingerNodeID(2, 0)
		expectFingerNodeID(3, 2)
	}
	{
		expectPredecessorID, expectFingerNodeID := prepareNodeFingerTests(t, nodes[7])
		expectPredecessorID(6)
		expectFingerNodeID(1, 0)
		expectFingerNodeID(2, 1)
		expectFingerNodeID(3, 3)
	}
}

func prepareNodes(ids ...int64) []*LocalNode {
	nodes := make([]*LocalNode, len(ids))
	for i, s := range ids {
		nodes[i] = newLocalNode(fakeAddr(s), newHash64(s, M3))
	}
	return nodes
}

func prepareNodeFingerTests(t *testing.T, node *LocalNode) (func(int64), func(int, int64)) {
	expectPredecessorID := func(predecessorID int64) {
		if n := node.predecessor; !n.Eq(newHash64(predecessorID, M3)) {
			t.Errorf("{%v}.predecessor expected to be %v, was %v", node, predecessorID, n)
		}
	}
	expectFingerNodeID := func(finger int, nodeID int64) {
		if n, _ := node.Finger(finger).Node(); !n.Eq(newHash64(nodeID, M3)) {
			t.Errorf("{%v}.finger(%v).node expected to be %v, was %v.", node, finger, nodeID, n)
		}
	}
	return expectPredecessorID, expectFingerNodeID
}
