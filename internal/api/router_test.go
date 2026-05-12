package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sindwrr/test_storage/internal/config"
)

func TestNewRouter_ReturnsHandler(t *testing.T) {
	tmpDir := t.TempDir()
	templatesDir := filepath.Join(tmpDir, "web", "templates")
	os.MkdirAll(templatesDir, 0755)
	os.WriteFile(filepath.Join(templatesDir, "login.html"), []byte("<html>Login</html>"), 0644)
	os.WriteFile(filepath.Join(templatesDir, "index.html"), []byte("<html>Index</html>"), 0644)

	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Chdir: %v", err)
	}

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	cfg := config.Config{
		ArtifactVolume: filepath.Join(tmpDir, "artifacts"),
		MaxFileBytes:   1024,
	}

	router := NewRouter(db, cfg)
	if router == nil {
		t.Fatal("NewRouter returned nil")
	}

	tests := []struct {
		method string
		target string
	}{
		{"GET", "/login"},
		{"POST", "/login"},
		{"POST", "/logout"},
		{"GET", "/health/alive"},
		{"GET", "/health/ready"},
		{"GET", "/docs/index.html"},
	}

	for _, tc := range tests {
		req := httptest.NewRequest(tc.method, tc.target, nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		if rec.Code == http.StatusNotFound {
			t.Errorf("route %s %s should not return 404", tc.method, tc.target)
		}
	}
}
