package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAliveHandler_OK(t *testing.T) {
	svc := &mockHealthService{aliveErr: nil}
	handler := AliveHandler(svc)

	req := httptest.NewRequest("GET", "/health/alive", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if body["status"] != "alive" {
		t.Errorf("expected status 'alive', got %q", body["status"])
	}
}

func TestAliveHandler_Unavailable(t *testing.T) {
	svc := &mockHealthService{aliveErr: fmt.Errorf("dead")}
	handler := AliveHandler(svc)

	req := httptest.NewRequest("GET", "/health/alive", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", rec.Code)
	}
}

func TestReadyHandler_OK(t *testing.T) {
	svc := &mockHealthService{readyErr: nil}
	handler := ReadyHandler(svc)

	req := httptest.NewRequest("GET", "/health/ready", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if body["status"] != "ready" {
		t.Errorf("expected status 'ready', got %q", body["status"])
	}
}

func TestReadyHandler_Unavailable(t *testing.T) {
	svc := &mockHealthService{readyErr: fmt.Errorf("db down")}
	handler := ReadyHandler(svc)

	req := httptest.NewRequest("GET", "/health/ready", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", rec.Code)
	}
}
