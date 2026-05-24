package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockAuthService struct {
	getGroupFn func(username string) (int, error)
}

func (m *mockAuthService) Validate(username, password string) bool    { return false }
func (m *mockAuthService) SetUserActive(username string, active bool) {}
func (m *mockAuthService) GetUserGroup(username string) (int, error) {
	if m.getGroupFn != nil {
		return m.getGroupFn(username)
	}
	return 0, nil
}

func TestRequireAdmin_AdminAccess(t *testing.T) {
	svc := &mockAuthService{getGroupFn: func(u string) (int, error) { return 2, nil }}
	adminMiddleware := RequireAdmin(svc)

	nextCalled := false
	handler := adminMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	ctx := context.WithValue(req.Context(), UserContextKey, "admin_user")
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.True(t, nextCalled, "expected next handler to be called for admin")
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRequireAdmin_RegularUserForbidden(t *testing.T) {
	svc := &mockAuthService{getGroupFn: func(u string) (int, error) { return 1, nil }}
	adminMiddleware := RequireAdmin(svc)

	nextCalled := false
	handler := adminMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	}))

	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	ctx := context.WithValue(req.Context(), UserContextKey, "regular_user")
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.False(t, nextCalled, "expected next handler not to be called for non‑admin")
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestRequireAdmin_NoUserContext(t *testing.T) {
	svc := &mockAuthService{getGroupFn: func(u string) (int, error) { return 2, nil }}
	adminMiddleware := RequireAdmin(svc)

	nextCalled := false
	handler := adminMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	}))

	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.False(t, nextCalled, "expected next handler not to be called without context")
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestRequireAdmin_DBError(t *testing.T) {
	svc := &mockAuthService{getGroupFn: func(u string) (int, error) { return 0, errors.New("db error") }}
	adminMiddleware := RequireAdmin(svc)

	nextCalled := false
	handler := adminMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	}))

	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	ctx := context.WithValue(req.Context(), UserContextKey, "admin_user")
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.False(t, nextCalled, "expected next handler not to be called on error")
	assert.Equal(t, http.StatusForbidden, rec.Code)
}
