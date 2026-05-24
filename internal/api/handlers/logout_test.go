package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogoutHandler_DeactivatesUser(t *testing.T) {
	var calledSetActive bool
	var capturedUsername string
	var capturedActive bool

	authSvc := &mockAuthService{
		setActiveFn: func(username string, active bool) {
			calledSetActive = true
			capturedUsername = username
			capturedActive = active
		},
	}
	handler := NewLogoutHandler(authSvc)

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "testuser"})
	rec := httptest.NewRecorder()
	handler.Handle(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected status %d, got %d", http.StatusSeeOther, rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "/login" {
		t.Errorf("expected redirect to /login, got %q", loc)
	}

	cookies := rec.Result().Cookies()
	var cleared bool
	for _, c := range cookies {
		if c.Name == "session" && c.Value == "" && c.MaxAge < 0 {
			cleared = true
			break
		}
	}
	if !cleared {
		t.Error("expected session cookie to be cleared")
	}

	if !calledSetActive {
		t.Error("expected SetUserActive to be called")
	}
	if capturedUsername != "testuser" {
		t.Errorf("expected username 'testuser', got %q", capturedUsername)
	}
	if capturedActive != false {
		t.Errorf("expected active=false, got %v", capturedActive)
	}
}

func TestLogoutHandler_NoCookie(t *testing.T) {
	var calledSetActive bool
	authSvc := &mockAuthService{
		setActiveFn: func(username string, active bool) {
			calledSetActive = true
		},
	}
	handler := NewLogoutHandler(authSvc)

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	rec := httptest.NewRecorder()
	handler.Handle(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected status %d, got %d", http.StatusSeeOther, rec.Code)
	}
	if calledSetActive {
		t.Error("SetUserActive should not be called when there is no session cookie")
	}
}

func TestLogoutHandler_MethodNotAllowed(t *testing.T) {
	handler := NewLogoutHandler(&mockAuthService{})
	req := httptest.NewRequest(http.MethodGet, "/logout", nil)
	rec := httptest.NewRecorder()
	handler.Handle(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
	}
}
