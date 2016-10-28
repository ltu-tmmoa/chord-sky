package chord

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var (
	errBodyMissing = errors.New("No body in request. Required.")
)

// HTTPService manages a local Chord node, exposing it as an HTTP service by
// implementing the http.Handler interface.
type HTTPService struct {
	pool     *NodePool
	router   *mux.Router
	isJoined bool
}

// NewHTTPService creates a new HTTP node, exposable as a service on the
// identified local TCP interface.
func NewHTTPService(laddr *net.TCPAddr) *HTTPService {
	service := HTTPService{
		pool:   NewNodePool(laddr),
		router: mux.NewRouter(),
	}

	pool := service.pool
	lnode := pool.lnode

	service.router.
		HandleFunc("/fingers/{i:[0-9]+}", func(w http.ResponseWriter, req *http.Request) {
			i, _ := strconv.Atoi(mux.Vars(req)["i"])
			node := <-lnode.FingerNode(i)
			httpWrite(w, http.StatusOK, node.TCPAddr())
		}).
		Methods(http.MethodGet)

	service.router.
		HandleFunc("/fingers/{i:[0-9]+}", func(w http.ResponseWriter, req *http.Request) {
			i, _ := strconv.Atoi(mux.Vars(req)["i"])
			addr, err := httpReadBodyAsAddr(req)
			if err != nil {
				httpWrite(w, http.StatusBadRequest, err.Error())
				return
			}
			node := pool.GetOrCreateNode(addr)
			lnode.SetFingerNode(i, node)
			w.WriteHeader(http.StatusNoContent)
		}).
		Methods(http.MethodPut)

	service.router.
		HandleFunc("/successor", func(w http.ResponseWriter, req *http.Request) {
			succ := <-lnode.Successor()
			httpWrite(w, http.StatusOK, succ.TCPAddr())
		}).
		Methods(http.MethodGet)

	service.router.
		HandleFunc("/predecessor", func(w http.ResponseWriter, req *http.Request) {
			pred := <-lnode.Predecessor()
			httpWrite(w, http.StatusOK, pred.TCPAddr())
		}).
		Methods(http.MethodGet)

	service.router.
		HandleFunc("/successors", func(w http.ResponseWriter, req *http.Request) {
			id, err := httpReadQueryID(req)
			if err != nil {
				httpWrite(w, http.StatusBadRequest, err.Error())
				return
			}
			node := <-lnode.FindSuccessor(id)
			httpWrite(w, http.StatusOK, node.TCPAddr())
		}).
		Methods(http.MethodGet)

	service.router.
		HandleFunc("/predecessors", func(w http.ResponseWriter, req *http.Request) {
			id, err := httpReadQueryID(req)
			if err != nil {
				httpWrite(w, http.StatusBadRequest, err.Error())
				return
			}
			node := <-lnode.FindPredecessor(id)
			httpWrite(w, http.StatusOK, node.TCPAddr())
		}).
		Methods(http.MethodGet)

	service.router.
		HandleFunc("/successor", func(w http.ResponseWriter, req *http.Request) {
			addr, err := httpReadBodyAsAddr(req)
			if err != nil {
				httpWrite(w, http.StatusBadRequest, err.Error())
				return
			}
			succ := pool.GetOrCreateNode(addr)
			lnode.SetSuccessor(succ)
			w.WriteHeader(http.StatusNoContent)
		}).
		Methods(http.MethodPut)

	service.router.
		HandleFunc("/predecessor", func(w http.ResponseWriter, req *http.Request) {
			addr, err := httpReadBodyAsAddr(req)
			if err != nil {
				httpWrite(w, http.StatusBadRequest, err.Error())
				return
			}
			pred := pool.GetOrCreateNode(addr)
			lnode.SetPredecessor(pred)
			w.WriteHeader(http.StatusNoContent)
		}).
		Methods(http.MethodPut)

	return &service
}

func httpWrite(w http.ResponseWriter, status int, body interface{}) {
	w.WriteHeader(status)
	fmt.Fprint(w, body)
}

func httpReadBody(req *http.Request) (string, error) {
	body := req.Body
	if body == nil {
		return "", errBodyMissing
	}
	defer body.Close()
	arr, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return "", err
	}
	if len(arr) == 0 {
		return "", errBodyMissing
	}
	return string(arr), nil
}

func httpReadBodyAsAddr(req *http.Request) (*net.TCPAddr, error) {
	body, err := httpReadBody(req)
	if err != nil {
		return nil, err
	}
	return net.ResolveTCPAddr("tcp", body)
}

func httpReadQueryID(req *http.Request) (*ID, error) {
	strID := req.URL.Query().Get("id")
	if len(strID) == 0 {
		return nil, errors.New("Query parameter `id` required.")
	}
	id, ok := ParseID(strID)
	if !ok {
		return nil, errors.New("Query parameter `id` is not valid.")
	}
	return id, nil
}

// Join makes this Chord HTTP service attempt to join a Chord ring available
// via a peer node at specified TCP address. Providing an `addr` being `nil`
// causes the service to form its own ring.
func (service *HTTPService) Join(addr *net.TCPAddr) {
	var peer Node
	if addr != nil {
		peer = service.pool.GetOrCreateNode(addr)
	}
	service.pool.lnode.Join(peer)
}

// Refresh causes the HTTP service to refresh its data.
//
// This method should be called at sensible intervals in order for the service
// to maintain its integrity.
func (service *HTTPService) Refresh() error {
	lnode := service.pool.lnode

	if err := lnode.Stabilize(); err != nil {
		return err
	}
	if err := lnode.FixRandomFinger(); err != nil {
		return err
	}
	return nil
}

func (service *HTTPService) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			http.Error(w, fmt.Sprint(r), http.StatusInternalServerError)
		}
	}()
	fmt.Println("Incoming request:", req.Method, req.URL)
	service.router.ServeHTTP(w, req)
}