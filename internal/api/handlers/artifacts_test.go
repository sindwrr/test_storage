package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sindwrr/test_storage/internal/models"
)

func TestListArtifacts_Success(t *testing.T) {
	expected := []models.ArtifactInfo{
		{ID: 1, FileName: "log1.txt", Component: "core"},
		{ID: 2, FileName: "log2.txt", Component: "ui"},
	}
	svc := &mockMetadataService{
		getArtifactInfoFn: func(component, build, suite string, fromTime, toTime time.Time) ([]models.ArtifactInfo, error) {
			return expected, nil
		},
	}
	h := NewArtifactsHandler(svc)

	req := httptest.NewRequest("GET", "/artifacts?component=core", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "admin"})
	rec := httptest.NewRecorder()
	h.List(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	if rec.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", rec.Header().Get("Content-Type"))
	}

	var artifacts []models.ArtifactInfo
	if err := json.NewDecoder(rec.Body).Decode(&artifacts); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if len(artifacts) != 2 || artifacts[0].FileName != "log1.txt" {
		t.Errorf("unexpected artifacts: %+v", artifacts)
	}
}

func TestListArtifacts_Error(t *testing.T) {
	svc := &mockMetadataService{
		getArtifactInfoFn: func(component, build, suite string, fromTime, toTime time.Time) ([]models.ArtifactInfo, error) {
			return nil, errSome
		},
	}
	h := NewArtifactsHandler(svc)

	req := httptest.NewRequest("GET", "/artifacts", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "admin"})
	rec := httptest.NewRecorder()
	h.List(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rec.Code)
	}
}

func TestListArtifacts_EmptyResult(t *testing.T) {
	svc := &mockMetadataService{
		getArtifactInfoFn: func(component, build, suite string, fromTime, toTime time.Time) ([]models.ArtifactInfo, error) {
			return nil, nil
		},
	}
	h := NewArtifactsHandler(svc)

	req := httptest.NewRequest("GET", "/artifacts", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "admin"})
	rec := httptest.NewRecorder()
	h.List(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	if body != "[]\n" && body != "[]" {
		t.Errorf("expected empty JSON array, got %q", body)
	}
}
