package chord

import (
	  "errors"
	  "github.com/gorilla/mux"
	  "github.com/ltu-tmmoa/chord-sky/data"
	  "github.com/ltu-tmmoa/chord-sky/log"

	  "net/http"
	  "net"
	  "fmt"
	  "io/ioutil"
	  "runtime/debug"
)

// HTTPStorageService manages local storage, exposing it as an HTTP service by
// implementing the data.storage interface.
type HTTPStorageService struct {
	  storage  *data.MemoryStorage
	  router   *mux.Router
}

// HTTPStorageService creates a new HTTP storage, exposable as a service on the
// identified local TCP interface.
func NewHTTPStorageService(laddr *net.TCPAddr) *HTTPStorageService {
	  service := HTTPStorageService{
		    storage:	data.NewMemoryStorage(),
		    router: 	mux.NewRouter(),
	  }

	  storage := service.storage
	  router := service.router


	  router.
	  HandleFunc("/{id}", func(w http.ResponseWriter, req *http.Request) {

		    if req.Body != nil {
				req.Body.Close()
		    }
		    strID, _ := mux.Vars(req)["id"]
		    id, ok := parseID(strID)
		    if !ok {
				err := errors.New("file `id` is not valid.")
				httpWrite(w, http.StatusBadRequest, err.Error())
				return
		    }
		    value, _ := storage.Get(id)
		    httpStorageWrite(w, http.StatusOK, value)


	  }).
		    Methods(http.MethodGet)
	  router.
	  HandleFunc("/{id}", func(w http.ResponseWriter, req *http.Request) {

		    if req.Body != nil {
				defer req.Body.Close()
		    }
		    strID, _ := mux.Vars(req)["id"]
		    id, ok := parseID(strID)
		    if !ok {
				err := errors.New("file `id` is not valid.")
				httpWrite(w, http.StatusBadRequest, err.Error())
				return
		    }
		    arr, err := ioutil.ReadAll(req.Body)
		    if err != nil {
				err := errors.New("file `id` is not valid.")
				httpWrite(w, http.StatusBadRequest, err.Error())
				return
		    }
		    storage.Set(id, arr)
		    httpStorageWrite(w, http.StatusOK, nil)

	  }).
		    Methods(http.MethodPut)
	  router.
	  HandleFunc("/{id}", func(w http.ResponseWriter, req *http.Request) {

		    if req.Body != nil {
				defer req.Body.Close()
		    }
		    strID, _ := mux.Vars(req)["id"]
		    id, ok := parseID(strID)
		    if !ok {
				err := errors.New("file `id` is not valid.")
				httpWrite(w, http.StatusBadRequest, err.Error())
				return
		    }
		    storage.Remove(id)
		    httpStorageWrite(w, http.StatusOK, nil)

	  }).
		    Methods(http.MethodDelete)

	  router.
	  HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {

		    if req.Body != nil {
				defer req.Body.Close()
		    }
		    strfromKey, _ := mux.Vars(req)["fromKey"]
		    strtoKey, _ := mux.Vars(req)["toKey"]


		    fromKey, ok1 := parseID(strfromKey)
		    toKey, ok2 := parseID(strtoKey)
		    if !ok1 || !ok2 {
				err := errors.New("file `id` is not valid.")
				httpWrite(w, http.StatusBadRequest, err.Error())
				return
		    }
		    keys, err := storage.GetKeyRange(fromKey, toKey)
		    if err!=nil{
				httpWrite(w, http.StatusInternalServerError, err.Error())
				return
		    }
		    httpStorageWrite(w, http.StatusOK, keys)

	  }).
		    Methods(http.MethodDelete)

	  return &service
}

func httpStorageWrite(w http.ResponseWriter, status int, body interface{}) {
	  w.WriteHeader(status)
	  fmt.Fprint(w, body)
}

func (service *HTTPStorageService) ServeHTTP(w http.ResponseWriter, req *http.Request) {
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