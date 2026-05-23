package handlers

import (
	"net/http"

	"github.com/sindwrr/test_storage/internal/preview"
)

type PreviewHandler struct {
	svc preview.PreviewService
}

func NewPreviewHandler(svc preview.PreviewService) *PreviewHandler {
	return &PreviewHandler{svc: svc}
}

// @Summary      Предпросмотр артефакта
// @Description  Возвращает содержимое файла-артефакта по его ID для отображения в браузере.
// @Tags         artifacts
// @Param        id  query  int  true  "Идентификатор артефакта"
// @Produce      octet-stream
// @Success      200  {file}   file    "Содержимое артефакта"
// @Failure      400  {object}  map[string]string  "Некорректный запрос (отсутствует или невалидный ID)"
// @Failure      404  {object}  map[string]string  "Артефакт не найден"
// @Failure      500  {object}  map[string]string  "Внутренняя ошибка сервера"
// @Router       /preview [get]
func (h *PreviewHandler) ServePreview(w http.ResponseWriter, r *http.Request) {
	h.svc.ServePreview(w, r)
}
