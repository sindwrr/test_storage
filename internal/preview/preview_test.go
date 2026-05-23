package preview

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/sindwrr/test_storage/internal/models"
	"github.com/stretchr/testify/assert"
)

type mockMetadataService struct {
	getFilePathByIDFn func(ctx context.Context, id int64) (string, error)
}

func (m *mockMetadataService) CreateArtifact(string, int64, string, string, string, string) error {
	return nil
}
func (m *mockMetadataService) GetArtifactInfo(string, string, string, time.Time, time.Time) ([]models.ArtifactInfo, error) {
	return nil, nil
}
func (m *mockMetadataService) GetFilePathByID(ctx context.Context, id int64) (string, error) {
	return m.getFilePathByIDFn(ctx, id)
}

type mockStorageService struct {
	openFn func(path string) (io.ReadCloser, error)
}

func (m *mockStorageService) Save(multipart.File, *multipart.FileHeader) (string, error) {
	return "", nil
}
func (m *mockStorageService) Open(path string) (io.ReadCloser, error) {
	return m.openFn(path)
}

type errReader struct{}

func (e errReader) Read(p []byte) (int, error) { return 0, errors.New("read error") }
func (e errReader) Close() error               { return nil }

func TestServePreview_Success(t *testing.T) {
	metaSvc := &mockMetadataService{
		getFilePathByIDFn: func(ctx context.Context, id int64) (string, error) {
			return "/tmp/test.log", nil
		},
	}
	storeSvc := &mockStorageService{
		openFn: func(path string) (io.ReadCloser, error) {
			return io.NopCloser(strings.NewReader("log content")), nil
		},
	}
	svc := NewService(metaSvc, storeSvc)

	req := httptest.NewRequest("GET", "/preview?id=1", nil)
	rec := httptest.NewRecorder()
	svc.ServePreview(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "text/plain; charset=utf-8", rec.Header().Get("Content-Type"))
	assert.Equal(t, "log content", rec.Body.String())
}

func TestServePreview_MissingID(t *testing.T) {
	svc := NewService(&mockMetadataService{}, &mockStorageService{})
	req := httptest.NewRequest("GET", "/preview", nil)
	rec := httptest.NewRecorder()
	svc.ServePreview(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "id parameter is required")
}

func TestServePreview_InvalidID(t *testing.T) {
	svc := NewService(&mockMetadataService{}, &mockStorageService{})
	req := httptest.NewRequest("GET", "/preview?id=abc", nil)
	rec := httptest.NewRecorder()
	svc.ServePreview(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid id")
}

func TestServePreview_ArtifactNotFound(t *testing.T) {
	metaSvc := &mockMetadataService{
		getFilePathByIDFn: func(ctx context.Context, id int64) (string, error) {
			return "", errors.New("not found")
		},
	}
	svc := NewService(metaSvc, &mockStorageService{})
	req := httptest.NewRequest("GET", "/preview?id=1", nil)
	rec := httptest.NewRecorder()
	svc.ServePreview(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "artifact not found")
}

func TestServePreview_FileOpenError(t *testing.T) {
	metaSvc := &mockMetadataService{
		getFilePathByIDFn: func(ctx context.Context, id int64) (string, error) {
			return "/tmp/test.log", nil
		},
	}
	storeSvc := &mockStorageService{
		openFn: func(path string) (io.ReadCloser, error) {
			return nil, errors.New("permission denied")
		},
	}
	svc := NewService(metaSvc, storeSvc)
	req := httptest.NewRequest("GET", "/preview?id=1", nil)
	rec := httptest.NewRecorder()
	svc.ServePreview(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "failed to open file")
}

func TestServePreview_StreamError(t *testing.T) {
	metaSvc := &mockMetadataService{
		getFilePathByIDFn: func(ctx context.Context, id int64) (string, error) {
			return "/tmp/test.log", nil
		},
	}
	storeSvc := &mockStorageService{
		openFn: func(path string) (io.ReadCloser, error) {
			return errReader{}, nil
		},
	}
	svc := NewService(metaSvc, storeSvc)
	req := httptest.NewRequest("GET", "/preview?id=1", nil)
	rec := httptest.NewRecorder()
	svc.ServePreview(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "failed to stream file")
}

func TestServePreview_ContentType(t *testing.T) {
	tests := []struct {
		ext         string
		contentType string
	}{
		{".txt", "text/plain; charset=utf-8"},
		{".log", "text/plain; charset=utf-8"},
		{".csv", "text/plain; charset=utf-8"},
		{".json", "application/json"},
		{".xml", "application/xml"},
		{".png", "image/png"},
		{".jpg", "image/jpeg"},
		{".jpeg", "image/jpeg"},
		{".bmp", "image/bmp"},
		{".gif", "image/gif"},
		{".pdf", "application/octet-stream"},
	}

	for _, tt := range tests {
		t.Run(tt.ext, func(t *testing.T) {
			metaSvc := &mockMetadataService{
				getFilePathByIDFn: func(ctx context.Context, id int64) (string, error) {
					return "/tmp/file" + tt.ext, nil
				},
			}
			storeSvc := &mockStorageService{
				openFn: func(path string) (io.ReadCloser, error) {
					return io.NopCloser(strings.NewReader("data")), nil
				},
			}
			svc := NewService(metaSvc, storeSvc)
			req := httptest.NewRequest("GET", "/preview?id=1", nil)
			rec := httptest.NewRecorder()
			svc.ServePreview(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, tt.contentType, rec.Header().Get("Content-Type"))
		})
	}
}
