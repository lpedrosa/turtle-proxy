package handlers

import (
	"fmt"
	"net/http"

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

		fmt.Fprintf(w, "Storing %s", entry)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method %s not supported", method)
	}
}

func parseDelayRequest(r *http.Request) (sd *delayed.DelayedDownload, err error) {
	return &delayed.DelayedDownload{Slug: "lol", URL: nil, Delay: 10}, nil
}
