package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/lpedrosa/turtle-proxy/handlers"
)

const defaultPort int = 6000

func main() {
	fmt.Printf("Starting turtle-proxy on port %d...\n", defaultPort)

	apiMux := http.NewServeMux()

	//adminMux.HandleFunc("/delay", handlers.HandleRegisterDelayed)
	//adminMux.HandleFunc("/delay/", handlers.HandleGetDelayed)

	// for monitoring
	apiMux.HandleFunc("/ping", handlePing)

	api := newConnector(defaultPort+1, adminMux)

	proxyHandlers := handlers.NewProxyHandlers("localhost", 8000)

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

func newConnector(port int, handler http.Handler) *http.Server {
	addr := ":" + strconv.Itoa(port)
	server := &http.Server{Addr: addr, Handler: handler}
	return server
}
