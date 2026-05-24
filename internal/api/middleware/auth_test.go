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

func TestRequireAuth_PassesWithValidCookie(t *testing.T) {
	nextCalled := false
	handler := RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		username, ok := r.Context().Value(UserContextKey).(string)
		if !ok || username != "testuser" {
			t.Errorf("expected username 'testuser' in context, got %q", username)
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "testuser"})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if !nextCalled {
		t.Error("expected next handler to be called")
	}
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}
