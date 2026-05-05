package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/sindwrr/test_storage/internal/analytics"
)

type AnalyticsHandler struct {
	svc analytics.AnalyticsService
}

func NewAnalyticsHandler(svc analytics.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{svc: svc}
}

// @Summary      Количество артефактов по дням
// @Description  Возвращает агрегированные данные о количестве загруженных артефактов с группировкой по дням.
// @Tags         analytics
// @Produce      json
// @Success      200  {array}   analytics.DayCount  "Успешный ответ"
// @Failure      500  {object}  map[string]string   "Внутренняя ошибка сервера"
// @Router       /analytics/artifacts-per-day [get]
func (h *AnalyticsHandler) ArtifactsPerDay(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.ArtifactsPerDay(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// @Summary      Распределение статусов тестовых прогонов
// @Description  Возвращает распределение прогонов по статусам (Passed, Failed).
// @Tags         analytics
// @Produce      json
// @Success      200  {array}   analytics.StatusCount  "Успешный ответ"
// @Failure      500  {object}  map[string]string      "Внутренняя ошибка сервера"
// @Router       /analytics/status-distribution [get]
func (h *AnalyticsHandler) StatusDistribution(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.StatusDistribution(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// @Summary      Страница аналитики
// @Description  Возвращает страницу с аналитикой в виде графиков.
// @Tags         analytics
// @Produce      json
// @Success      200  {string}  string  "HTML-страница"
// @Failure      500  {string}  string  "Ошибка сервера"
// @Router       /analytics [get]
func AnalyticsPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/analytics.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
