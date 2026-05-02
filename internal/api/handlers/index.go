package handlers

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/sindwrr/test_storage/internal/metadata"
	"github.com/sindwrr/test_storage/internal/models"
)

type IndexHandler struct {
	metaSvc metadata.MetadataService
	tmpl    *template.Template
}

func NewIndexHandler(metaSvc metadata.MetadataService) *IndexHandler {
	tmpl := template.Must(template.ParseFiles("web/templates/index.html"))
	return &IndexHandler{metaSvc: metaSvc, tmpl: tmpl}
}

// @Summary      Главная страница
// @Description  Отображает главную HTML-страницу для авторизованного пользователя.
// @Tags         general
// @Produce      html
// @Success      200  {string}  string  "HTML-страница"
// @Failure      500  {string}  string  "Ошибка сервера"
// @Router       / [get]
func (h *IndexHandler) Handle(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	component := r.URL.Query().Get("component")
	build := r.URL.Query().Get("build")
	suite := r.URL.Query().Get("suite")

	// TODO: Add time filters
	fromTime := time.Time{}
	toTime := time.Now()

	artifacts, err := h.metaSvc.GetArtifactInfo(component, build, suite, fromTime, toTime)
	log.Printf("Found %d artifacts", len(artifacts))
	if err != nil {
		http.Error(w, "Failed to load artifacts", http.StatusInternalServerError)
		return
	}

	data := struct {
		Username  string
		Artifacts []models.ArtifactInfo
	}{
		Username:  cookie.Value,
		Artifacts: artifacts,
	}

	if err := h.tmpl.Execute(w, data); err != nil {
		http.Error(w, "Render error", http.StatusInternalServerError)
		return
	}
}
