package chord

import (
	"fmt"
	"math/big"
)

// Finger represents a Chord node finger interval.
type Finger struct {
	start ID
	stop  ID
}

func newFinger(id *ID, i int) *Finger {
	return &Finger{
		start: *fingerStart(id, i),
		stop:  *fingerStart(id, i+1),
	}
}

// (n + 2^(i-1)) mod (2^m)
func fingerStart(id *ID, i int) *ID {
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

	return newID(result, id.Bits())
}

// Start yields the Chord node finger[i].start finger ID.
func (finger *Finger) Start() *ID {
	return &finger.start
}

// Stop yields the Chord node finger[i + 1].start finger ID.
func (finger *Finger) Stop() *ID {
	return &finger.stop
}

// String produces a canonical string representation of this Finger.
func (finger *Finger) String() string {
	return fmt.Sprintf("[%v, %v)", finger.start, finger.stop)
}
