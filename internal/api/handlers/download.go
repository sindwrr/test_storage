package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/sindwrr/test_storage/internal/metadata"
	"github.com/sindwrr/test_storage/internal/storage"
)

type DownloadHandler struct {
	metadata metadata.MetadataService
	storage  storage.StorageService
}

func NewDownloadHandler(md metadata.MetadataService, st storage.StorageService) *DownloadHandler {
	return &DownloadHandler{
		metadata: md,
		storage:  st,
	}
}

// @Summary      Скачать артефакт
// @Description  Скачивает файл артефакта по его ID.
// @Tags         artifacts
// @Produce      octet-stream
// @Param        id   path      int  true  "ID артефакта"
// @Success      200  {file}    file  "Файл артефакта"
// @Failure      400  {object}  map[string]string  "Некорректный ID"
// @Failure      404  {object}  map[string]string  "Артефакт не найден"
// @Failure      500  {object}  map[string]string  "Внутренняя ошибка"
// @Router       /artifact/download/{id} [get]
func (h *DownloadHandler) Handle(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	filePath, err := h.metadata.GetFilePathByID(r.Context(), id)
	if err != nil {
		http.Error(w, `{"error":"artifact not found"}`, http.StatusNotFound)
		return
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, `{"error":"file not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(filePath)))
	w.Header().Set("X-Accel-Redirect", "/protected_files/"+filepath.Base(filePath))
	w.WriteHeader(http.StatusOK)
}
