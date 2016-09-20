package chord

import (
	"crypto/sha1"
	"fmt"
	"math/big"
)

// Hash identifies some Chord node or key.
type Hash struct {
	value big.Int
	bits  int
}

func newHash(value big.Int, bits int) *Hash {
	hash := new(Hash)
	hash.value = value
	hash.bits = bits
	return hash
}

func hash(a interface{}, bits int) *Hash {
	// ceil = 2^bits
	ceil := big.Int{}
	ceil.Exp(big.NewInt(2), big.NewInt(int64(bits)), nil)

	// sum = sha1(a)
	sum := sha1.Sum([]byte(fmt.Sprint(a)))

	// value = sum % ceil
	value := big.Int{}
	value.SetBytes(sum[:])
	value.Mod(&value, &ceil)

	return newHash(value, bits)
}

// AsInt turns hash into big.Int representation.
func (hash *Hash) AsInt() *big.Int {
	return &hash.value
}

// Bits returns amount of significant bits in hash.
func (hash *Hash) Bits() int {
	return hash.bits
}

// Cmp compares hash to given ID.
func (hash *Hash) Cmp(other ID) int {
	return hash.value.Cmp(other.AsInt())
}

// Diff calculates the difference between this hash and given ID.
func (hash *Hash) Diff(other ID) ID {
	diff := new(Hash)

	// diff = hash - other
	diff.value.Sub(&hash.value, other.AsInt())

	// ceil = 2^bits
	//ceil := big.Int{}
	//ceil.Exp(big.NewInt(2), big.NewInt(int64(hash.bits)), nil)

	// diff = diff % ceil
	//diff.value.Mod(&diff.value, &ceil)

	diff.bits = hash.bits
	return diff
}

// Eq determines if this hash and given ID are equal.
func (hash *Hash) Eq(other ID) bool {
	return hash.Cmp(other) == 0
}

// String produces a canonical string representation of this Hash.
func (hash *Hash) String() string {
	return hash.value.String()
}
