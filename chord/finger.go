package chord

import (
	"fmt"
	"math/big"
)

// Finger represents a Chord node finger.
type Finger struct {
	n *Node
	i int
}

func newFinger(n *Node, i int) *Finger {
	if 1 > i || i > *n.id.bits {
		panic(fmt.Sprintf("i must be in [1,%d]", *n.id.bits))
	}
	finger := new(Finger)
	finger.n = n
	finger.i = i
	return finger
}

// Start yields the Chord node finger[i].start finger ID.
func (finger *Finger) Start() *ID {
	return fingerStart(finger.n, finger.i, *finger.n.id.bits)
}

// (n + 2^(i-1)) mod (2^m)
func fingerStart(n *Node, i, m int) *ID {
	two := big.NewInt(2)

	// addend = 2^(i-1)
	addend := big.Int{}
	addend.Exp(two, big.NewInt(int64(i-1)), nil)

	// sum = n + addend
	sum := big.Int{}
	sum.Add(&n.id.value, &addend)

	// ceil = 2^m
	ceil := big.Int{}
	ceil.Exp(two, big.NewInt(int64(m)), nil)

	// sum = sum % ceil
	result := big.Int{}
	result.Mod(&sum, &ceil)

	id := new(ID)
	id.value = result
	return id
}

// Interval yields the [finger[i].start, finger[i + 1].start) finger ID range.
func (finger *Finger) Interval() *FingerInterval {
	m := *finger.n.id.bits

	fingerInterval := new(FingerInterval)
	fingerInterval.start = fingerStart(finger.n, finger.i, m)
	fingerInterval.stop = fingerStart(finger.n, finger.i+1, m)
	return fingerInterval
}

// Node yields Chord node associated with finger.
func (finger *Finger) Node() *Node {
	return finger.n // TODO: Perform actual node lookup.
}

// FingerInterval holds a Chord node finger interval.
type FingerInterval struct {
	start *ID
	stop  *ID
}

// Contains determines if given Chord is contained in this interval.
func (fingerInterval *FingerInterval) Contains(id *ID) bool {
	return fingerInterval.start.Cmp(id) >= 0 && fingerInterval.stop.Cmp(id) < 0
}

// String produces a canonical string representation of this finger interval.
func (fingerInterval *FingerInterval) String() string {
	return fmt.Sprintf("[%v,%v)", fingerInterval.start, fingerInterval.stop)
}
