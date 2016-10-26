package chord

import "math/big"

// FingerTable holds a collection of related Chord Node fingers.
type FingerTable []finger

func newFingerTable(id *ID) FingerTable {
	fingers := make([]finger, id.Bits()+1)
	for i := range fingers {
		fingers[i] = finger{
			start: *calcFingerStart(id, i),
			node:  nil,
		}
	}
	return FingerTable(fingers)
}

// (n + 2^i) mod (2^m)
func calcFingerStart(id *ID, i int) *ID {
	n := id.BigInt()
	m := big.NewInt(int64(id.Bits()))
	two := big.NewInt(2)

	// addend = 2^i
	addend := big.Int{}
	addend.Exp(two, big.NewInt(int64(i)), nil)

	// sum = n + addend
	sum := big.Int{}
	sum.Add(n, &addend)

	// ceil = 2^m
	ceil := big.Int{}
	ceil.Exp(two, m, nil)

	// result = sum % ceil
	result := big.Int{}
	result.Mod(&sum, &ceil)

	return NewID(&result, id.Bits())
}

// FingerStart resolves finger interval start ID at given finger table offset i.
//
// The result is only defined for i in [1,M+1], where M is the amount of table
// rows.
func (table FingerTable) FingerStart(i int) *ID {
	table.verifyTableIndexOrPanic(i)
	return &table[i-1].start
}

func (table FingerTable) verifyTableIndexOrPanic(i int) {
	verifyIndexOrPanic(len(table), i)
}

// FingerNode resolves Chord node at given finger table offset i.
//
// The result is only defined for i in [1,M], where M is the amount of table
// rows.
func (table FingerTable) FingerNode(i int) Node {
	table.verifyTableIndexOrPanic(i)
	return table.successor(i - 1)
}

// Scans table downwards from provided table offset i until it finds a node,
// which is returned.
//
// The result is only defined for i in [0,M), where M is the amount of table
// rows.
func (table FingerTable) successor(i int) Node {
	for _, sect := range [][]finger{table[i:], table[:i]} {
		for _, fing := range sect {
			if node := fing.node; node != nil {
				return node
			}
		}
	}
	return nil
}

// SetFingerNode attempts to set the node of the ith finger to n.
//
// The operation is only defined for i in [1,M], where M is the amount of table
// rows.
func (table FingerTable) SetFingerNode(i int, n Node) {
	table.verifyTableIndexOrPanic(i)
	table[i-1].node = n
}

// RemoveFingerNodesByID removes all nodes from table matching given ID.
func (table FingerTable) RemoveFingerNodesByID(id *ID) {
	for _, fing := range table {
		if n := fing.node; n != nil && n.ID().Eq(id) {
			fing.node = nil
		}
	}
}

type finger struct {
	start ID
	node  Node
}
