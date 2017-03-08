package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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

func NewProxyHandlers(hostPort string, storage *delay.RuleStorage) *ProxyHandlers {
	return &ProxyHandlers{
		target:      "http://" + hostPort,
		client:      &http.Client{},
		ruleStorage: storage}
}

func (p *ProxyHandlers) ProxyRequest(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	method := r.Method

	delay := p.checkForDelay(method, path)

	targetPath := p.buildTargetPath(r.URL)

	msg := "proxy: incoming %s request to: %s. Delaying it [pre: %dms, post: %dms]"
	log.Printf(msg, method, targetPath, delay.request, delay.response)

	// create request
	req, err := http.NewRequest(method, targetPath, r.Body)
	if err != nil {
		errorMsg := fmt.Sprintf("error while proxying: %s", err)
		log.Printf("proxy: %s\n", errorMsg)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, errorMsg)
		return
	}

	// copy all request headers
	req.Header = r.Header

	// delay the request
	turtleIt(delay.request)

	// reply original client
	res, err := p.client.Do(req)
	if err != nil {
		errorMsg := fmt.Sprintf("error while proxying: %s", err)
		log.Printf("proxy: %s\n", errorMsg)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, errorMsg)
		return
	}

	// delay the response
	turtleIt(delay.response)

	// copy all the response headers
	for k, v := range res.Header {
		w.Header().Set(k, strings.Join(v, ","))
	}

	w.WriteHeader(res.StatusCode)

	//pipe body to writer
	defer res.Body.Close()
	_, err = io.Copy(w, res.Body)
	if err != nil {
		log.Printf("proxy: error while copying response body: %s", err)
		return
	}
}

func (p *ProxyHandlers) buildTargetPath(url *url.URL) string {
	targetPath := p.target + url.Path
	query := url.RawQuery

	if query == "" {
		return targetPath
	}

	return targetPath + "?" + query
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
