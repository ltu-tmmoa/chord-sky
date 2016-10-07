package chord

import "math/big"

// ID identifies some Chord node or key.
type ID interface {
	// BigInt turns ID into big.Int representation.
	BigInt() *big.Int

	// Bits returns amount of significant bits in ID.
	Bits() int

	// Cmp compares this ID with given ID.
	//
	// Returns -1, 0 or 1 depending on if given other ID is lesser than, equal
	// to, or greater than this ID.
	Cmp(other ID) int

	// Diff calculates the difference between this and given other ID.
	Diff(other ID) ID

	// Eq determines if this and given other ID are equal.
	Eq(other ID) bool

	// String turns ID into its canonical string representation.
	String() string
}
