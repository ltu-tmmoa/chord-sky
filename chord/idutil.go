package chord

import (
	"crypto/sha1"
	"math/big"
	"net"

	"github.com/ltu-tmmoa/chord-sky/data"
)

const (
	idBits = sha1.Size * 8
)

func parseID(s string) (*data.ID, bool) {
	return data.ParseID(s, idBits)
}

func addrToID(addr *net.TCPAddr) *data.ID {
	value := new(big.Int)
	sum := sha1.Sum([]byte(addr.String()))
	value.SetBytes(sum[:])
	return data.NewID(value, idBits)
}
