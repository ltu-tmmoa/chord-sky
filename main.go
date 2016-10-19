package main

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"sync"
	"time"

	"github.com/ltu-tmmoa/chord-sky/chord"
	"github.com/ltu-tmmoa/chord-sky/log"
	cnet "github.com/ltu-tmmoa/chord-sky/net"
)

func main() {
	log.Logger.Println("Chord Sky")

	// Setup chord node representing this machine.
	var localNode *chord.LocalNode
	{
		localIPAddr, err := cnet.GetLocalIPAddr()
		if err != nil {
			log.Logger.Fatalln(err)
		}
		localNode = chord.NewLocalNode(localIPAddr)
		log.Logger.Println("Local address:", localNode.IPAddr())
	}

	localNodeMutex := &sync.RWMutex{}

	// Schedule recurring operations.
	go func() {
		for {
			time.Sleep(30 * time.Second)

			localNodeMutex.Lock()
			defer localNodeMutex.Unlock()

			log.Logger.Println("Stabilizing ...")
			localNode.Stabilize()

			log.Logger.Println("Fixing random finger table entry ...")
			localNode.FixRandomFinger()
		}
	}()

	// Expose local node as HTTP RPC service.
	{
		publicNode := chord.NewPublicNode(localNode, localNodeMutex)

		rpc.Register(publicNode)
		rpc.HandleHTTP()

		listner, err := net.Listen("tcp", fmt.Sprintf("%s:8080", localNode.IPAddr().String()))
		if err != nil {
			log.Logger.Fatalln(err)
		}

		http.Serve(listner, nil)
	}
}
