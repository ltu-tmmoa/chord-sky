package chord

import (
	"math/big"
	"testing"
)

const (
	// M3 denotes the use of 3 bits in Chord ID:s.
	M3 = 3
)

func TestFingerStart(t *testing.T) {
	expectNodeFingerStart := func(table *FingerTable, i int, expected int64) {
		actual := table.FingerStart(i)

		if actual.BigInt().Int64() != expected {
			t.Errorf("finger[%d].start %v != expected %v", i, actual, expected)
		}
	}
	// See figure 3(a), Chord paper.
	{
		a := newTable(1)
		expectNodeFingerStart(a, 1, 2)
		expectNodeFingerStart(a, 2, 3)
		expectNodeFingerStart(a, 3, 5)
	}
}

func TestFingerInterval(t *testing.T) {
	expectNodeFingerInterval := func(table *FingerTable, i int, expectedStart, expectedStop int64) {
		actualStart := table.FingerStart(i)
		actualStop := table.FingerStart(i + 1)

		if actualStart.BigInt().Int64() != expectedStart || actualStop.BigInt().Int64() != expectedStop {
			t.Errorf("finger[%d].interval [%v,%v) != expected [%v,%v)", i, actualStart, actualStop, expectedStart, expectedStop)
		}
	}
	// See figure 3(b), Chord paper.
	{
		a := newTable(0)
		expectNodeFingerInterval(a, 1, 1, 2)
		expectNodeFingerInterval(a, 2, 2, 4)
		expectNodeFingerInterval(a, 3, 4, 0)

		b := newTable(1)
		expectNodeFingerInterval(b, 1, 2, 3)
		expectNodeFingerInterval(b, 2, 3, 5)
		expectNodeFingerInterval(b, 3, 5, 1)

		c := newTable(3)
		expectNodeFingerInterval(c, 1, 4, 5)
		expectNodeFingerInterval(c, 2, 5, 7)
		expectNodeFingerInterval(c, 3, 7, 3)
	}
}

func newTable(id int64) *FingerTable {
	return newFingerTable(&RemoteNode{
		id: *NewID(big.NewInt(id), M3),
	})
}
