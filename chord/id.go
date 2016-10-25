package chord

import (
	"crypto/sha1"
	"fmt"
	"math/big"
)

const (
	// HashBitsMax represents the maximum allowed number of bits in ID objects.
	HashBitsMax = sha1.Size * 8
)

// ID identifies some Chord node or key.
type ID struct {
	value big.Int
	bits  int
}

func newID(value big.Int, bits int) *ID {
	id := new(ID)
	id.value = value
	id.bits = bits
	return id
}

// NewID creates ID from given object a.
func NewID(a interface{}) *ID {
	return identity(a, HashBitsMax)
}

func identity(a interface{}, bits int) *ID {
	// ceil = 2^bits
	ceil := big.Int{}
	ceil.Exp(big.NewInt(2), big.NewInt(int64(bits)), nil)

	// sum = sha1(a)
	sum := sha1.Sum([]byte(fmt.Sprint(a)))

	// value = sum % ceil
	value := big.Int{}
	value.SetBytes(sum[:])
	value.Mod(&value, &ceil)

	return newID(value, bits)
}

// BigInt turns id into big.Int representation.
func (id *ID) BigInt() *big.Int {
	return &id.value
}

// Bits returns amount of significant bits in id.
func (id *ID) Bits() int {
	return id.bits
}

// Cmp compares id to given ID.
func (id *ID) Cmp(other *ID) int {
	return id.value.Cmp(other.BigInt())
}

// Diff calculates the difference between this id and given ID.
func (id *ID) Diff(other *ID) *ID {
	diff := new(ID)

	// diff = id - other
	diff.value.Sub(&id.value, other.BigInt())

	// ceil = 2^bits
	ceil := big.Int{}
	ceil.Exp(big.NewInt(2), big.NewInt(int64(id.bits)), nil)

	// diff = diff % ceil
	diff.value.Mod(&diff.value, &ceil)

	diff.bits = id.bits
	return diff
}

// Eq determines if this id and given ID are equal.
func (id *ID) Eq(other *ID) bool {
	return id.Cmp(other) == 0
}

// String produces a canonical string representation of this ID.
func (id *ID) String() string {
	return id.value.String()
}
