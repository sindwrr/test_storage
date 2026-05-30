package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestDownloadHandler_Success(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "testfile.log")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	meta := &mockMetadataService{
		getFilePathByIDFn: func(ctx context.Context, id int64) (string, error) {
			return tmpFile.Name(), nil
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
