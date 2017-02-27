package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lpedrosa/turtle-proxy/delay"
	"github.com/lpedrosa/turtle-proxy/handlers"
)

const defaultPort int = 6000

func main() {
	fmt.Printf("Starting turtle-proxy on port %d...\n", defaultPort)

	storage := delay.DefaultStorage()

	apiHandlers := handlers.NewApiHandlers(storage)
	apiMux := mux.NewRouter()

	apiMux.HandleFunc("/delay", apiHandlers.CreateDelay).Methods("POST")
	apiMux.HandleFunc("/delay", apiHandlers.ClearDelays).Methods("DELETE")

	// for monitoring
	apiMux.HandleFunc("/ping", handlePing).Methods("GET")

	api := newConnector(defaultPort+1, apiMux)

	proxyHandlers := handlers.NewProxyHandlers("localhost", 8000, storage)

	proxyMux := http.NewServeMux()
	proxyMux.HandleFunc("/", proxyHandlers.ProxyRequest)

	proxy := newConnector(defaultPort, proxyMux)

	// start up connectors
	shutdown := make(chan error)

	go func() {
		shutdown <- api.ListenAndServe()
	}()

	go func() {
		shutdown <- proxy.ListenAndServe()
	}()

	err := <-shutdown
	log.Fatal(err)
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	reqVars := mux.Vars(r)
	fmt.Fprintf(w, "method: %s, id: %s", r.Method, reqVars["id"])
}

func newConnector(port int, handler http.Handler) *http.Server {
	addr := ":" + strconv.Itoa(port)
	server := &http.Server{Addr: addr, Handler: handler}
	return server
}
