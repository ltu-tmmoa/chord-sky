package chord

import (
	"fmt"
	"testing"
)

func TestNodeJoin2(t *testing.T) {
	nodes := make([]*Node, 2)
	for i, s := range []int64{0, 1} {
		nodes[i] = newNode(stringAddr(fmt.Sprintf("%02d", s)), newHash64(s, M3))
	}

	nodes[0].Join(nil)
	nodes[1].Join(nodes[0])

	for _, node := range nodes {
		fmt.Printf("Node: %v\n", node)
		for i := 1; i <= M3; i++ {
			finger := node.finger(i)
			fmt.Printf("  Finger: %v, %v, %v\n", i, finger.Interval(), finger.Node())
		}
		fmt.Printf("  Predecessor: %v\n", node.predecessor)
	}
}

func TestNodeJoin3(t *testing.T) {
	nodes := make([]*Node, 3)
	for i, s := range []int64{0, 1, 3} {
		nodes[i] = newNode(stringAddr(fmt.Sprintf("%02d", s)), newHash64(s, M3))
	}

	nodes[0].Join(nil)
	nodes[1].Join(nodes[0])
	nodes[2].Join(nodes[1])

	for _, node := range nodes {
		fmt.Printf("Node: %v\n", node)
		for i := 1; i <= M3; i++ {
			finger := node.finger(i)
			fmt.Printf("  Finger: %v, %v, %v\n", i, finger.Interval(), finger.Node())
		}
		fmt.Printf("  Predecessor: %v\n", node.predecessor)
	}
}

/*
func TestNodeJoin(t *testing.T) {
	nodes := make([]*Node, 3)
	for i := range nodes {
		hash := newHash64(int64(i), M3)
		nodes[i] = newNode(stringAddr(fmt.Sprintf("%02d", i)), hash)
	}

	nodes[0].Join(nil)
	nodes[1].Join(nodes[0])
	nodes[2].Join(nodes[1])
	nodes[3].Join(nodes[1])
	nodes[4].Join(nodes[1])
	nodes[5].Join(nodes[4])
	nodes[6].Join(nodes[3])
	nodes[7].Join(nodes[3])

	for _, node := range nodes {
		fmt.Printf("Node: %v\n", node)
		for i := 1; i <= M3; i++ {
			finger := node.finger(i)
			fmt.Printf("  Finger: %v, %v, %v\n", i, finger.Interval(), finger.Node())
		}
		fmt.Printf("  Predecessor: %v\n", node.predecessor)
	}

	//nodes[0].PrintRing()
}
*/