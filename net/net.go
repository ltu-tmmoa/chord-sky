package net

import (
	"errors"
	"net"
)

// GetLocalAddr returns a local non-loopback network address.
func GetLocalAddr() (net.Addr, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return addr, nil
			}
		}
	}
	return nil, errors.New("No suitable IP interface available.")
}
