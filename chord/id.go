package chord

import (
	"crypto/sha1"
	"fmt"
	"math/big"
	"net"
)

const (
	hashBitsMax = sha1.Size * 8
)

// ID identifies some Chord node or key.
type ID struct {
	value big.Int
	bits  int
}

// NewID creates from big integer and an amount of significant bits.
func NewID(value *big.Int, bits int) *ID {
	id := new(ID)
	id.value = *value
	id.bits = bits
	id.truncate()
	return id
}

// ParseID takes a string `s` and convert it into an `*ID`.
func ParseID(s string) (*ID, bool) {
	return parseID(s, hashBitsMax)
}

func parseID(s string, bits int) (*ID, bool) {
	value := new(big.Int)
	if _, ok := value.SetString(s, 16); !ok {
		return nil, false
	}
	return NewID(value, bits), true
}

func hashAddr(addr *net.TCPAddr) *ID {
	value := new(big.Int)
	sum := sha1.Sum([]byte(addr.String()))
	value.SetBytes(sum[:])
	return NewID(value, hashBitsMax)
}

func (id *ID) truncate() {
	// ceil = 2^bits
	ceil := big.Int{}
	ceil.Exp(big.NewInt(2), big.NewInt(int64(id.bits)), nil)

	// id = id % ceil
	id.value.Mod(&id.value, &ceil)
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
	diff.bits = id.bits

	// diff = id - other
	diff.value.Sub(&id.value, other.BigInt())

	diff.truncate()
	return diff
}

// Eq determines if this id and given ID are equal.
func (id *ID) Eq(other *ID) bool {
	return id.Cmp(other) == 0
}

// String produces a canonical string representation of this ID.
func (id *ID) String() string {
	return id.value.Text(16)
}

func verifyIndexOrPanic(len, i int) {
	if 1 > i || i > len {
		panic(fmt.Sprintf("%d not in [1,%d]", i, len))
	}
}
