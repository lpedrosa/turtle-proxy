package handlers

import (
	"strings"
	"testing"
)

func TestParseDelay(t *testing.T) {
	t.Run("ParseDisplay should work for well formed json", func(t *testing.T) {
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
	t.Run("ParseDisplay should default delay response values to 0 when not present", func(t *testing.T) {
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
	t.Run("ParseDisplay should default delay response values to 0 when not present", func(t *testing.T) {
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
}
