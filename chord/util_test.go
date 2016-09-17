package chord

import "math/big"

// stringAddr allows a regular string to be treated as a net.Addr.
type stringAddr string

func (s stringAddr) Network() string {
	return string(s)
}

func (s stringAddr) String() string {
	return string(s)
}

// newID creates new ID, without hashing given value a.
func newID(a, bits int) *ID {
	id := new(ID)
	id.value = *big.NewInt(int64(a))
	id.bits = &bits
	return id
}
