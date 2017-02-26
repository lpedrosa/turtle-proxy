package handlers

/*
import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleRegisterDelayed(t *testing.T) {
	t.Run("Supports only POST", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/delay", nil)

		resp := httptest.NewRecorder()

		HandleRegisterDelayed(resp, req)

		if status := resp.Code; status != http.StatusMethodNotAllowed {
			t.Errorf("expected status code %v, got %v",
				http.StatusMethodNotAllowed, status)
		}
	})
	t.Run("Expects only JSON", func(t *testing.T) {
		body := strings.NewReader("some crap")

		req := httptest.NewRequest("POST", "/delay", body)

		resp := httptest.NewRecorder()

		HandleRegisterDelayed(resp, req)

		if status := resp.Code; status != http.StatusBadRequest {
			t.Errorf("expected status code %v, got %v",
				http.StatusBadRequest, status)
		}
	})
	t.Run("Handles request successfully", func(t *testing.T) {
		body := strings.NewReader(`{"target": "http://example.com", "delay": 1}`)

		req := httptest.NewRequest("POST", "/delay", body)

		resp := httptest.NewRecorder()

		HandleRegisterDelayed(resp, req)

		if status := resp.Code; status != http.StatusOK {
			t.Errorf("expected status code %v, got %v",
				http.StatusOK, status)
		}

		var decodedResp map[string]string
		json.NewDecoder(resp.Body).Decode(&decodedResp)

		if _, ok := decodedResp["id"]; !ok {
			t.Errorf("expected slug to be present in the response, but go %v", decodedResp)
		}
	})
	t.Run("Rejects invalid target urls", func(t *testing.T) {
		body := strings.NewReader(`{"target": "nope", "delay": 1}`)

		req := httptest.NewRequest("POST", "/delay", body)

		resp := httptest.NewRecorder()

		HandleRegisterDelayed(resp, req)

		if status := resp.Code; status != http.StatusBadRequest {
			t.Errorf("expected status code %v, got %v",
				http.StatusBadRequest, status)
		}
	})
}

func TestHandleGetDelay(t *testing.T) {
	t.Run("Supports only GET", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/delay/1", nil)

		resp := httptest.NewRecorder()

		HandleGetDelayed(resp, req)

		if status := resp.Code; status != http.StatusMethodNotAllowed {
			t.Errorf("expected status code %v, got %v",
				http.StatusMethodNotAllowed, status)
		}
	})
	t.Run("Requires a delay id on the path", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/delay/", nil)

		resp := httptest.NewRecorder()

		HandleGetDelayed(resp, req)

		if status := resp.Code; status != http.StatusBadRequest {
			t.Errorf("expected status code %v, got %v",
				http.StatusBadRequest, status)
		}
	})
}

func TestParseGetDelayedRequest(t *testing.T) {
	expectedSlug := "1"
	req := httptest.NewRequest("GET", "/delay/"+expectedSlug, nil)

	slug, err := ParseGetDelayedRequest(req)

	if err != nil {
		t.Fatalf("expected slug to be %v, got error: %v", expectedSlug, err)
	}

	if slug != expectedSlug {
		t.Fatalf("expected slug to be %v, got %v", expectedSlug, slug)
	}
}
*/
