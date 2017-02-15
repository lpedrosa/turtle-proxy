package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/lpedrosa/turtle-proxy/delayed"
)

func HandleRegisterDelayed(w http.ResponseWriter, r *http.Request) {
	method := r.Method

	switch method {
	case "POST":
		entry, err := parseDelayRequest(r)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error: %s", err)
			return
		}

		responseContent := map[string]string{"id": entry.Slug}

		json.NewEncoder(w).Encode(responseContent)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method %s not supported", method)
	}
}

func parseDelayRequest(r *http.Request) (sd *delayed.DelayedDownload, err error) {
	jsonDecoder := json.NewDecoder(r.Body)

	var parsedReq struct {
		Target string
		Delay  uint
	}
	// read body as json
	err = jsonDecoder.Decode(&parsedReq)
	if err != nil {
		return nil, err
	}

	// check if target is a valid url
	targetURL, err := url.ParseRequestURI(parsedReq.Target)
	if err != nil {
		return nil, errors.New("target is not a valid url")
	}

	//jsonEncoder.Encode(
	return &delayed.DelayedDownload{
		Slug:  "lol",
		URL:   targetURL,
		Delay: parsedReq.Delay}, nil
}

func HandleGetDelayed(w http.ResponseWriter, r *http.Request) {
	method := r.Method

	switch method {
	case "GET":
		slug, err := parseGetDelayedRequest(r)

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

func parseGetDelayedRequest(r *http.Request) (slug string, err error) {
	slug = r.FormValue("id")

	if len(slug) == 0 {
		err = errors.New("Delay id not found")
		return "", err
	}

	return slug, nil
}
