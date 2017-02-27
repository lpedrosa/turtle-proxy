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
	a.ruleStorage.Store(*d.Target, d.Delay.Request)

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

type DelayConfig struct {
	Request  int
	Response int
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
	Method *string
	Target *string
	Delay  DelayConfig
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
