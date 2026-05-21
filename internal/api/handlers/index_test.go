package handlers

import (
	"errors"
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sindwrr/test_storage/internal/models"
)

func setupTestEnvironment(t *testing.T, templateContent string) (*IndexHandler, func()) {
	t.Helper()

	tmpDir := t.TempDir()
	templatesDir := filepath.Join(tmpDir, "web", "templates")
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(templatesDir, "index.html"), []byte(templateContent), 0644); err != nil {
		t.Fatalf("write template: %v", err)
	}

	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	cleanup := func() {
		os.Chdir(origDir)
	}

	tmpl := template.Must(template.ParseFiles("web/templates/index.html"))
	handler := &IndexHandler{
		metaSvc: &mockMetadataService{},
		tmpl:    tmpl,
	}

	return handler, cleanup
}

func TestIndexHandler_Handle_NoSessionRedirect(t *testing.T) {
	handler, cleanup := setupTestEnvironment(t, `<html>ok</html>`)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.Handle(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303 See Other, got %d", rec.Code)
	}
	loc := rec.Header().Get("Location")
	if loc != "/login" {
		t.Errorf("expected redirect to /login, got %q", loc)
	}
}

func TestIndexHandler_Handle_Success(t *testing.T) {
	handler, cleanup := setupTestEnvironment(t, `<html>{{ range .Artifacts }}{{ .FileName }}{{ end }}</html>`)
	defer cleanup()

	handler.metaSvc = &mockMetadataService{
		getArtifactInfoFn: func(component, build, suite string, fromTime, toTime time.Time) ([]models.ArtifactInfo, error) {
			return []models.ArtifactInfo{
				{FileName: "test.log", Component: "core"},
			}, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "demo"})
	rec := httptest.NewRecorder()

	handler.Handle(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "test.log") {
		t.Errorf("expected body to contain 'test.log', got %s", body)
	}
}

func TestIndexHandler_Handle_ArtifactLoadError(t *testing.T) {
	handler, cleanup := setupTestEnvironment(t, `<html>ok</html>`)
	defer cleanup()

	handler.metaSvc = &mockMetadataService{
		getArtifactInfoFn: func(component, build, suite string, fromTime, toTime time.Time) ([]models.ArtifactInfo, error) {
			return nil, errors.New("db error")
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "demo"})
	rec := httptest.NewRecorder()

	handler.Handle(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Failed to load artifacts") {
		t.Errorf("expected error message, got %s", rec.Body.String())
	}
}
