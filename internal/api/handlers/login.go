package handlers

import (
	"html/template"
	"net/http"

	"github.com/sindwrr/test_storage/internal/auth"
)

type LoginHandler struct {
	auth auth.AuthService
	tmpl *template.Template
}

func NewLoginHandler(a auth.AuthService) *LoginHandler {
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

// @Summary      Страница входа
// @Description  Отображает HTML-страницу с формой для входа в систему.
// @Tags         auth
// @Produce      html
// @Success      200  {string}  string  "HTML-страница входа"
// @Failure      500  {object}  map[string]string  "Ошибка сервера"
// @Router       /login [get]
func (h *LoginHandler) ShowLogin(w http.ResponseWriter, r *http.Request) {
	err := h.tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Login: template error", http.StatusInternalServerError)
	}
}

// @Summary      Войти в систему
// @Description  Аутентифицирует пользователя по логину и паролю из HTML-формы.
// @Tags         auth
// @Accept       application/x-www-form-urlencoded
// @Param        username  formData  string  true  "Логин"
// @Param        password  formData  string  true  "Пароль"
// @Success      303  "Редирект на главную страницу"
// @Failure      400  {string}  string  "Плохой запрос"
// @Failure      401  {string}  string  "Неверные учетные данные"
// @Router       /login [post]
func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Login: bad request", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if !h.auth.Validate(username, password) {
		http.Error(w, "Login: invalid credentials", http.StatusUnauthorized)
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
