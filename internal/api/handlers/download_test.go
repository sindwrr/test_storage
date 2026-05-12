package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDownloadHandler_Success(t *testing.T) {
	meta := &mockMetadataService{
		getFilePathByIDFn: func(ctx context.Context, id int64) (string, error) {
			return "/tmp/testfile.log", nil
		},
	}
	store := &mockStorageService{}
	h := NewDownloadHandler(meta, store)

	req := httptest.NewRequest("GET", "/artifact/download/{id}", nil)
	req.SetPathValue("id", "1")
	req.AddCookie(&http.Cookie{Name: "session", Value: "admin"})
	rec := httptest.NewRecorder()
	h.Handle(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestDownloadHandler_MissingFile(t *testing.T) {
	meta := &mockMetadataService{
		getFilePathByIDFn: func(ctx context.Context, id int64) (string, error) {
			return "", errors.New("not found")
		},
	}
	store := &mockStorageService{}
	h := NewDownloadHandler(meta, store)

	req := httptest.NewRequest("GET", "/artifact/download/{id}", nil)
	req.SetPathValue("id", "999")
	req.AddCookie(&http.Cookie{Name: "session", Value: "admin"})
	rec := httptest.NewRecorder()
	h.Handle(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}
