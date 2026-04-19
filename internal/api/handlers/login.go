package handlers

import (
	"html/template"
	"net/http"

	"github.com/sindwrr/test_storage/internal/auth"
)

type LoginHandler struct {
	auth *auth.Service
	tmpl *template.Template
}

func NewLoginHandler(a *auth.Service) *LoginHandler {
	tmpl := template.Must(template.ParseFiles("web/templates/login.html"))
	return &LoginHandler{
		auth: a,
		tmpl: tmpl,
	}
}

func (h *LoginHandler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.ShowLogin(w, r)
	case http.MethodPost:
		h.Login(w, r)
	default:
		http.Error(w, "Login: method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *LoginHandler) ShowLogin(w http.ResponseWriter, r *http.Request) {
	err := h.tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Login: template error", http.StatusInternalServerError)
	}
}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Login: bad request", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if !h.auth.Validate(username, password) {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    username,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   3600,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
