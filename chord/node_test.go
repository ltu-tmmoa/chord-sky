package chord

import (
	"fmt"
	"testing"
)

func TestNodeJoin(t *testing.T) {
	nodes := make([]*Node, 8)
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
			start, stop := finger.Interval()
			fmt.Printf("  Finger:      %v\n", i)
			fmt.Printf("    Start:       %v\n", finger.start)
			fmt.Printf("    Interval:    [%v, %v)\n", start.BigInt().Int64(), stop.BigInt().Int64())
			fmt.Printf("    Node:        %v\n", finger.Node())
		}
		fmt.Printf("  Predecessor: %v\n", node.predecessor)
	}

	//nodes[0].PrintRing()
}
