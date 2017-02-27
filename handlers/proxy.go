package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/lpedrosa/turtle-proxy/delay"
)

var zeroDelay *delayConfig = new(delayConfig)

type ProxyHandlers struct {
	target      string
	client      *http.Client
	ruleStorage *delay.RuleStorage
}

func NewProxyHandlers(host string, port int, storage *delay.RuleStorage) *ProxyHandlers {
	return &ProxyHandlers{
		target:      "http://" + host + ":" + strconv.Itoa(port),
		client:      &http.Client{},
		ruleStorage: storage}
}

func (p *ProxyHandlers) ProxyRequest(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	method := r.Method

	log.Printf("Incoming %s request to: %s", method, path)

	// create request
	req, err := http.NewRequest(method, p.target+path, r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error while proxying: %s", err)
		return
	}

	// delay it
	delay := p.checkForDelay(method, path)
	turtleIt(delay.request)

	// reply original client
	res, err := p.client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error while proxying: %s", err)
		return
	}

	// copy all the headers
	for k, v := range res.Header {
		w.Header().Set(k, strings.Join(v, ","))
	}

	w.WriteHeader(res.StatusCode)

	//pipe body to writer
	defer res.Body.Close()
	_, err = io.Copy(w, res.Body)
	if err != nil {
		log.Printf("Error found while copying response body: %s", err)
		return
	}
}

type delayConfig struct {
	request  int
	response int
}

func (p *ProxyHandlers) checkForDelay(method string, path string) *delayConfig {
	rule, ok := p.ruleStorage.Get(method, path)
	if !ok {
		// no delay found
		return zeroDelay
	}

	return &delayConfig{request: rule.RequestDelay, response: rule.ResponseDelay}
}

func turtleIt(delay int) {
	time.Sleep(time.Duration(delay) * time.Millisecond)
}
