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

func newHash64(value int64, bits int) *Hash {
	return newHash(*big.NewInt(value), bits)
}
