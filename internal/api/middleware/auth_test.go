package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequireAuth_RedirectsToLogin(t *testing.T) {
	req := httptest.NewRequest("GET", "/artifacts", nil)
	rec := httptest.NewRecorder()
	handler := RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther && rec.Code != http.StatusFound {
		t.Errorf("expected redirect status, got %d", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "/login" {
		t.Errorf("expected redirect to /login, got %q", loc)
	}
}
