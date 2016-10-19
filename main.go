package main

import (
	"time"

	"github.com/ltu-tmmoa/chord-sky/chord"
	"github.com/ltu-tmmoa/chord-sky/log"
	"github.com/ltu-tmmoa/chord-sky/net"
)

func main() {
	log.Logger.Println("Chord Sky")

	localIPAddr, err := net.GetLocalIPAddr()
	if err != nil {
		log.Logger.Fatalln(err)
	}
	localNode := chord.NewLocalNode(localIPAddr)

	log.Logger.Println("Local address:", localNode.IPAddr())

	go func() {
		for {
			log.Logger.Println("Stabilizing ...")
			localNode.Stabilize()

			log.Logger.Println("Fixing random finger table entry ...")
			localNode.FixRandomFinger()

			time.Sleep(30 * time.Second)
		}
	}()

	// TODO: Accept incoming connections.
}
