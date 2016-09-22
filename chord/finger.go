package chord

import (
	"math/big"
	"fmt"
)

// Finger represents a Chord node finger.
type Finger struct {
	start ID
	stop  ID
	node  *Node
}

func newFinger(id ID, i int) *Finger {
	finger := new(Finger)
	finger.start = fingerStart(id, i)
	finger.stop = fingerStart(id, i+1)
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
	return finger.start
}

// Interval yields finger[i].start and finger[i + 1].start.
func (finger *Finger) Interval() (ID, ID) {
	return finger.start, finger.stop
}

// Node yields Chord node associated with finger.
func (finger *Finger) Node() *Node {
	return finger.node
}

// String produces a canonical string representation of this Finger.
func (finger *Finger) String() string {
	return fmt.Sprintf("[%v, %v) (%v)", finger.start, finger.stop, finger.node)
}
