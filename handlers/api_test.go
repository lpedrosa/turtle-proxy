package handlers

import (
	"strings"
	"testing"
)

func TestParseDelay(t *testing.T) {

	r := strings.NewReader(`{"method": "GET", "target": "/members", "delay": 5000}`)

	d, err := ParseDelay(r)

	if err != nil {
		t.Fatalf("Expected no error, got %s", err)
	}

	if d.Method != "GET" {
		t.Fatalf("Expected GET, got: %s", d.Method)
	}

	if d.Target != "/members" {
		t.Fatalf("Expected /members, got: %s", d.Target)
	}

	if d.Delay != 5000 {
		t.Fatalf("Expected 5000, got: %d", d.Delay)
	}
}
