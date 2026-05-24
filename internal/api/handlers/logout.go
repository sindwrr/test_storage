package handlers

import (
	"net/http"

	"github.com/sindwrr/test_storage/internal/auth"
)

type LogoutHandler struct {
	authSvc auth.AuthService
}

func NewLogoutHandler(authSvc auth.AuthService) *LogoutHandler {
	return &LogoutHandler{authSvc: authSvc}
}

// @Summary      Выйти из системы
// @Description  Завершает сеанс пользователя и перенаправляет на страницу входа.
// @Tags         auth
// @Success      303  "Редирект на /login"
// @Failure      405  {object}  map[string]string  "Метод не разрешен"
// @Router       /logout [post]
func (h *LogoutHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Logout: method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if cookie, err := r.Cookie("session"); err == nil && cookie.Value != "" {
		h.authSvc.SetUserActive(cookie.Value, false)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
