package main

import (
	"github.com/ltu-tmmoa/chord-sky/log"
	"github.com/ltu-tmmoa/chord-sky/chord"
	"net"
	"errors"
	"time"
)

func main() {
	log.Logger.Println("Chord Sky")

	localAddr, err := getLocalAddr()
	if err != nil {
		log.Logger.Fatalln(err)
	}
	localNode := chord.NewLocalNode(localAddr)

	log.Logger.Println("Local address:", localNode.Addr())

	go func() {
		for {
			log.Logger.Println("Stabilizing ...")
			localNode.Stabilize()
			time.Sleep(30 * time.Second)
		}
	}()

	go func() {
		for {
			log.Logger.Println("Fixing random finger table entry ...")
			localNode.FixRandomFinger()
			time.Sleep(30 * time.Second)
		}
	}()

	// TODO: Accept incoming connections.

}

func getLocalAddr() (net.Addr, error) {
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