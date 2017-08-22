package grafana

import (
	"context"
	"fmt"
	"net/url"
	"testing"
)

func TestClient_NewRequest_Authorization(t *testing.T) {
	baseURL, _ := url.Parse("http://localhost/")
	token := "token"
	c := NewClient(baseURL, token, nil)
	r, err := c.NewRequest(context.Background(), "GET", "path", nil)
	if err != nil {
		t.Fatalf("NewRequest returned error: %v", err)
	}

	got := r.Header.Get("Authorization")
	want := fmt.Sprintf("Bearer %s", token)
	if got != want {
		t.Errorf("Authorization header is invalid. Got %s, want %s", got, want)

	}
}

func TestClient_NewRequest_ContentType(t *testing.T) {
	ts := []struct {
		body     interface{}
		expected string
	}{
		{nil, ""},
		{struct{}{}, "application/json"},
	}

	baseURL, _ := url.Parse("http://localhost/")
	c := NewClient(baseURL, "token", nil)
	for _, tt := range ts {
		r, err := c.NewRequest(context.Background(), "GET", "path", tt.body)
		if err != nil {
			t.Fatalf("NewRequest returned error: %v", err)
		}

		got := r.Header.Get("Content-Type")
		if got != tt.expected {
			t.Errorf("Content-Type header is invalid. Got %s, want %s", got, tt.expected)

		}
	}
}
