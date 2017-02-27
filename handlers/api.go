package handlers

import (
	"encoding/json"
	"errors"
	"io"
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
	d, err := ParseDelay(r.Body)
	if err != nil {
		writeError(w, err)
		return
	}

	// store delay
	a.ruleStorage.Store(d.Target, d.Delay)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

type apiError struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	encErr := json.NewEncoder(w).Encode(&apiError{Error: err.Error()})
	if encErr != nil {
		log.Printf("Error marshalling response: %s", encErr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type DelayRequest struct {
	Method string
	Target string
	Delay  int
}

func ParseDelay(r io.Reader) (d *DelayRequest, err error) {
	jsonDecoder := json.NewDecoder(r)

	// read body
	err = jsonDecoder.Decode(&d)
	if err != nil {
		return nil, err
	}

	// check if Target is a valid url
	_, err = url.ParseRequestURI(d.Target)
	if err != nil {
		return nil, errors.New("target '" + d.Target + "' is not a valid url")
	}

	if d.Delay <= 0 {
		return nil, errors.New("delay should be positive")
	}

	return d, nil
}
