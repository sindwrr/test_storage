package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockPreviewService struct{}

func (m *mockPreviewService) ServePreview(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestPreviewHandler_ServePreview_Success(t *testing.T) {
	handler := NewPreviewHandler(&mockPreviewService{})
	req := httptest.NewRequest("GET", "/preview?id=1", nil)
	rec := httptest.NewRecorder()
	handler.ServePreview(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}
