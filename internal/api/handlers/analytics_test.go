package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/sindwrr/test_storage/internal/models/analytics"
)

func TestArtifactsPerDay_ReturnsJSON(t *testing.T) {
	svc := &mockAnalyticsService{
		artifactsPerDayFn: func(ctx context.Context) ([]analytics.DayCount, error) {
			return []analytics.DayCount{{Date: "2026-05-01", Count: 5}}, nil
		},
	}
	h := NewAnalyticsHandler(svc)
	req := httptest.NewRequest("GET", "/analytics/artifacts-per-day", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "admin"})
	rec := httptest.NewRecorder()
	h.ArtifactsPerDay(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	if rec.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", rec.Header().Get("Content-Type"))
	}
}

func TestArtifactsPerDay_Error(t *testing.T) {
	svc := &mockAnalyticsService{
		artifactsPerDayFn: func(ctx context.Context) ([]analytics.DayCount, error) {
			return nil, errSome
		},
	}
	h := NewAnalyticsHandler(svc)
	req := httptest.NewRequest("GET", "/analytics/artifacts-per-day", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "admin"})
	rec := httptest.NewRecorder()
	h.ArtifactsPerDay(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rec.Code)
	}
}

func TestStatusDistribution_ReturnsJSON(t *testing.T) {
	svc := &mockAnalyticsService{
		statusDistFn: func(ctx context.Context) ([]analytics.StatusCount, error) {
			return []analytics.StatusCount{{Status: "passed", Count: 10}}, nil
		},
	}
	h := NewAnalyticsHandler(svc)
	req := httptest.NewRequest("GET", "/analytics/status-distribution", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "admin"})
	rec := httptest.NewRecorder()
	h.StatusDistribution(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestAnalyticsPageHandler_Success(t *testing.T) {
	tmpDir := t.TempDir()
	templatesDir := filepath.Join(tmpDir, "web", "templates")
	err := os.MkdirAll(templatesDir, 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(filepath.Join(templatesDir, "analytics.html"), []byte("<html></html>"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tmpDir)

	req := httptest.NewRequest("GET", "/analytics", nil)
	rec := httptest.NewRecorder()
	AnalyticsPageHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestAnalyticsPageHandler_TemplateError(t *testing.T) {
	req := httptest.NewRequest("GET", "/analytics", nil)
	rec := httptest.NewRecorder()
	AnalyticsPageHandler(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rec.Code)
	}
}
