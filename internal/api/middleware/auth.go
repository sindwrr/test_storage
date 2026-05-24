package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const UserContextKey contextKey = "username"

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), UserContextKey, cookie.Value)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
