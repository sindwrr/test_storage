package api

import (
	"net/http"

	"github.com/sindwrr/test_storage/internal/api/handlers"
	"github.com/sindwrr/test_storage/internal/api/middleware"
	"github.com/sindwrr/test_storage/internal/auth"
	"github.com/sindwrr/test_storage/internal/config"
)

func NewRouter(cfg config.Config) http.Handler {
	mux := http.NewServeMux()

	authService := auth.NewService()
	loginHandler := handlers.NewLoginHandler(authService)

	// healthSvc := health.NewService(db)

	mux.HandleFunc("/login", loginHandler.Handle)
	mux.HandleFunc("/logout", handlers.LogoutHandler)
	mux.HandleFunc("/", middleware.RequireAuth(handlers.HelloHandler))

	return mux
}
