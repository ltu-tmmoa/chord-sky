package chord

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"

	"github.com/ltu-tmmoa/chord-sky/chord"
	"github.com/ltu-tmmoa/chord-sky/log"
)

// Node is used to expose Chord Node operations via RPC.
type Node struct {
	node chord.Node
}

// NewNode creates a new public node, wrapping given regular node.
func NewNode(node chord.Node) *Node {
	httpNode := new(Node)
	httpNode.node = node
	return httpNode
}

// ServeHTTP routes incoming HTTP requests.
func (node *Node) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v", r)
		}
	}()

	log.Logger.Printf("Incoming request: %s %s", r.Method, r.URL.String())

	switch r.URL.Path {
	case "/node/info":
		switch r.Method {
		case http.MethodGet:
			node.getInfo(w, r)
			return
		}
	case "/node/fingers":
		switch r.Method {
		case http.MethodDelete:
			node.deleteFingers(w, r)
			return

		case http.MethodGet:
			node.getFingers(w, r)
			return

		case http.MethodPut:
			node.putFingers(w, r)
			return
		}
	case "/node/successor":
		switch r.Method {
		case http.MethodGet:
			node.getSuccessor(w, r)
			return
		case http.MethodPut:
			node.putSuccessor(w, r)
			return
		}
	case "/node/successors":
		switch r.Method {
		case http.MethodGet:
			node.getSuccessors(w, r)
			return
		}
	case "/node/predecessor":
		switch r.Method {
		case http.MethodGet:
			node.getPredecessor(w, r)
			return
		case http.MethodPut:
			node.putPredecessor(w, r)
			return
		}
	case "/node/predecessors":
		switch r.Method {
		case http.MethodGet:
			node.getPredecessors(w, r)
			return
		}
	}
	node.notFound(w)
}

func (node *Node) getInfo(w http.ResponseWriter, r *http.Request) {
	bits := node.node.ID().Bits()
	var buffer bytes.Buffer
	buffer.WriteString("Predecessor:\n")
	{
		pred, err := node.node.Predecessor()
		if err != nil {
			node.internalServerError(w, err)
			return
		}
		buffer.WriteString(pred.String())
	}
	buffer.WriteString("\n\n")

	buffer.WriteString("Successor:\n")
	{
		succ, err := node.node.Successor()
		if err != nil {
			node.internalServerError(w, err)
			return
		}
		buffer.WriteString(succ.String())
	}
	buffer.WriteString("\n\n")

	buffer.WriteString("Finger table:\n")
	for i := 1; i <= bits; i++ {
		fing, err := node.node.FingerNode(i)
		if err != nil {
			node.internalServerError(w, err)
			return
		}
		buffer.WriteString(fmt.Sprintf("%3d: %v\n", i, fing))
	}
	buffer.WriteString("\n")

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, buffer.String())
}

// Handles HTTP DELETE fingers request.
//
// Example: curl -X DELETE -v '<ip:port>/node/fingers?id=7asdt1dg7awdgj2aj2'
func (node *Node) deleteFingers(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if len(id) == 0 {
		node.badRequest(w, "id param not provided")
		return
	}
	err := node.node.RemoveFingerNodesByID(chord.Identity(id))
	if err != nil {
		node.internalServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Handles HTTP GET fingers request.
//
// Example: curl -v '<ip:port>/node/fingers?id=2'
func (node *Node) getFingers(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if len(id) == 0 {
		node.badRequest(w, "id param not provided")
		return
	}
	index, err := strconv.Atoi(id)
	if err != nil {
		node.badRequest(w, err.Error())
		return
	}
	finger, err := node.node.FingerNode(index)
	if err != nil {
		node.internalServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, finger.IPAddr())
}

// Handles HTTP PUT fingers request.
//
// Example: curl -X PUT -d '127.0.0.1' -v '<ip:port>/node/fingers?id=4'
func (node *Node) putFingers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		node.internalServerError(w, err)
		return
	}
	if len(body) == 0 {
		node.badRequest(w, "no IP provided in body")
		return
	}
	ipAddr, err := net.ResolveIPAddr("ip", string(body))
	if err != nil {
		node.internalServerError(w, err)
		return
	}
	id := r.URL.Query().Get("id")
	if len(id) == 0 {
		node.badRequest(w, "id param not provided")
		return
	}
	index, err := strconv.Atoi(id)
	if err != nil {
		node.badRequest(w, err.Error())
		return
	}
	node.node.SetFingerNode(index, chord.NewRemoteNode(ipAddr))
	w.WriteHeader(http.StatusNoContent)
}

// Handles HTTP GET successor request.
//
// Example: curl -v '<ip:port>/node/successor'
func (node *Node) getSuccessor(w http.ResponseWriter, r *http.Request) {
	successor, err := node.node.Successor()
	if err != nil {
		node.internalServerError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, successor.IPAddr())
}

// Handles HTTP PUT successor request.
//
// Example: curl -X PUT -d '127.0.0.1' -v '<ip:port>/node/successor'
func (node *Node) putSuccessor(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		node.internalServerError(w, err)
		return
	}
	if len(body) == 0 {
		node.badRequest(w, "no IP provided in body")
		return
	}
	ipAddr, err := net.ResolveIPAddr("ip", string(body))
	if err != nil {
		node.internalServerError(w, err)
		return
	}
	node.node.SetSuccessor(chord.NewRemoteNode(ipAddr))
	w.WriteHeader(http.StatusNoContent)
}

// Handles HTTP GET successors request.
//
// Example: curl -v '<ip:port>/node/successors?id=1'
func (node *Node) getSuccessors(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if len(id) == 0 {
		node.badRequest(w, "id param not provided")
		return
	}
	successor, err := node.node.FindSuccessor(chord.Identity(id))
	if err != nil {
		node.internalServerError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, successor.IPAddr())
}

// Handles HTTP GET predecessor request.
//
// Example: curl -v '<ip:port>/node/predecessor'
func (node *Node) getPredecessor(w http.ResponseWriter, r *http.Request) {
	predecessor, err := node.node.Predecessor()
	if err != nil {
		node.internalServerError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, predecessor.IPAddr())
}

// Handles HTTP PUT predecessor request.
//
// Example: curl -X PUT -d '127.0.0.1' -v '<ip:port>/node/predecessor'
func (node *Node) putPredecessor(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		node.internalServerError(w, err)
		return
	}
	if len(body) == 0 {
		node.badRequest(w, "no IP provided in body")
		return
	}
	ipAddr, err := net.ResolveIPAddr("ip", string(body))
	if err != nil {
		node.internalServerError(w, err)
		return
	}
	node.node.SetPredecessor(chord.NewRemoteNode(ipAddr))
	w.WriteHeader(http.StatusNoContent)
}

// Handles HTTP GET predecessors request.
//
// Example: curl -v '<ip:port>/node/predecessors?id=4'
func (node *Node) getPredecessors(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if len(id) == 0 {
		node.badRequest(w, "id param not provided")
		return
	}
	predecessor, err := node.node.FindPredecessor(chord.Identity(id))
	if err != nil {
		node.internalServerError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, predecessor.IPAddr())
}

func (node *Node) notFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func (node *Node) internalServerError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, err.Error())
}

func (node *Node) badRequest(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, message)
}
