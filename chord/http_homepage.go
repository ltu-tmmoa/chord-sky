package chord

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ltu-tmmoa/chord-sky/log"
	"html/template"
	"net/http"
	"path/filepath"
	"runtime/debug"
	  "os"
)

type HTTPHomepage struct {
	router *mux.Router
}

func NewHTTPHomepage() *HTTPHomepage {
	service := HTTPHomepage{
		router: mux.NewRouter(),
	}
	router := service.router
	router.
		HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {

			if req.Body != nil {
				defer req.Body.Close()
			}
		  	goPath := os.Getenv("GOPATH")
			absPath, _ := filepath.Abs(goPath+"/src/github.com/ltu-tmmoa/chord-sky/template/index.html")
			t, _ := template.ParseFiles(absPath)
			t.Execute(w, nil)

		}).
		Methods(http.MethodGet)

	return &service
}

func (service *HTTPHomepage) ServeHTTP(w http.ResponseWriter, req *http.Request) {
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