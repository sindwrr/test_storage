package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sindwrr/test_storage/internal/metadata"
	"github.com/sindwrr/test_storage/internal/models"
)

type ArtifactsHandler struct {
	metadataSvc metadata.MetadataService
}

func NewArtifactsHandler(metadataSvc metadata.MetadataService) *ArtifactsHandler {
	return &ArtifactsHandler{metadataSvc: metadataSvc}
}

// @Summary      Список артефактов с фильтрацией
// @Description  Возвращает список артефактов в формате JSON, с возможностью фильтрации по компоненту, сборке, набору тестов и временному диапазону.
// @Tags         artifacts
// @Produce      json
// @Param        component  query     string  false  "Название компонента"
// @Param        build      query     string  false  "Номер сборки"
// @Param        suite      query     string  false  "Название тестового набора"
// @Param        from       query     string  false  "Начальная граница времени загрузки (RFC3339)"
// @Param        to         query     string  false  "Конечная граница времени загрузки (RFC3339)"
// @Success      200        {array}   models.ArtifactInfo  "Список артефактов"
// @Failure      400        {object}  map[string]string    "Некорректные параметры запроса"
// @Failure      500        {object}  map[string]string    "Внутренняя ошибка сервера"
// @Router       /artifacts [get]
func (h *ArtifactsHandler) List(w http.ResponseWriter, r *http.Request) {
	component := r.URL.Query().Get("component")
	build := r.URL.Query().Get("build")
	suite := r.URL.Query().Get("suite")
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	var fromTime, toTime time.Time
	if fromStr != "" {
		fromTime, _ = time.Parse(time.RFC3339, fromStr)
	}
	if toStr != "" {
		toTime, _ = time.Parse(time.RFC3339, toStr)
	}

	artifacts, err := h.metadataSvc.GetArtifactInfo(component, build, suite, fromTime, toTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if artifacts == nil {
		artifacts = []models.ArtifactInfo{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artifacts)
}
