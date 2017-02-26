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

	adminMux := http.NewServeMux()

	adminMux.HandleFunc("/delay", handlers.HandleRegisterDelayed)
	adminMux.HandleFunc("/delay/", handlers.HandleGetDelayed)

	// for monitoring
	adminMux.HandleFunc("/ping", handlePing)

	admin := newConnector(defaultPort+1, adminMux)

	proxyMux := http.NewServeMux()

	proxyMux.HandleFunc("/", handleHello)

	proxy := newConnector(defaultPort, proxyMux)

	// start up connectors
	shutdown := make(chan error)

	go func() {
		shutdown <- admin.ListenAndServe()
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

func handleHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world")
}

func newConnector(port int, handler http.Handler) *http.Server {
	addr := ":" + strconv.Itoa(port)
	server := &http.Server{Addr: addr, Handler: handler}
	return server
}
