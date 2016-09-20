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
	/*nodes[2].Join(nodes[1])
	nodes[3].Join(nodes[1])
	nodes[4].Join(nodes[1])
	nodes[5].Join(nodes[4])
	nodes[6].Join(nodes[3])
	nodes[7].Join(nodes[3])*/

	nodes[0].PrintRing()
}
