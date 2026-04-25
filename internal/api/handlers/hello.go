package handlers

import (
	"html/template"
	"net/http"
)

// @Summary      Главная страница
// @Description  Отображает главную HTML-страницу для авторизованного пользователя.
// @Tags         general
// @Produce      html
// @Success      200  {string}  string  "HTML-страница"
// @Failure      500  {string}  string  "Ошибка сервера"
// @Router       / [get]
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	username := cookie.Value
	data := struct {
		Username string
	}{
		Username: username,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Render error", http.StatusInternalServerError)
		return
	}
}
