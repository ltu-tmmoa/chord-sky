package chord

import (
	"bytes"
	"fmt"
	"math/big"
)

// FingerTable holds a collection of related Chord Node fingers.
type FingerTable struct {
	owner   Node
	fingers []finger
}

func newFingerTable(owner Node) *FingerTable {
	id := owner.ID()
	fingers := make([]finger, id.Bits()+1)
	for i := range fingers {
		fingers[i] = finger{
			start: *calcFingerStart(id, i),
			node:  nil,
		}
	}
	return &FingerTable{
		owner:   owner,
		fingers: fingers,
	}
}

// n + 2^i
func calcFingerStart(id *ID, i int) *ID {
	n := id.BigInt()

	// addend = 2^i
	addend := big.Int{}
	addend.Exp(big.NewInt(2), big.NewInt(int64(i)), nil)

	// result = n + addend
	result := big.Int{}
	result.Add(n, &addend)

	return NewID(&result, id.Bits())
}

// FingerStart resolves finger interval start ID at given finger table offset i.
//
// The result is only defined for i in [1,M+1], where M is the amount of table
// rows.
func (table *FingerTable) FingerStart(i int) *ID {
	table.verifyTableIndexOrPanic(i)
	return &table.fingers[i-1].start
}

func (table *FingerTable) verifyTableIndexOrPanic(i int) {
	verifyIndexOrPanic(len(table.fingers), i)
}

// FingerNode resolves Chord node at given finger table offset i.
//
// The result is only defined for i in [1,M], where M is the amount of table
// rows.
func (table *FingerTable) FingerNode(i int) Node {
	table.verifyTableIndexOrPanic(i)
	return table.successor(i - 1)
}

// Scans table downwards from provided table offset i until it finds a node,
// which is returned.
//
// The result is only defined for i in [0,M), where M is the amount of table
// rows.
func (table *FingerTable) successor(i int) Node {
	for _, sect := range [][]finger{table.fingers[i:], table.fingers[:i]} {
		for _, fing := range sect {
			if node := fing.node; node != nil {
				return node
			}
		}
	}
	return table.owner
}

// SetFingerNode attempts to set the node of the ith finger to n.
//
// The operation is only defined for i in [1,M], where M is the amount of table
// rows.
func (table *FingerTable) SetFingerNode(i int, n Node) {
	table.verifyTableIndexOrPanic(i)
	table.fingers[i-1].node = n
}

// RemoveFingerNodesByID removes all nodes from table matching given ID.
func (table *FingerTable) RemoveFingerNodesByID(id *ID) {
	for i, fing := range table.fingers {
		if n := fing.node; n != nil && n.ID().Eq(id) {
			table.fingers[i].node = nil
		}
	}
}

func (table *FingerTable) String() string {
	buf := &bytes.Buffer{}
	for i, finger := range table.fingers {
		fmt.Fprintf(buf, "%3d: %s\r\n", i, finger.node)
	}
	return string(buf.Bytes())
}

type finger struct {
	start ID
	node  Node
}
