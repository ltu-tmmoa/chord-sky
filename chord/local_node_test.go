package chord

import "testing"

func TestNodeJoin2(t *testing.T) {
	nodes := prepareNodes(0, 1)

	nodes[0].join(nil)
	nodes[1].join(nodes[0])

	for _, node := range nodes {
		node.fixSuccessorList()
		node.fixAllFingers()
	}
	{
		expectPredecessorID, expectSuccessorIDs, expectFingerNodeID := prepareNodeFingerTests(t, nodes[0])
		expectPredecessorID(1)
		expectSuccessorIDs(1, 0, 1)
		expectFingerNodeID(1, 1)
		expectFingerNodeID(2, 0)
		expectFingerNodeID(3, 0)
	}
	{
		expectPredecessorID, expectSuccessorIDs, expectFingerNodeID := prepareNodeFingerTests(t, nodes[1])
		expectPredecessorID(0)
		expectSuccessorIDs(0, 1, 0)
		expectFingerNodeID(1, 0)
		expectFingerNodeID(2, 0)
		expectFingerNodeID(3, 0)
	}
}

func TestNodeJoin3(t *testing.T) {
	nodes := prepareNodes(0, 1, 3)

	nodes[0].join(nil)
	nodes[1].join(nodes[0])
	nodes[2].join(nodes[1])

	for _, node := range nodes {
		node.fixSuccessorList()
		node.fixAllFingers()
	}
	{
		expectPredecessorID, expectSuccessorIDs, expectFingerNodeID := prepareNodeFingerTests(t, nodes[0])
		expectPredecessorID(3)
		expectSuccessorIDs(1, 3, 0)
		expectFingerNodeID(1, 1)
		expectFingerNodeID(2, 3)
		expectFingerNodeID(3, 0)
	}
	{
		expectPredecessorID, expectSuccessorIDs, expectFingerNodeID := prepareNodeFingerTests(t, nodes[1])
		expectPredecessorID(0)
		expectSuccessorIDs(3, 0, 1)
		expectFingerNodeID(1, 3)
		expectFingerNodeID(2, 3)
		expectFingerNodeID(3, 0)
	}
	{
		expectPredecessorID, expectSuccessorIDs, expectFingerNodeID := prepareNodeFingerTests(t, nodes[2])
		expectPredecessorID(1)
		expectSuccessorIDs(0, 1, 3)
		expectFingerNodeID(1, 0)
		expectFingerNodeID(2, 0)
		expectFingerNodeID(3, 0)
	}
}

func TestNodeJoin4(t *testing.T) {
	nodes := prepareNodes(0, 1, 3, 6)

	nodes[0].join(nil)
	nodes[1].join(nodes[0])
	nodes[2].join(nodes[1])
	nodes[3].join(nodes[2])

	for _, node := range nodes {
		node.fixSuccessorList()
		node.fixAllFingers()
	}
	{
		expectPredecessorID, expectSuccessorIDs, expectFingerNodeID := prepareNodeFingerTests(t, nodes[0])
		expectPredecessorID(6)
		expectSuccessorIDs(1, 3, 6)
		expectFingerNodeID(1, 1)
		expectFingerNodeID(2, 3)
		expectFingerNodeID(3, 6)
	}
	{
		expectPredecessorID, expectSuccessorIDs, expectFingerNodeID := prepareNodeFingerTests(t, nodes[1])
		expectPredecessorID(0)
		expectSuccessorIDs(3, 6, 0)
		expectFingerNodeID(1, 3)
		expectFingerNodeID(2, 3)
		expectFingerNodeID(3, 6)
	}
	{
		expectPredecessorID, expectSuccessorIDs, expectFingerNodeID := prepareNodeFingerTests(t, nodes[2])
		expectPredecessorID(1)
		expectSuccessorIDs(6, 0, 1)
		expectFingerNodeID(1, 6)
		expectFingerNodeID(2, 6)
		expectFingerNodeID(3, 0)
	}
	{
		expectPredecessorID, expectSuccessorIDs, expectFingerNodeID := prepareNodeFingerTests(t, nodes[3])
		expectPredecessorID(3)
		expectSuccessorIDs(0, 1, 3)
		expectFingerNodeID(1, 0)
		expectFingerNodeID(2, 0)
		expectFingerNodeID(3, 3)
	}
}

func TestNodeJoin8(t *testing.T) {
	nodes := prepareNodes(0, 1, 2, 3, 4, 5, 6, 7)

	nodes[0].join(nil)
	nodes[1].join(nodes[0])
	nodes[2].join(nodes[1])
	nodes[3].join(nodes[1])
	nodes[4].join(nodes[1])
	nodes[5].join(nodes[4])
	nodes[6].join(nodes[3])
	nodes[7].join(nodes[3])

	for _, node := range nodes {
		node.fixSuccessorList()
		node.fixAllFingers()
	}
	{
		expectPredecessorID, expectSuccessorIDs, expectFingerNodeID := prepareNodeFingerTests(t, nodes[0])
		expectPredecessorID(7)
		expectSuccessorIDs(1, 2, 3)
		expectFingerNodeID(1, 1)
		expectFingerNodeID(2, 2)
		expectFingerNodeID(3, 4)
	}
	{
		expectPredecessorID, expectSuccessorIDs, expectFingerNodeID := prepareNodeFingerTests(t, nodes[1])
		expectPredecessorID(0)
		expectSuccessorIDs(2, 3, 4)
		expectFingerNodeID(1, 2)
		expectFingerNodeID(2, 3)
		expectFingerNodeID(3, 5)
	}
	{
		expectPredecessorID, expectSuccessorIDs, expectFingerNodeID := prepareNodeFingerTests(t, nodes[2])
		expectPredecessorID(1)
		expectSuccessorIDs(3, 4, 5)
		expectFingerNodeID(1, 3)
		expectFingerNodeID(2, 4)
		expectFingerNodeID(3, 6)
	}
	{
		expectPredecessorID, expectSuccessorIDs, expectFingerNodeID := prepareNodeFingerTests(t, nodes[3])
		expectPredecessorID(2)
		expectSuccessorIDs(4, 5, 6)
		expectFingerNodeID(1, 4)
		expectFingerNodeID(2, 5)
		expectFingerNodeID(3, 7)
	}
	{
		expectPredecessorID, expectSuccessorIDs, expectFingerNodeID := prepareNodeFingerTests(t, nodes[4])
		expectPredecessorID(3)
		expectSuccessorIDs(5, 6, 7)
		expectFingerNodeID(1, 5)
		expectFingerNodeID(2, 6)
		expectFingerNodeID(3, 0)
	}
	{
		expectPredecessorID, expectSuccessorIDs, expectFingerNodeID := prepareNodeFingerTests(t, nodes[5])
		expectPredecessorID(4)
		expectSuccessorIDs(6, 7, 0)
		expectFingerNodeID(1, 6)
		expectFingerNodeID(2, 7)
		expectFingerNodeID(3, 1)
	}
	{
		expectPredecessorID, expectSuccessorIDs, expectFingerNodeID := prepareNodeFingerTests(t, nodes[6])
		expectPredecessorID(5)
		expectSuccessorIDs(7, 0, 1)
		expectFingerNodeID(1, 7)
		expectFingerNodeID(2, 0)
		expectFingerNodeID(3, 2)
	}
	{
		expectPredecessorID, expectSuccessorIDs, expectFingerNodeID := prepareNodeFingerTests(t, nodes[7])
		expectPredecessorID(6)
		expectSuccessorIDs(0, 1, 2)
		expectFingerNodeID(1, 0)
		expectFingerNodeID(2, 1)
		expectFingerNodeID(3, 3)
	}
}

func prepareNodes(ids ...int64) []*localNode {
	nodes := make([]*localNode, len(ids))
	for i, s := range ids {
		nodes[i] = newLocalNodeID(fakeAddr(byte(s)), newID64(s, M3))
	}
	return nodes
}

func prepareNodeFingerTests(t *testing.T, node *localNode) (func(int64), func(...int64), func(int, int64)) {
	expectPredecessorID := func(predecessorID int64) {
		if n := node.predecessor; !n.ID().Eq(newID64(predecessorID, M3)) {
			t.Errorf("{%v}.predecessor expected to be %v, was %v", node, predecessorID, n)
		}
	}
	expectFingerNodeID := func(finger int, nodeID int64) {
		if n := node.fingerNode(finger); !n.ID().Eq(newID64(nodeID, M3)) {
			t.Errorf("{%v}.finger(%v).node expected to be %v, was %v.", node, finger, nodeID, n)
		}
	}
	expectSuccessorIDs := func(succIDs ...int64) {
		if len(succIDs) != len(node.succlist) {
			t.Errorf("len({%v}.succlist) expected to be %v, was %v", node, len(succIDs), len(node.succlist))
			return
		}
		for i, succID := range succIDs {
			succ := node.succlist[i]
			if x := succ.ID().BigInt().Int64(); x != succID {
				t.Errorf("{%v}.succlist[%d] expected to be %v, was %v", node, i, succID, x)
			}
		}
	}
	return expectPredecessorID, expectSuccessorIDs, expectFingerNodeID
}
