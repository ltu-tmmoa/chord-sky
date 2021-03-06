package chord

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/ltu-tmmoa/chord-sky/data"
	"github.com/ltu-tmmoa/chord-sky/log"
)

var (
	errBodyMissing = errors.New("No body in request. Required.")
)

// HTTPService manages a local Chord node, exposing it as an HTTP service by
// implementing the http.Handler interface.
type HTTPService struct {
	pool     *nodePool
	router   *mux.Router
	isJoined bool
}

// NewHTTPService creates a new HTTP node, exposable as a service on the
// identified local TCP interface.
func NewHTTPService(laddr *net.TCPAddr) *HTTPService {
	service := HTTPService{
		pool:   newNodePool(laddr),
		router: mux.NewRouter(),
	}

	pool := service.pool
	router := service.router
	lnode := pool.lnode

	router.
		HandleFunc("/info", func(w http.ResponseWriter, req *http.Request) {
			var pred string
			{
				predNode, predErr := lnode.Predecessor()
				if predErr != nil {
					pred = predErr.Error()
				} else {
					pred = predNode.String()
				}
			}
			buf := &bytes.Buffer{}
			fmt.Fprintf(buf, "ID:          %s\r\n", lnode.ID())
			fmt.Fprintf(buf, "Successor:   %s\r\n", lnode.successor())
			fmt.Fprintf(buf, "Predecessor: %s\r\n", pred)

			fmt.Fprint(buf, "\r\nSuccessor List:\r\n")
			for i, succ := range lnode.succlist {
				fmt.Fprintf(buf, "%3d:         %s\r\n", i, succ)
			}
			fmt.Fprint(buf, "\r\nFinger Table:\r\n")
			m := lnode.ID().Bits()
			for i := 1; i <= m; i++ {
				fmt.Fprintf(buf, "%3d:         %s\r\n", i, lnode.fingerNode(i))
			}
			w.WriteHeader(http.StatusOK)
			w.Write(buf.Bytes())
		}).
		Methods(http.MethodGet)

	router.
		HandleFunc("/info/ring", func(w http.ResponseWriter, req *http.Request) {
			if req.Body != nil {
				req.Body.Close()
			}
			w.WriteHeader(http.StatusOK)
			lnode.writeRingTextTo(w)
		}).
		Methods(http.MethodGet)

	router.
		HandleFunc("/info/fix", func(w http.ResponseWriter, req *http.Request) {
			if req.Body != nil {
				req.Body.Close()
			}
			if err := lnode.fixAllFingers(); err != nil {
				panic(err)
			}
			w.WriteHeader(http.StatusNoContent)
		}).
		Methods(http.MethodGet)

	router.
		HandleFunc("/fingers/{i:[0-9]+}", func(w http.ResponseWriter, req *http.Request) {
			if req.Body != nil {
				req.Body.Close()
			}
			i, _ := strconv.Atoi(mux.Vars(req)["i"])
			node := lnode.fingerNode(i)
			httpWrite(w, http.StatusOK, node.TCPAddr())
		}).
		Methods(http.MethodGet)

	router.
		HandleFunc("/fingers/{i:[0-9]+}", func(w http.ResponseWriter, req *http.Request) {
			i, _ := strconv.Atoi(mux.Vars(req)["i"])
			addr, err := httpReadBodyAsAddr(req)
			if err != nil {
				httpWrite(w, http.StatusBadRequest, err.Error())
				return
			}
			node := pool.getOrCreateNode(addr)
			lnode.SetFingerNode(i, node)
			w.WriteHeader(http.StatusNoContent)
		}).
		Methods(http.MethodPut)

	router.
		HandleFunc("/heartbeat", func(w http.ResponseWriter, req *http.Request) {
			if req.Body != nil {
				req.Body.Close()
			}
			httpWrite(w, http.StatusOK, "\u2764")
		}).
		Methods(http.MethodGet)

	router.
		HandleFunc("/successor", func(w http.ResponseWriter, req *http.Request) {
			if req.Body != nil {
				req.Body.Close()
			}
			succ := lnode.successor()
			httpWrite(w, http.StatusOK, succ.TCPAddr())
		}).
		Methods(http.MethodGet)

	router.
		HandleFunc("/predecessor", func(w http.ResponseWriter, req *http.Request) {
			if req.Body != nil {
				req.Body.Close()
			}
			pred, err := lnode.Predecessor()
			if err != nil {
				httpWrite(w, http.StatusFailedDependency, err.Error())
				return
			}
			httpWrite(w, http.StatusOK, pred.TCPAddr())
		}).
		Methods(http.MethodGet)

	router.
		HandleFunc("/successors", func(w http.ResponseWriter, req *http.Request) {
			if req.Body != nil {
				req.Body.Close()
			}
			var id *data.ID
			{
				var err error
				id, err = httpReadQueryID(req)
				if err != nil {
					httpWrite(w, http.StatusBadRequest, err.Error())
					return
				}
			}
			if id == nil {
				succs, _ := lnode.SuccessorList()
				buf := &bytes.Buffer{}
				for _, succ := range succs {
					fmt.Fprintf(buf, "%s\r\n", succ.TCPAddr())
				}
				httpWrite(w, http.StatusOK, string(buf.Bytes()))
				return
			}
			node, err := lnode.FindSuccessor(id)
			if err != nil {
				httpWrite(w, http.StatusFailedDependency, err.Error())
				return
			}
			httpWrite(w, http.StatusOK, node.TCPAddr())
		}).
		Methods(http.MethodGet)

	router.
		HandleFunc("/predecessors", func(w http.ResponseWriter, req *http.Request) {
			id, err := httpReadQueryID(req)
			if err != nil {
				httpWrite(w, http.StatusBadRequest, err.Error())
				return
			}
			if id == nil {
				httpWrite(w, http.StatusBadRequest, "Query parameter `id` required.")
				return
			}
			node, err := lnode.FindPredecessor(id)
			if err != nil {
				httpWrite(w, http.StatusFailedDependency, err.Error())
				return
			}
			httpWrite(w, http.StatusOK, node.TCPAddr())
		}).
		Methods(http.MethodGet)

	router.
		HandleFunc("/successor", func(w http.ResponseWriter, req *http.Request) {
			addr, err := httpReadBodyAsAddr(req)
			if err != nil {
				httpWrite(w, http.StatusBadRequest, err.Error())
				return
			}
			succ := pool.getOrCreateNode(addr)
			if err = lnode.SetSuccessor(succ); err != nil {
				httpWrite(w, http.StatusFailedDependency, err.Error())
				return
			}
			w.WriteHeader(http.StatusNoContent)
		}).
		Methods(http.MethodPut)

	router.
		HandleFunc("/predecessor", func(w http.ResponseWriter, req *http.Request) {
			addr, err := httpReadBodyAsAddr(req)
			if err != nil {
				httpWrite(w, http.StatusBadRequest, err.Error())
				return
			}
			pred := pool.getOrCreateNode(addr)
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
	if len(body) == 0 {
		return nil, errors.New("Cannot resolve empty body into TCP address.")
	}
	return net.ResolveTCPAddr("tcp", body)
}

func httpReadBodyAsAddrs(req *http.Request) ([]*net.TCPAddr, error) {
	body, err := httpReadBody(req)
	if err != nil {
		return nil, err
	}
	tokens := strings.Split(body, "\r\n")
	addrs := make([]*net.TCPAddr, 0, len(tokens))
	for _, token := range tokens {
		if len(token) == 0 {
			continue
		}
		addr, err := net.ResolveTCPAddr("tcp", token)
		if err != nil {
			return nil, err
		}
		if addr != nil {
			addrs = append(addrs, addr)
		}
	}
	return addrs, nil
}

func httpReadQueryID(req *http.Request) (*data.ID, error) {
	strID := req.URL.Query().Get("id")
	if len(strID) == 0 {
		return nil, nil
	}
	id, ok := parseID(strID)
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
		peer = service.pool.getOrCreateNode(addr)
	}
	service.pool.lnode.join(peer)
}

// Refresh causes the HTTP service to refresh its data.
//
// This method should be called at sensible intervals in order for the service
// to maintain its integrity.
func (service *HTTPService) Refresh() error {
	return service.pool.refresh()
}

func (service *HTTPService) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			http.Error(w, fmt.Sprint(r), http.StatusInternalServerError)
			log.Logger.Println("Recovered:", r)
			log.Logger.Println(string(debug.Stack()))
		}
	}()
	log.Logger.Println(req.Method, req.URL)
	service.router.ServeHTTP(w, req)
}
