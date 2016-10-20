package chord

import (
	"sync"
	"github.com/ltu-tmmoa/chord-sky/chord"
	"net/http"
	"fmt"
	  "strconv"
	  "net"
)

// PublicNode is used to expose Chord Node operations via RPC.
type Node struct {
	node  chord.Node
	mutex *sync.RWMutex
}

// NewPublicNode creates a new public node, wrapping given regular node, synchronizing it using provided mutex.
func NewNode(node chord.Node, mutex *sync.RWMutex) *Node {
	httpNode := new(Node)
	httpNode.node = node
	httpNode.mutex = mutex
	return httpNode
}

// HTTP router
func (node *Node) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/node/successor":
		switch r.Method {
		case "GET":
			node.getSuccessor(w, r)
			return
		case "PUT":
			node.putSuccessor(w, r)
			return
		}
	case "/node/successors":
		switch r.Method {
		case "GET":
			node.getSuccessors(w, r)
			return
		}
	case "/node/predecessor":
		switch r.Method {
		case "GET":
			node.getPredecessor(w, r)
			return
		case "PUT":
			node.putPredecessor(w, r)
			return
		}
	case "/node/predecessors":
		switch r.Method {
		case "GET":
			node.getPredecessors(w, r)
			return
		}
	case "/node/finger":
		  switch r.Method {
		  case "GET":
			    node.getFinger(w, r)
			    return
		  }
	}
	node.notFound(w)
}

// Handles HTTP GET successor request.
//
// Example: curl -v '<ip:port>/node/successor'

func (node *Node) getFinger(w http.ResponseWriter, r *http.Request) {
	  node.mutex.RLock()
	  defer node.mutex.RUnlock()

	  fing := r.URL.Query().Get("finger")
	  if len(fing) == 0 {
		    node.badRequest(w, "finger param not provided")
		    return
	  }
	  finger, _ := strconv.Atoi(fing)

	  res := node.node.Finger(finger).Node()

	  w.WriteHeader(200)
	  fmt.Fprint(w, res.IPAddr())
}

func (node *Node) getSuccessor(w http.ResponseWriter, r *http.Request) {
	node.mutex.RLock()
	defer node.mutex.RUnlock()

	successor, err := node.node.Successor()
	if err != nil {
		node.internalServerError(w, err)
		return
	}
	w.WriteHeader(200)
	fmt.Fprint(w, successor.IPAddr())
}

func (node *Node) getSuccessors(w http.ResponseWriter, r *http.Request) {
	node.mutex.RLock()
	defer node.mutex.RUnlock()

	id := r.URL.Query().Get("id")
	if len(id) == 0 {
		node.badRequest(w, "id param not provided")
		return
	}
	successor, err := node.node.FindSuccessor(chord.NewHash(id))
	if err != nil {
		node.internalServerError(w, err)
		return
	}
	w.WriteHeader(200)
	fmt.Fprint(w, successor.IPAddr())
}

func (node *Node) getPredecessor(w http.ResponseWriter, r *http.Request) {
	node.mutex.RLock()
	defer node.mutex.RUnlock()

	predecessor, err := node.node.Predecessor()
	if err != nil {
		node.internalServerError(w, err)
		return
	}
	w.WriteHeader(200)
	fmt.Fprint(w, predecessor.IPAddr())
}

func (node *Node) getPredecessors(w http.ResponseWriter, r *http.Request) {
	node.mutex.RLock()
	defer node.mutex.RUnlock()

	id := r.URL.Query().Get("id")
	if len(id) == 0 {
		node.badRequest(w, "id param not provided")
		return
	}
	predecessor, err := node.node.FindPredecessor(chord.NewHash(id))
	if err != nil {
		node.internalServerError(w, err)
		return
	}
	w.WriteHeader(200)
	fmt.Fprint(w, predecessor.IPAddr())
}

func (node *Node) putSuccessor(w http.ResponseWriter, r *http.Request) {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	ipPara := r.URL.Query().Get("ip")
	if len(ipPara) == 0 {
		node.badRequest(w, "ip param not provided\n")
		return
	}
	  addr , err := net.ResolveIPAddr("ip", ipPara)

	  if err != nil {
		    node.internalServerError(w, err)
		    return
	  }
	node.node.SetSuccessor(chord.NewRemoteNode(addr))
	w.WriteHeader(201)
}

func (node *Node) putPredecessor(w http.ResponseWriter, r *http.Request) {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	  ipPara := r.URL.Query().Get("ip")
	  if len(ipPara) == 0 {
		    node.badRequest(w, "ip param not provided\n")
		    return
	  }
	  addr , err := net.ResolveIPAddr("ip", ipPara)

	  if err != nil {
		    node.internalServerError(w, err)
		    return
	  }
	  node.node.SetPredecessor(chord.NewRemoteNode(addr))
	  w.WriteHeader(201)
}


func (node *Node) notFound(w http.ResponseWriter) {
	w.WriteHeader(404)
}

func (node *Node) internalServerError(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	fmt.Fprint(w, err.Error())
}

func (node *Node) badRequest(w http.ResponseWriter, message string) {
	w.WriteHeader(400)
	fmt.Fprint(w, message)
}