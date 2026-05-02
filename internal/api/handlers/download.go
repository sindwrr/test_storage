package handlers

import (
	"fmt"
	"io"
	"net/http"
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

	file, err := h.storage.Open(filePath)
	if err != nil {
		http.Error(w, `{"error":"file not found"}`, http.StatusNotFound)
		return
	}
	defer file.Close()

	buf := make([]byte, 512)
	n, _ := file.Read(buf)
	contentType := http.DetectContentType(buf[:n])

	if seeker, ok := file.(io.Seeker); ok {
		seeker.Seek(0, io.SeekStart)
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(filePath)))

	io.Copy(w, file)
}
