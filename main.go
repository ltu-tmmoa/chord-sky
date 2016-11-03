package main

import (
	"flag"
	"net"
	"net/http"
	"runtime"
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
	// Force goroutine scheduling to be confined to one thread to avoid having
	// to lock anything.
	runtime.GOMAXPROCS(1)

	flag.Parse()

	log.Logger.Println("Chord Sky")
	http.DefaultClient.Timeout = 5 * time.Second

	laddr, err := cnet.GetLocalTCPAddr(port)
	if err != nil {
		log.Logger.Fatalln(err)
	}
	chordService := chord.NewHTTPService(laddr)

	trimmedPeer := strings.TrimSpace(peer)
	if len(trimmedPeer) == 0 {
		log.Logger.Println("No peer specified. Forming new ring ...")
		chordService.Join(nil)

	} else {
		addr, err := net.ResolveTCPAddr("tcp", trimmedPeer)
		if err != nil {
			log.Logger.Fatalln(err)
		}
		log.Logger.Println("Joining ring via", trimmedPeer, "...")
		chordService.Join(addr)
	}

	go func() {
		for {
			time.Sleep(10 * time.Second)
			log.Logger.Println("Refreshing ...")
			if err := chordService.Refresh(); err != nil {
				log.Logger.Printf("Refresh error: %s", err.Error())
			}
		}
	}()

	storageService := chord.NewHTTPStorageService()

	log.Logger.Println("Accepting incoming connections on", laddr, "...")
	http.Handle("/node/", http.StripPrefix("/node", chordService))
	http.Handle("/storage/", http.StripPrefix("/storage", storageService))
	httpServer := http.Server{
		Addr:         laddr.String(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	httpServer.SetKeepAlivesEnabled(false)
	log.Logger.Fatalln(httpServer.ListenAndServe())
}
