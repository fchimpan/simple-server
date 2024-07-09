package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fchimpan/simple-server/config"
)

func TestNewMux(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/health", nil)
	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}
	m, _, err := NewMux(context.Background(), cfg)
	if err != nil {
		t.Fatalf("failed to create mux: %v", err)
	}
	m.ServeHTTP(w, r)
	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, resp.StatusCode)
	}
	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("failed to read response body: %v", err)
	}
	want := `{"status":"ok"}`
	if string(got) != want {
		t.Errorf("want %q, got %q", want, got)
	}
}
