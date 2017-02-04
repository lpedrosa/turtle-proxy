package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleRegisterDelayed(t *testing.T) {
	t.Run("Supports only POST", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/delayed", nil)

		resp := httptest.NewRecorder()

		HandleRegisterDelayed(resp, req)

		if status := resp.Code; status != http.StatusMethodNotAllowed {
			t.Errorf("expected status code %v, got %v",
				status, http.StatusMethodNotAllowed)
		}
	})
	t.Run("Expects only JSON", func(t *testing.T) {
	})
}
