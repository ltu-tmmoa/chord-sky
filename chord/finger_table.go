package chord

import (
	"bytes"
	"fmt"
	"math/big"
)

// Holds a collection of related Chord Node fingers.
type fingerTable struct {
	owner   Node
	fingers []finger
}

func newFingerTable(owner Node) *fingerTable {
	id := owner.ID()
	fingers := make([]finger, id.Bits()+1)
	for i := range fingers {
		fingers[i] = finger{
			start: *calcfingerStart(id, i),
			node:  nil,
		}
	}
	return &fingerTable{
		owner:   owner,
		fingers: fingers,
	}
}

// n + 2^i
func calcfingerStart(id *ID, i int) *ID {
	n := id.BigInt()

	// addend = 2^i
	addend := big.Int{}
	addend.Exp(big.NewInt(2), big.NewInt(int64(i)), nil)

	// result = n + addend
	result := big.Int{}
	result.Add(n, &addend)

	return NewID(&result, id.Bits())
}

// Resolves finger interval start ID at given finger table offset i.
//
// The result is only defined for i in [1,M+1], where M is the amount of table
// rows.
func (table *fingerTable) fingerStart(i int) *ID {
	table.verifyTableIndexOrPanic(i)
	return &table.fingers[i-1].start
}

func (table *fingerTable) verifyTableIndexOrPanic(i int) {
	verifyIndexOrPanic(len(table.fingers), i)
}

// Resolves Chord node at given finger table offset i.
//
// The result is only defined for i in [1,M], where M is the amount of table
// rows.
func (table *fingerTable) fingerNode(i int) Node {
	table.verifyTableIndexOrPanic(i)
	return table.successor(i - 1)
}

// Scans table downwards from provided table offset i until it finds a node,
// which is returned.
//
// The result is only defined for i in [0,M), where M is the amount of table
// rows.
func (table *fingerTable) successor(i int) Node {
	for _, sect := range [][]finger{table.fingers[i:], table.fingers[:i]} {
		for _, fing := range sect {
			if node := fing.node; node != nil {
				return node
			}
		}
	}
	return table.owner
}

// Attempts to set the node of the ith finger to n.
//
// The operation is only defined for i in [1,M], where M is the amount of table
// rows.
func (table *fingerTable) setFingerNode(i int, n Node) {
	table.verifyTableIndexOrPanic(i)
	table.fingers[i-1].node = n
}

// Removes all nodes from table matching given ID.
func (table *fingerTable) removeFingerNodesByID(id *ID) {
	for i, fing := range table.fingers {
		if n := fing.node; n != nil && n.ID().Eq(id) {
			table.fingers[i].node = nil
		}
	}
}

func (table *fingerTable) String() string {
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
