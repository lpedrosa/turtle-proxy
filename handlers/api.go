package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"

	"github.com/lpedrosa/turtle-proxy/delay"
)

type ApiHandlers struct {
	ruleStorage *delay.RuleStorage
}

func NewApiHandlers(storage *delay.RuleStorage) *ApiHandlers {
	return &ApiHandlers{ruleStorage: storage}
}

func (a *ApiHandlers) CreateDelay(w http.ResponseWriter, r *http.Request) {
	d, err := parseDelay(r)
	if err != nil {
		log.Printf("Error parsing json: %s", err)
		writeError(w, err)
		return
	}

	// store delay
	a.ruleStorage.Store(d.target, d.delay)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

type apiError struct {
	err string
}

func writeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	encErr := json.NewEncoder(w).Encode(&apiError{err: err.Error()})
	if encErr != nil {
		log.Printf("Error marshalling response: %s", encErr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type delayRequest struct {
	method string
	target string
	delay  int
}

func parseDelay(r *http.Request) (d *delayRequest, err error) {
	jsonDecoder := json.NewDecoder(r.Body)

	// read body
	err = jsonDecoder.Decode(&d)
	if err != nil {
		return nil, err
	}

	log.Printf("Got: %#v", d)

	// check if target is a valid url
	_, err = url.ParseRequestURI(d.target)
	if err != nil {
		return nil, errors.New("target is not a valid url")
	}

	if d.delay <= 0 {
		return nil, errors.New("delay should be > 0")
	}

	return d, nil
}
