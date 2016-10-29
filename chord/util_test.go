package chord

import (
	"math/big"
	"net"

	"github.com/ltu-tmmoa/chord-sky/data"
)

func fakeAddr(id byte) *net.TCPAddr {
	return &net.TCPAddr{
		IP:   []byte{192, 168, 1, id},
		Port: 8080,
	}
}

func newID64(value int64, bits int) *data.ID {
	return data.NewID(big.NewInt(value), bits)
}
