package handlers

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUploadHandler_MethodNotAllowed(t *testing.T) {
	svc := &mockStorageService{}
	meta := &mockMetadataService{}
	h := NewUploadHandler(svc, meta, 1024)
	req := httptest.NewRequest(http.MethodGet, "/upload", nil)
	rec := httptest.NewRecorder()
	h.Handle(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}

func TestUploadHandler_MissingFile(t *testing.T) {
	svc := &mockStorageService{}
	meta := &mockMetadataService{}
	h := NewUploadHandler(svc, meta, 1024)
	req := httptest.NewRequest(http.MethodPost, "/upload?component=c&build=b&suite=s", nil)
	rec := httptest.NewRecorder()
	h.Handle(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing file, got %d", rec.Code)
	}
}

func TestUploadHandler_MissingComponent(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	part.Write([]byte("hello"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload?build=b&suite=s", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()
	h := NewUploadHandler(&mockStorageService{}, &mockMetadataService{}, 1024)
	h.Handle(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing component, got %d", rec.Code)
	}
}

func TestUploadHandler_FileTooLarge(t *testing.T) {
	h := NewUploadHandler(&mockStorageService{}, &mockMetadataService{}, 5) // small file size limit
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.bin")
	part.Write([]byte("aaaaaa"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload?component=c&build=b&suite=s", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()
	h.Handle(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for file too large, got %d", rec.Code)
	}
}

func TestUploadHandler_Success(t *testing.T) {
	svc := &mockStorageService{
		saveFn: func(file multipart.File, header *multipart.FileHeader) (string, error) {
			return "/tmp/testfile.bin", nil
		},
	}
	meta := &mockMetadataService{
		createArtifactFn: func(filePath string, fileSize int64, component, build, suite, result string) error {
			return nil
		},
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "testfile.bin")
	part.Write([]byte("content"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload?component=core&build=v1&suite=smoke&result=Passed", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()
	h := NewUploadHandler(svc, meta, 1<<20)
	h.Handle(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rec.Code)
	}
	if rec.Body.String() != "File was successfully uploaded!" {
		t.Errorf("unexpected body: %q", rec.Body.String())
	}
}

func TestUploadHandler_SaveError(t *testing.T) {
	svc := &mockStorageService{
		saveFn: func(file multipart.File, header *multipart.FileHeader) (string, error) {
			return "", errSome
		},
	}
	meta := &mockMetadataService{}
	h := NewUploadHandler(svc, meta, 1<<20)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "file.bin")
	part.Write([]byte("data"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload?component=c&build=b&suite=s&result=r", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()
	h.Handle(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rec.Code)
	}
}
