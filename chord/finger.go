package chord

import (
	"fmt"
	"math/big"
)

// Finger represents a Chord node finger.
type Finger struct {
	interval FingerInterval
	node     func() (Node, error)
}

func newFinger(id ID, i int) *Finger {
	finger := new(Finger)
	finger.interval = FingerInterval{
		start: fingerStart(id, i),
		stop:  fingerStart(id, i+1),
	}
	finger.node = nil
	return finger
}

// (n + 2^(i-1)) mod (2^m)
func fingerStart(id ID, i int) ID {
	n := id.BigInt()
	m := big.NewInt(int64(id.Bits()))
	two := big.NewInt(2)

	// addend = 2^(i-1)
	addend := big.Int{}
	addend.Exp(two, big.NewInt(int64(i-1)), nil)

	// sum = n + addend
	sum := big.Int{}
	sum.Add(n, &addend)

	// ceil = 2^m
	ceil := big.Int{}
	ceil.Exp(two, m, nil)

	// result = sum % ceil
	result := big.Int{}
	result.Mod(&sum, &ceil)

	return newHash(result, id.Bits())
}

// Start yields the Chord node finger[i].start finger ID.
func (finger *Finger) Start() ID {
	return finger.interval.start
}

// Interval yields finger[i].start and finger[i + 1].start.
func (finger *Finger) Interval() *FingerInterval {
	return &finger.interval
}

// Node yields Chord node associated with finger.
func (finger *Finger) Node() (Node, error) {
	return finger.node()
}

// SetNode sets known node as finger node.
func (finger *Finger) SetNode(node Node) {
	finger.node = func() (Node, error) {
		return node, nil
	}
}

// SetsNodeLazy sets function used to resolve finger node when requested.
func (finger *Finger) SetNodeLazy(node func() (Node, error)) {
	finger.node = node
}

// String produces a canonical string representation of this Finger.
func (finger *Finger) String() string {
	var node interface{}
	node, err := finger.Node()
	if err != nil {
		node = err.Error()
	}
	return fmt.Sprintf("%v (%v)", finger.interval.String(), node)
}

// FingerInterval holds two ID:s, representing a [start, stop) range of ID:s.
type FingerInterval struct {
	start ID
	stop  ID
}

// Start yields the Chord node finger[i].start finger ID.
func (interval *FingerInterval) Start() ID {
	return interval.start
}

// Stop yields the Chord node finger[i + 1].start finger ID.
func (interval *FingerInterval) Stop() ID {
	return interval.stop
}

// Contains determines if given ID is contained within the interval.
func (interval *FingerInterval) Contains(other ID) bool {
	return idIntervalContainsIE(interval.start, interval.stop, other)
}

// String produces a canonical string representation of this Finger.
func (interval *FingerInterval) String() string {
	return fmt.Sprintf("[%v, %v)", interval.start, interval.stop)
}
