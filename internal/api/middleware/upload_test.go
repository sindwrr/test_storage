package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequireUpload_AdminAllowed(t *testing.T) {
	svc := &mockAuthService{getGroupFn: func(u string) (int, error) { return 2, nil }}
	mw := RequireUpload(svc)
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/upload", nil)
	ctx := context.WithValue(req.Context(), UserContextKey, "admin")
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRequireUpload_ATFAllowed(t *testing.T) {
	svc := &mockAuthService{getGroupFn: func(u string) (int, error) { return 3, nil }}
	mw := RequireUpload(svc)
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/upload", nil)
	ctx := context.WithValue(req.Context(), UserContextKey, "atf_user")
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRequireUpload_RegularUserForbidden(t *testing.T) {
	svc := &mockAuthService{getGroupFn: func(u string) (int, error) { return 1, nil }}
	mw := RequireUpload(svc)
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("should not be called")
	}))

	req := httptest.NewRequest(http.MethodPost, "/upload", nil)
	ctx := context.WithValue(req.Context(), UserContextKey, "regular_user")
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestRequireUpload_NoUserContext(t *testing.T) {
	svc := &mockAuthService{getGroupFn: func(u string) (int, error) { return 2, nil }}
	mw := RequireUpload(svc)
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("should not be called")
	}))

	req := httptest.NewRequest(http.MethodPost, "/upload", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestRequireUpload_DBError(t *testing.T) {
	svc := &mockAuthService{getGroupFn: func(u string) (int, error) { return 0, errors.New("db error") }}
	mw := RequireUpload(svc)
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("should not be called")
	}))

	req := httptest.NewRequest(http.MethodPost, "/upload", nil)
	ctx := context.WithValue(req.Context(), UserContextKey, "admin")
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}
