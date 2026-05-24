package middleware

import (
	"net/http"

	"github.com/sindwrr/test_storage/internal/auth"
)

func RequireAdmin(authSvc auth.AuthService) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			username, ok := r.Context().Value(UserContextKey).(string)
			if !ok || username == "" {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			groupID, err := authSvc.GetUserGroup(username)
			if err != nil || groupID != 2 {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		}
	}
}
