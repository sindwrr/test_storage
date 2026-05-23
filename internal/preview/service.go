package preview

import (
	"context"
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/sindwrr/test_storage/internal/metadata"
	"github.com/sindwrr/test_storage/internal/storage"
)

type previewService struct {
	metaSvc  metadata.MetadataService
	storeSvc storage.StorageService
}

func NewService(metaSvc metadata.MetadataService, storeSvc storage.StorageService) PreviewService {
	return &previewService{metaSvc: metaSvc, storeSvc: storeSvc}
}

func (s *previewService) ServePreview(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "id parameter is required", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	filePath, err := s.metaSvc.GetFilePathByID(context.Background(), id)
	if err != nil {
		http.Error(w, "artifact not found", http.StatusNotFound)
		return
	}

	reader, err := s.storeSvc.Open(filePath)
	if err != nil {
		http.Error(w, "failed to open file", http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	ext := filepath.Ext(filePath)
	var contentType string
	switch ext {
	case ".txt", ".log", ".csv":
		contentType = "text/plain; charset=utf-8"
	case ".json":
		contentType = "application/json"
	case ".xml":
		contentType = "application/xml"
	case ".png":
		contentType = "image/png"
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".bmp":
		contentType = "image/bmp"
	case ".gif":
		contentType = "image/gif"
	default:
		contentType = "application/octet-stream"
	}

	w.Header().Set("Content-Type", contentType)
	if _, err := io.Copy(w, reader); err != nil {
		http.Error(w, "failed to stream file", http.StatusInternalServerError)
		return
	}
}
