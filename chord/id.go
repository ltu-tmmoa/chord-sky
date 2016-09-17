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

func (id *ID) cmp(other *ID) int {
	return id.value.Cmp(&other.value)
}

func (id *ID) diff(other *big.Int) *ID {
	diff := new(ID)

	// diff = id - other
	(&diff.value).Sub(&id.value, other)

	// ceil = 2^bits
	ceil := big.Int{}
	ceil.Exp(big.NewInt(2), big.NewInt(int64(*id.bits)), nil)

	// diff = diff % ceil
	(&diff.value).Mod(&diff.value, &ceil)

	diff.bits = id.bits
	return diff
}

func (id *ID) eq(other *ID) bool {
	return id.cmp(other) == 0
}

// String produces a canonical string representation of this ID.
func (id *ID) String() string {
	return id.value.String()
}
