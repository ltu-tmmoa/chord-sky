package net

import (
	"errors"
	"net"
)

// GetLocalIPAddr returns a local non-loopback network address.
func GetLocalIPAddr() (*net.IPAddr, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return addr.(*net.IPAddr), nil
			}
		}
	}
	return nil, errors.New("No suitable IP interface available.")
}
