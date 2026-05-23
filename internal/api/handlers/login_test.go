package handlers

import (
	"errors"
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShowLogin_ReturnsOK(t *testing.T) {
	tmpl := template.Must(template.New("login").Parse("<html>Login</html>"))
	h := &LoginHandler{
		auth: &mockAuthService{},
		tmpl: tmpl,
	}
	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	rec := httptest.NewRecorder()
	h.ShowLogin(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	if rec.Body.String() != "<html>Login</html>" {
		t.Errorf("unexpected body: %q", rec.Body.String())
	}
}

func TestLogin_Success(t *testing.T) {
	svc := &mockAuthService{
		validateFn: func(u, p string) bool {
			return u == "admin" && p == "123"
		},
	}
	tmpl := template.Must(template.New("login").Parse("irrelevant"))
	h := &LoginHandler{
		auth: svc,
		tmpl: tmpl,
	}

	body := "username=admin&password=123"
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	h.Login(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303 See Other, got %d", rec.Code)
	}

	setCookie := rec.Header().Get("Set-Cookie")
	if !strings.Contains(setCookie, "session=admin") {
		t.Errorf("expected session cookie, got %q", setCookie)
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	svc := &mockAuthService{
		validateFn: func(u, p string) bool { return false },
	}
	tmpl := template.Must(template.New("login").Parse("irrelevant"))
	h := &LoginHandler{
		auth: svc,
		tmpl: tmpl,
	}

	body := "username=admin&password=wrong"
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	h.Login(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestHandle_MethodNotAllowed(t *testing.T) {
	tmpl := template.Must(template.New("login").Parse("irrelevant"))
	h := &LoginHandler{
		auth: &mockAuthService{},
		tmpl: tmpl,
	}
	req := httptest.NewRequest(http.MethodPut, "/login", nil)
	rec := httptest.NewRecorder()
	h.Handle(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

type errReader struct{}

func (e errReader) Read(p []byte) (int, error) {
	return 0, errors.New("read error")
}

func TestLogin_ParseFormError(t *testing.T) {
	h := &LoginHandler{
		auth: &mockAuthService{},
		tmpl: template.Must(template.New("login").Parse("irrelevant")),
	}
	req := httptest.NewRequest(http.MethodPost, "/login", nil)
	req.Body = io.NopCloser(errReader{})
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	h.Login(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestShowLogin_TemplateError(t *testing.T) {
	tmpl := template.Must(template.New("login").Parse(`{{ template "nonexistent" }}`))
	h := &LoginHandler{
		auth: &mockAuthService{},
		tmpl: tmpl,
	}
	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	rec := httptest.NewRecorder()
	h.ShowLogin(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
