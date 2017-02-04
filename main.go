package main

import (
	"errors"
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
	http.HandleFunc("/get-delay", handleGetDelay)

	// for monitoring
	http.HandleFunc("/ping", handlePing)

	log.Fatal(http.ListenAndServe(addr, nil))
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

func handleGetDelay(w http.ResponseWriter, r *http.Request) {
	method := r.Method

	switch method {
	case "GET":
		slug, err := parseGetDelayRequest(r)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error: %s", err)
			return
		}

		fmt.Fprintf(w, "Would download %s", slug)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method %s not supported", method)
	}
}

func parseGetDelayRequest(r *http.Request) (slug string, err error) {
	slug = r.FormValue("id")

	if len(slug) == 0 {
		err = errors.New("Delay id not found")
		return "", err
	}

	return slug, nil
}
