package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sindwrr/test_storage/internal/health"
)

// @Summary      Проверка, жив ли сервис
// @Description  Возвращает статус "alive", если сервис работает.
// @Tags         health
// @Produce      json
// @Success      200  {object}  map[string]string  "Alive"
// @Router       /health/alive [get]
func AliveHandler(svc health.HealthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := svc.Alive(r.Context()); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]string{"status": "not alive"})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "alive"})
	}
}

// @Summary      Проверка готовности сервиса
// @Description  Проверяет доступность БД и других зависимостей.
// @Tags         health
// @Produce      json
// @Success      200  {object}  map[string]string  "Ready"
// @Failure      503  {object}  map[string]string  "Not ready"
// @Router       /health/ready [get]
func ReadyHandler(svc health.HealthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := svc.Ready(r.Context()); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]string{"status": "not ready"})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
	}
}
