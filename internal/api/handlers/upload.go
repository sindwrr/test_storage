package handlers

import (
	"net/http"

	"github.com/sindwrr/test_storage/internal/metadata"
	"github.com/sindwrr/test_storage/internal/storage"
)

type uploadHandler struct {
	storage      storage.StorageService
	metadata     metadata.MetadataService
	maxFileBytes int64
}

func NewUploadHandler(storage storage.StorageService, metadata metadata.MetadataService, maxFileBytes int64) *uploadHandler {
	return &uploadHandler{storage: storage, metadata: metadata, maxFileBytes: maxFileBytes}
}

// @Summary      Загрузить файл
// @Description  Загружает в систему файл артефакта и метаданные.
// @Tags         upload
// @Accept       multipart/form-data
// @Produce      plain
// @Param        file       formData  file    true  "Файл артефакта"
// @Param        component  query     string  true  "Тестируемый компонент"
// @Param        build      query     string  true  "Номер сборки"
// @Param        suite      query     string  true  "Набор тестов"
// @Success      201  {string}  string  "Файл успешно загружен"
// @Failure      400  {string}  string  "Ошибка в запросе"
// @Failure      500  {string}  string  "Внутренняя ошибка сервера"
// @Router       /upload [post]
func (h *uploadHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Upload: method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, h.maxFileBytes)
	if err := r.ParseMultipartForm(h.maxFileBytes); err != nil {
		http.Error(w, "File too large or bad multipart form", http.StatusBadRequest)
		return
	}
	defer r.MultipartForm.RemoveAll()

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Missing file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	component := r.FormValue("component")
	if component == "" {
		http.Error(w, "Missing component", http.StatusBadRequest)
		return
	}

	build := r.FormValue("build")
	if build == "" {
		http.Error(w, "Missing build", http.StatusBadRequest)
		return
	}

	suite := r.FormValue("suite")
	if suite == "" {
		http.Error(w, "Missing test suite", http.StatusBadRequest)
		return
	}

	_, err = h.storage.Save(file, header)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// TODO: Add insert to DB in metadata service

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("File was successfully uploaded!"))
}
