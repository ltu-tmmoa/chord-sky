package main

import (
	"flag"
	"net"
	"strings"
	"time"

	"github.com/ltu-tmmoa/chord-sky/chord"
	"github.com/ltu-tmmoa/chord-sky/log"
	cnet "github.com/ltu-tmmoa/chord-sky/net"
)

var peer string
var port int

func init() {
	flag.StringVar(&peer, "peer", "", "<IP:PORT> of Chord Sky Node to join. If not given a new ring is created.")
	flag.IntVar(&port, "port", 8080, "Network port number to use for receiving incoming connections.")
}

func main() {
	flag.Parse()

	log.Logger.Println("Chord Sky")

	// Setup chord node representing this machine.
	var localNode *chord.LocalNode
	{
		localTCPAddr, err := cnet.GetLocalTCPAddr(port)
		if err != nil {
			log.Logger.Fatalln(err)
		}
		localNode = chord.NewLocalNode(localTCPAddr)
	}

	// Join new or existing Chord ring.
	{
		trimmedPeer := strings.TrimSpace(peer)
		if len(trimmedPeer) == 0 {
			log.Logger.Println("No peer specified. Forming new ring ...")

			localNode.Join(nil)

		} else {
			log.Logger.Println("Joining ring via", trimmedPeer, "...")

			tcpAddr, err := net.ResolveTCPAddr("ip", trimmedPeer)
			if err != nil {
				log.Logger.Fatalln(err)
			}
			localNode.Join(chord.NewRemoteNode(tcpAddr))
		}
	}

	// Schedule recurring operations.
	time.Sleep(10 * time.Second)
	for {
		func() {
			log.Logger.Println("Stabilizing ...")
			localNode.Stabilize()

			log.Logger.Println("Fixing random finger table entry ...")
			localNode.FixRandomFinger()
		}()
		time.Sleep(30 * time.Second)
	}
}
