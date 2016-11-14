package chord

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ltu-tmmoa/chord-sky/data"
	"github.com/ltu-tmmoa/chord-sky/log"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"runtime/debug"
	  "os"
)

// HTTPStorageService manages local storage, exposing it as an HTTP service by
// implementing the data.storage interface.
type HTTPStorageService struct {
	storage *data.MemoryStorage
	router  *mux.Router
}

// HTTPStorageService creates a new HTTP storage, exposable as a service on the
// identified local TCP interface.
func NewHTTPStorageService() *HTTPStorageService {
	service := HTTPStorageService{
		storage: data.NewMemoryStorage(),
		router:  mux.NewRouter(),
	}

	storage := service.storage
	router := service.router

	router.
		HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {

			if req.Body != nil {
				defer req.Body.Close()
			}
			req.ParseForm()

			strID := req.Form["key"][0]
			strValue := req.Form["value"][0]

			fmt.Printf("Key: %s, Value: %s", strID, strValue)
			id, ok := parseID(strID)
			if !ok {
				err := errors.New("file `id` is not valid.")
				httpWrite(w, http.StatusBadRequest, err.Error())
				return
			} else {
				  goPath := os.Getenv("GOPATH")
				  absPath, _ := filepath.Abs(goPath+"/src/github.com/ltu-tmmoa/chord-sky/template/index.html")
				t, _ := template.ParseFiles(absPath)
				t.Execute(w, nil)
			}
			arr := []byte(strValue)
			storage.Set(id, arr)

		}).
		Methods(http.MethodPost)

	router.
		HandleFunc("/keys", func(w http.ResponseWriter, req *http.Request) {

			if req.Body != nil {
				defer req.Body.Close()
			}
			strfromKey := req.URL.Query().Get("from")
			strtoKey := req.URL.Query().Get("to")
			// If provided a key in the query
			if len(strfromKey) > 0 || len(strtoKey) > 0 {
				fromKey, ok1 := parseID(strfromKey)
				toKey, ok2 := parseID(strtoKey)
				if !ok1 || !ok2 {
					strErr := fmt.Sprintf("The `from` is %v and `to` is %v", ok1, ok2)
					err := errors.New(strErr)
					httpWrite(w, http.StatusBadRequest, err.Error())
					return
				}
				keys, err := storage.GetKeyRange(fromKey, toKey)
				if err != nil {
					httpWrite(w, http.StatusInternalServerError, err.Error())
					return
				}
				var buffer bytes.Buffer
				for _, v := range keys {
					buffer.WriteString(v.String())
					buffer.WriteString("\n")
				}
				httpStorageWrite(w, http.StatusOK, buffer.String())
			} else { // else send all the local keys
				keys, err := storage.GetAllKeys()
				if err != nil {
					httpWrite(w, http.StatusInternalServerError, err.Error())
					return
				}
				var buffer bytes.Buffer
				for _, v := range keys {
					buffer.WriteString(v.String())
					buffer.WriteString("\n")
				}
				httpStorageWrite(w, http.StatusOK, buffer.String())
			}

		}).
		Methods(http.MethodGet)

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
			httpStorageWrite(w, http.StatusOK, string(value))

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
			httpStorageWrite(w, http.StatusOK, "")

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
