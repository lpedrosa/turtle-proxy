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
	fmt.Printf("Starting turtle-proxy on port %d...", defaultPort)

	addr := ":" + strconv.Itoa(defaultPort)

	http.HandleFunc("/delay", handlers.HandleRegisterDelayed)
	http.HandleFunc("/delay/", handlers.HandleGetDelayed)

	// for monitoring
	http.HandleFunc("/ping", handlePing)

	log.Fatal(http.ListenAndServe(addr, nil))
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}
