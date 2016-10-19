package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ltu-tmmoa/chord-sky/chord"
	"github.com/ltu-tmmoa/chord-sky/log"
	cnet "github.com/ltu-tmmoa/chord-sky/net"
	chttp "github.com/ltu-tmmoa/chord-sky/net/http"
)

var peer string

func init() {
	flag.StringVar(&peer, "peer", "", "<IP:PORT> of Chord Sky Node to join. If not given a new ring is created.")
}

func main() {
	flag.Parse()

	log.Logger.Println("Chord Sky")

	// Setup chord node representing this machine.
	var localNode *chord.LocalNode
	{
		localIPAddr, err := cnet.GetLocalIPAddr()
		if err != nil {
			log.Logger.Fatalln(err)
		}
		localNode = chord.NewLocalNode(localIPAddr)
	}

	// Join new or existing Chord ring.
	{
		trimmedPeer := strings.TrimSpace(peer)
		if len(trimmedPeer) == 0 {
			log.Logger.Println("No peer specified. Forming new ring ...")

			localNode.Join(nil)

		} else {
			log.Logger.Println("Joining ring via", trimmedPeer, "...")

			ipAddr, err := net.ResolveIPAddr("ip", trimmedPeer)
			if err != nil {
				log.Logger.Fatalln(err)
			}
			localNode.Join(chord.NewRemoteNode(ipAddr))
		}
	}

	localNodeMutex := &sync.RWMutex{}

	// Schedule recurring operations.
	go func() {
		time.Sleep(10 * time.Second)
		for {
			func() {
				// TODO: Heartbeat?

				localNodeMutex.Lock()
				defer localNodeMutex.Unlock()

				log.Logger.Println("Stabilizing ...")
				localNode.Stabilize()

				log.Logger.Println("Fixing random finger table entry ...")
				localNode.FixRandomFinger()
			}()
			time.Sleep(30 * time.Second)
		}
	}()

	// Expose local node as HTTP RPC service.
	{
		httpNode := chttp.NewNode(localNode, localNodeMutex)
		http.Handle("/", httpNode)

		localAddr := fmt.Sprintf("%s:8080", localNode.IPAddr().String())
		log.Logger.Println("Listening on", localAddr, "...")
		log.Logger.Fatal(http.ListenAndServe(localAddr, nil))
	}
}
