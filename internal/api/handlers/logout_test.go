package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogoutHandler_POST(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	rec := httptest.NewRecorder()
	LogoutHandler(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected status %d, got %d", http.StatusSeeOther, rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "/login" {
		t.Errorf("expected redirect to /login, got %q", loc)
	}

	cookies := rec.Result().Cookies()
	var sessionCleared bool
	for _, c := range cookies {
		if c.Name == "session" && c.Value == "" && c.MaxAge < 0 {
			sessionCleared = true
			break
		}
	}
	if !sessionCleared {
		t.Error("expected session cookie to be cleared")
	}
}

func TestLogoutHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/logout", nil)
	rec := httptest.NewRecorder()
	LogoutHandler(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
	}
}
