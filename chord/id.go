package chord

import (
	"crypto/sha1"
	"fmt"
	"math/big"
)

// ID identifies some Chord node or key.
type ID struct {
	value big.Int
	bits  *int
}

func hash(a interface{}, bits *int) *ID {
	// ceil = 2^bits
	ceil := big.Int{}
	ceil.Exp(big.NewInt(2), big.NewInt(int64(*bits)), nil)

	// hash = sha1(a)
	hash := sha1.Sum([]byte(fmt.Sprint(a)))

	// value = hash % ceil
	value := big.Int{}
	value.SetBytes(hash[:])
	value.Mod(&value, &ceil)

	id := new(ID)
	id.value = value
	id.bits = bits
	return id
}

// Cmp resolves ordering of this and given ID:s.
//
// Returns -1, 0 or 1 if given node is lesser than, equal to, or greater than
// this node.
func (id *ID) Cmp(other *ID) int {
	return id.value.Cmp(&other.value)
}

// Diff calculates the Chord ring distance between this and given ID:s.
func (id *ID) Diff(other *ID) *ID {
	diff := new(ID)

	// diff = id - other
	(&diff.value).Sub(&id.value, &other.value)

	// ceil = 2^bits
	ceil := big.Int{}
	ceil.Exp(big.NewInt(2), big.NewInt(int64(*id.bits)), nil)

	// diff = diff % ceil
	(&diff.value).Mod(&diff.value, &ceil)

	diff.bits = id.bits
	return diff
}

// Eq determines if this and given ID:s are equal.
func (id *ID) Eq(other *ID) bool {
	return id.Cmp(other) == 0
}

// String produces a canonical string representation of this ID.
func (id *ID) String() string {
	return id.value.String()
}
