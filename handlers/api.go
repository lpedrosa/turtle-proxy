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
	rule := delay.Rule{
		Method:        *d.Method,
		Path:          *d.Target,
		RequestDelay:  d.Delay.Request,
		ResponseDelay: d.Delay.Response}

	a.ruleStorage.Store(rule)
	log.Printf("api: added rule: %#v", rule)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (a *ApiHandlers) ClearDelays(w http.ResponseWriter, r *http.Request) {
	a.ruleStorage.Clear()
	log.Println("api: cleared rules")

	w.WriteHeader(http.StatusNoContent)
}

//-------------------------------
// Error handling
//-------------------------------

type ApiError struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	encErr := json.NewEncoder(w).Encode(&ApiError{Error: err.Error()})
	if encErr != nil {
		log.Printf("api: error marshalling response: %s", encErr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

//-------------------------------
// Delay Parsing
//-------------------------------

type DelayConfig struct {
	Request  int `json:request`
	Response int `json:response`
}

func (dc *DelayConfig) validate() error {
	if dc.Request < 0 {
		return errors.New("request delay should be positive")
	}
	if dc.Response < 0 {
		return errors.New("response delay should be positive")
	}

	return nil
}

type DelayRequest struct {
	Method *string     `json:method`
	Target *string     `json:target`
	Delay  DelayConfig `json:delay`
}

func (dr *DelayRequest) validate() error {
	if dr.Method == nil {
		return errors.New("method is required")
	}
	if dr.Target == nil {
		return errors.New("target is required")
	}

	// check if Target is a valid url
	_, err := url.ParseRequestURI(*dr.Target)
	if err != nil {
		return errors.New("target " + *dr.Target + " is not a valid url")
	}

	err = dr.Delay.validate()
	if err != nil {
		return err
	}

	return nil
}

func ParseDelay(r io.Reader) (d *DelayRequest, err error) {
	jsonDecoder := json.NewDecoder(r)

	// read body
	err = jsonDecoder.Decode(&d)
	if err != nil {
		return nil, err
	}

	err = d.validate()
	if err != nil {
		return nil, err
	}

	return d, nil
}
