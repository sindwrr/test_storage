package handlers

import "net/http"

// @Summary      Выйти из системы
// @Description  Завершает сеанс пользователя и перенаправляет на страницу входа.
// @Tags         auth
// @Success      303  "Редирект на /login"
// @Failure      405  {object}  map[string]string  "Метод не разрешен"
// @Router       /logout [post]
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Logout: method not allowed", http.StatusMethodNotAllowed)
		return
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
