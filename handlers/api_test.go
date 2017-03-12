package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lpedrosa/turtle-proxy/delay"
)

func TestCreateDelayHandler(t *testing.T) {
	api := NewApiHandlers(delay.DefaultStorage())

	t.Run("CreateDelay should reply with status 400 when it fails to parse a rule", func(t *testing.T) {
		badBody := strings.NewReader(`{"hello": "world"}`)
		request := httptest.NewRequest("POST", "/delay", badBody)

		recorder := httptest.NewRecorder()

		api.CreateDelay(recorder, request)

		if recorder.Code != http.StatusBadRequest {
			t.Fatalf("Expected %d, got %d", http.StatusBadRequest, recorder.Code)
		}

		var apiError ApiError
		err := json.NewDecoder(recorder.Body).Decode(&apiError)

		if err != nil {
			t.Fatal("Expected response to be an ApiError structure")
		}

		if len(apiError.Error) == 0 {
			t.Fatalf("Expected an error message, got %s", apiError.Error)
		}
	})

	t.Run("CreateDelay should reply with status 201 when it receives a well formed rule", func(t *testing.T) {
		rule := strings.NewReader(`{"method": "GET", "target": "/members", "delay": {"request": 5000, "response": 1000}}`)
		request := httptest.NewRequest("POST", "/delay", rule)

		recorder := httptest.NewRecorder()

		api.CreateDelay(recorder, request)

		if recorder.Code != http.StatusCreated {
			t.Fatalf("Expected %d, got %d", http.StatusBadRequest, recorder.Code)
		}

		if recorder.Body.Len() != 0 {
			t.Fatal("Expected response to be empty")
		}
	})
}

func TestClearDelaysHandler(t *testing.T) {
	t.Run("ClearDelays should reply with status 204 when it is successful", func(t *testing.T) {
		api := NewApiHandlers(delay.DefaultStorage())
		request := httptest.NewRequest("DELETE", "/delay", nil)

		recorder := httptest.NewRecorder()

		api.ClearDelays(recorder, request)

		if recorder.Code != http.StatusNoContent {
			t.Fatalf("Expected %d, got %d", http.StatusNoContent, recorder.Code)
		}

		if recorder.Body.Len() != 0 {
			t.Fatal("Expected response to be empty")
		}
	})
}

func TestParseDelay(t *testing.T) {
	t.Run("ParseDelay should work for well formed json", func(t *testing.T) {
		r := strings.NewReader(`{"method": "GET", "target": "/members", "delay": {"request": 5000, "response": 1000}}`)

		d, err := ParseDelay(r)

		if err != nil {
			t.Fatalf("Expected no error, got %s", err)
		}

		if *d.Method != "GET" {
			t.Fatalf("Expected GET, got: %s", *d.Method)
		}

		if *d.Target != "/members" {
			t.Fatalf("Expected /members, got: %s", *d.Target)
		}

		if d.Delay.Request != 5000 {
			t.Fatalf("Expected 5000, got: %d", d.Delay.Request)
		}

		if d.Delay.Response != 1000 {
			t.Fatalf("Expected 1000, got: %d", d.Delay.Response)
		}
	})
	t.Run("ParseDelay should default delay response values to 0 when not present", func(t *testing.T) {
		r := strings.NewReader(`{"method": "GET", "target": "/members", "delay": {"request": 5000}}`)

		d, err := ParseDelay(r)

		if err != nil {
			t.Fatalf("Expected no error, got %s", err)
		}

		if d.Delay.Request != 5000 {
			t.Fatalf("Expected 5000, got: %d", d.Delay.Request)
		}

		if d.Delay.Response != 0 {
			t.Fatalf("Expected 0, got: %d", d.Delay.Response)
		}
	})
	t.Run("ParseDelay should default delay response values to 0 when not present", func(t *testing.T) {
		r := strings.NewReader(`{"method": "GET", "target": "/members", "delay": {"response": 5000}}`)

		d, err := ParseDelay(r)

		if err != nil {
			t.Fatalf("Expected no error, got %s", err)
		}

		if d.Delay.Request != 0 {
			t.Fatalf("Expected 0, got: %d", d.Delay.Request)
		}

		if d.Delay.Response != 5000 {
			t.Fatalf("Expected 5000, got: %d", d.Delay.Response)
		}
	})
	t.Run("ParseDelay should fail if delay request or response fields are negative", func(t *testing.T) {
		r := strings.NewReader(`{"method": "GET", "target": "/members", "delay": {"response": -100}}`)

		_, err := ParseDelay(r)

		if err == nil {
			t.Fatal("Expected an error for a negative response value")
		}

		r = strings.NewReader(`{"method": "GET", "target": "/members", "delay": {"request": -100}}`)

		_, err = ParseDelay(r)

		if err == nil {
			t.Fatal("Expected an error for a negative request value")
		}
	})
	t.Run("ParseDelay should fail if target is not a valid url", func(t *testing.T) {
		r := strings.NewReader(`{"method": "GET", "target": "la\\la"`)

		_, err := ParseDelay(r)

		if err == nil {
			t.Fatal("Expected an error for a invalid url value")
		}
	})
}
