package handlers

import (
	"net/http"

	"github.com/sindwrr/test_storage/internal/metadata"
	"github.com/sindwrr/test_storage/internal/storage"
)

type uploadHandler struct {
	storage  storage.StorageService
	metadata metadata.MetadataService
}

func NewUploadHandler(storage storage.StorageService, metadata metadata.MetadataService) *uploadHandler {
	return &uploadHandler{storage: storage, metadata: metadata}
}

// @Summary      Загрузить файл
// @Description  (Еще не реализован!) Загружает файл артефакта и метаданные.
// @Tags         upload
// @Success      200  {string}  string  "Файл успешно загружен"
// @Failure      500  {object}  map[string]string  "Ошибка сервера"
// @Router       /upload [post]
func (*uploadHandler) Handle(w http.ResponseWriter, r *http.Request) {
	return // TODO: реализовать хендлер
}
