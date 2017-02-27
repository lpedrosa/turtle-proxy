package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lpedrosa/turtle-proxy/delay"
	"github.com/lpedrosa/turtle-proxy/handlers"
)

const defaultPort int = 6000

type Config struct {
	Host        string
	ProxyPort   int
	ApiPort     int
	ProxyTarget string
}

func parseConfig() (*Config, error) {
	host := flag.String("host", "0.0.0.0", "hostname")
	proxyPort := flag.Int("port", defaultPort, "proxy port")
	apiPort := flag.Int("api-port", (*proxyPort)+1, "api port")
	target := flag.String("target", "", "proxy target")

	flag.Parse()

	if target == nil || *target == "" {
		return nil, errors.New("you must specify a target")
	}

	return &Config{
		Host:        *host,
		ProxyPort:   *proxyPort,
		ApiPort:     *apiPort,
		ProxyTarget: *target}, nil
}

func main() {
	config, err := parseConfig()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}

	storage := delay.DefaultStorage()

	apiHandlers := handlers.NewApiHandlers(storage)
	apiMux := mux.NewRouter()

	apiMux.HandleFunc("/delay", apiHandlers.CreateDelay).Methods("POST")
	apiMux.HandleFunc("/delay", apiHandlers.ClearDelays).Methods("DELETE")

	// for monitoring
	apiMux.HandleFunc("/ping", handlePing).Methods("GET")

	api := newConnector(config.Host, config.ApiPort, apiMux)

	proxyHandlers := handlers.NewProxyHandlers(config.ProxyTarget, storage)

	proxyMux := http.NewServeMux()
	proxyMux.HandleFunc("/", proxyHandlers.ProxyRequest)

	proxy := newConnector(config.Host, config.ProxyPort, proxyMux)

	fmt.Printf("Starting turtle-proxy [port: %d, api: %d]...\n", config.ProxyPort, config.ApiPort)
	// start up connectors
	shutdown := make(chan error)

	go func() {
		shutdown <- api.ListenAndServe()
	}()

	go func() {
		shutdown <- proxy.ListenAndServe()
	}()

	err = <-shutdown
	log.Fatal(err)
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	reqVars := mux.Vars(r)
	fmt.Fprintf(w, "method: %s, id: %s", r.Method, reqVars["id"])
}

func newConnector(host string, port int, handler http.Handler) *http.Server {
	addr := host + ":" + strconv.Itoa(port)
	server := &http.Server{Addr: addr, Handler: handler}
	return server
}
